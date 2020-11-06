package rsd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/cuigh/auxo/app"
)

var javaTypes = map[string]typeMap{
	"bool":     {"boolean", "Boolean"},
	"string":   {"String", "String"},
	"bytes":    {"byte[]", "byte[]"},
	"float":    {"float", "Float"},
	"double":   {"double", "Double"},
	"int32":    {"int", "Integer"},
	"uint32":   {"int", "Integer"},
	"sint32":   {"int", "Integer"},
	"fixed32":  {"int", "Integer"},
	"sfixed32": {"int", "Integer"},
	"int64":    {"long", "Long"},
	"uint64":   {"long", "Long"},
	"sint64":   {"long", "Long"},
	"fixed64":  {"long", "Long"},
	"sfixed64": {"long", "Long"},
}

type typeMap struct {
	PrimaryType string
	ObjectType  string
}

// EnumModel 枚举数据
type EnumModel struct {
	ToolVersion string
	Package     string
	*Enum
}

// DtoModel DTO 数据
type DtoModel struct {
	ToolVersion string
	Package     string
	Name        string
	Imports     []*Import
	Types       []*Type
	Errors      []*ErrorType
	Enums       []*Enum
}

func newDtoModel(pkg, name string, imports []*Import, types []*Type, errTypes []*ErrorType, enums []*Enum) *DtoModel {
	m := &DtoModel{
		ToolVersion: app.Version,
		Package:     pkg,
		Name:        name,
		Imports:     imports,
		Types:       types,
		Errors:      errTypes,
		Enums:       enums,
	}
	m.initImports()
	return m
}

func (m *DtoModel) initImports() {
	m.addImports("lark.pb.annotation.ProtoField", "lark.pb.annotation.ProtoMessage", "lark.pb.field.FieldType",
		"lombok.Getter", "lombok.Setter")
	if len(m.Errors) > 0 {
		m.addImports("lark.core.lang.Error")
	}
	for _, e := range m.Enums {
		m.addImports(m.Package + ".constant." + e.Name)
	}

	checks := map[string]bool{
		"List":          false,
		"Map":           false,
		"ZonedDateTime": false,
		"LocalDateTime": false,
		"LocalDate":     false,
		"LocalTime":     false,
		"BigDecimal":    false,
	}
	for _, t := range m.Types {
		for _, f := range t.Fields {
			if f.Modifier == "repeated" {
				checks["List"] = true
			}
			if strings.HasPrefix(f.FullType, "Map<") {
				checks["Map"] = true
			}

			matchType := ""
			for k, v := range checks {
				if !v && f.FullType == k {
					matchType = k
					break
				}
			}
			if matchType != "" {
				checks[matchType] = true
			}
		}
	}
	if checks["List"] {
		m.addImports("java.util.ArrayList", "java.util.Collection", "java.util.List")
	}
	if checks["Map"] {
		m.addImports("java.util.Map")
	}
	if checks["ZonedDateTime"] {
		m.addImports("java.time.ZonedDateTime")
	}
	if checks["LocalDateTime"] {
		m.addImports("java.time.LocalDateTime")
	}
	if checks["LocalDate"] {
		m.addImports("java.time.LocalDate")
	}
	if checks["LocalTime"] {
		m.addImports("java.time.LocalTime")
	}
	if checks["BigDecimal"] {
		m.addImports("java.math.BigDecimal")
	}

	sort.Slice(m.Imports, func(i, j int) bool {
		return m.Imports[i].Path < m.Imports[j].Path
	})
}

func (m *DtoModel) addImports(pathes ...string) {
	for _, path := range pathes {
		missing := true
		for _, imp := range m.Imports {
			if imp.Path == path {
				missing = false
				break
			}
		}
		if missing {
			m.Imports = append(m.Imports, newImport(path, ""))
		}
	}
}

// ServiceModel 服务数据
type ServiceModel struct {
	ToolVersion string
	Package     string
	*Service
}

// Service 服务信息
type Service struct {
	Name        string
	Alias       string
	Fail        string
	Description string
	Imports     []*Import
	Methods     []*Method
}

// Method 方法信息
type Method struct {
	Name        string
	Alias       string
	Invoke      string
	Description string
	Request     []*Parameter
	Response    *Parameter
}

// ErrorType 错误定义类型
type ErrorType struct {
	Name        string
	Description string
	Errors      []*Error
}

// Error 错误定义
type Error struct {
	Name    string `xml:"name,attr"`
	Code    string `xml:"code,attr"`
	Message string `xml:"message,attr"`
}

// Parameter 参数信息
type Parameter struct {
	Type string
	Name string
}

func newParameter(paramType, paramName string) *Parameter {
	return &Parameter{
		Type: paramType,
		Name: paramName,
	}
}

// Type 数据类型信息
type Type struct {
	Name        string
	Description string
	Extends     string
	Fields      []*Field
}

func newType(name, desc, extends string, fields []*Field) *Type {
	return &Type{
		Name:        name,
		Description: desc,
		Extends:     strings.Replace(strings.Replace(extends, "[", "<", -1), "]", ">", -1),
		Fields:      fields,
	}
}

// Field 字段信息
type Field struct {
	Name        string
	Modifier    string
	Order       string
	Description string
	Type        string
	FullType    string
	FieldType   string
	Annotations []string
}

func newField(fi *FieldInfo, kind string) (f *Field) {
	f = &Field{
		Name:        fi.Name,
		Modifier:    fi.Modifier,
		Order:       fi.Order,
		Description: fi.Description,
	}
	if fi.JavaAnnotations != "" {
		f.Annotations = strings.Split(fi.JavaAnnotations, ";")
	}

	if kind == "enum" {
		f.Type = fi.Type
		f.FieldType = "ENUM"
		if fi.Modifier == "repeated" {
			f.FullType = "List<" + f.Type + ">"
		} else {
			f.FullType = f.Type
		}
	} else {
		tm, ok := javaTypes[fi.Type]
		if ok {
			f.FieldType = strings.ToUpper(fi.Type)
			if fi.JavaType != "" {
				f.Type = fi.JavaType
			} else if fi.Modifier == "repeated" {
				f.Type = tm.ObjectType
			} else {
				f.Type = tm.PrimaryType
			}

			if fi.Modifier == "repeated" {
				f.FullType = "List<" + f.Type + ">"
			} else {
				f.FullType = f.Type
			}
		} else if strings.HasPrefix(fi.Type, "map[") {
			f.Type = fmt.Sprintf("Map<%s>", buildMapTypes(fi.Type[4:len(fi.Type)-1]))
			f.FieldType = "MESSAGE"
			f.FullType = f.Type
		} else {
			f.Type = fi.Type
			f.FieldType = "MESSAGE"
			if fi.Modifier == "repeated" {
				f.FullType = "List<" + f.Type + ">"
			} else {
				f.FullType = f.Type
			}
		}
	}
	return
}

func buildMapTypes(types string) string {
	arr := strings.Split(types, ",")
	if len(arr) != 2 {
		panic("invalid map types: " + types)
	}

	for i := 0; i < 2; i++ {
		arr[i] = strings.TrimSpace(arr[i])
		if tm, ok := javaTypes[arr[i]]; ok {
			arr[i] = tm.ObjectType
		}
	}

	return strings.Join(arr, ", ")
}

// Enum 枚举信息
type Enum struct {
	Name        string `xml:"name,attr"`
	Description string `xml:"description,attr"`
	Fields      []*struct {
		Name        string `xml:"name,attr"`
		Value       string `xml:"value,attr"`
		Description string `xml:"description,attr"`
	} `xml:"field"`
}

// Import 包引入定义
type Import struct {
	Path  string `xml:"path,attr"`
	Alias string `xml:"alias,attr"`
}

func newImport(path, alias string) *Import {
	return &Import{
		Path:  path,
		Alias: alias,
	}
}
