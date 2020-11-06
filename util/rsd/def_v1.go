package rsd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/cuigh/auxo/app"
)

// definitionV1 V1 版本服务定义
type definitionV1 struct {
	JavaPackage string `xml:"javaPackage,attr"`
	Version     string `xml:"version,attr"`
	Service     *struct {
		Name        string    `xml:"name,attr"`
		Alias       string    `xml:"alias,attr"`
		Fail        string    `xml:"fail,attr"`
		Description string    `xml:"description,attr"`
		Imports     []*Import `xml:"imports>import"`
		Methods     []*struct {
			Name        string `xml:"name,attr"`
			Alias       string `xml:"alias,attr"`
			Invoke      string `xml:"invoke,attr"`
			Description string `xml:"description,attr"`
			Request     *struct {
				Multiple bool           `xml:"multiple,attr"`
				Extends  string         `xml:"extends,attr"`
				Fields   []*fieldInfoV1 `xml:"field"`
			} `xml:"request"`
			Response *struct {
				Multiple bool           `xml:"multiple,attr"`
				Extends  string         `xml:"extends,attr"`
				Fields   []*fieldInfoV1 `xml:"field"`
			} `xml:"response"`
			Errors []*Error `xml:"errors>error"`
		} `xml:"method"`
	} `xml:"service"`
	ProtoImports []*Import `xml:"types>protoImports>protoImport"`
	Imports      []*Import `xml:"types>imports>import"`
	Types        []*struct {
		Name        string         `xml:"name,attr"`
		Description string         `xml:"description,attr"`
		Extends     string         `xml:"extends,attr"`
		Fields      []*fieldInfoV1 `xml:"field"`
	} `xml:"types>type"`
	Enums  []*Enum  `xml:"enums>enum"`
	Errors []*Error `xml:"errors>error"`
}

// fieldInfoV1 字段定义
type fieldInfoV1 struct {
	*FieldInfo
	Kind     string `xml:"kind,attr"`
	Inherits bool   `xml:"inherits,attr"`
}

func (d *definitionV1) GetPackage() string {
	return d.JavaPackage
}

func (d *definitionV1) GetEnumModels() (enums []*EnumModel) {
	for _, e := range d.Enums {
		em := &EnumModel{
			ToolVersion: app.Version,
			Package:     d.JavaPackage,
			Enum:        e,
		}
		enums = append(enums, em)
	}
	return
}

func (d *definitionV1) GetServiceModel() *ServiceModel {
	if d.Service == nil || len(d.Service.Methods) == 0 {
		return nil
	}

	s := &Service{
		Name:        d.Service.Name,
		Alias:       d.Service.Alias,
		Fail:        d.Service.Fail,
		Description: d.Service.Description,
		Imports:     d.Service.Imports,
		Methods:     make([]*Method, len(d.Service.Methods)),
	}
	for i, mi := range d.Service.Methods {
		m := &Method{
			Name:        mi.Name,
			Alias:       mi.Alias,
			Invoke:      mi.Invoke,
			Description: mi.Description,
		}
		if mi.Request != nil {
			if mi.Request.Multiple {
				for _, f := range mi.Request.Fields {
					name := strings.ToLower(f.Name[:1]) + f.Name[1:]
					pi := newParameter(d.getParamType(f), name)
					m.Request = append(m.Request, pi)
				}
			} else {
				m.Request = []*Parameter{newParameter(mi.Name+"Request", "request")}
			}
		}
		if mi.Response != nil {
			if mi.Response.Multiple {
				m.Response = newParameter(d.getParamType(mi.Response.Fields[0]), "")
			} else {
				m.Response = newParameter(mi.Name+"Response", "")
			}
		}
		s.Methods[i] = m
	}

	return &ServiceModel{
		ToolVersion: app.Version,
		Package:     d.JavaPackage,
		Service:     s,
	}
}

func (d *definitionV1) GetDtoModel() *DtoModel {
	types := make([]*Type, 0)
	errTypes := make([]*ErrorType, 0)
	for _, ti := range d.Types {
		t := newType(ti.Name, ti.Description, ti.Extends, d.createFields(ti.Fields))
		types = append(types, t)
	}
	if d.Service != nil {
		for _, m := range d.Service.Methods {
			if m.Request != nil && !m.Request.Multiple {
				t := newType(m.Name+"Request", m.Name+" 请求参数", m.Request.Extends, d.createFields(m.Request.Fields))
				types = append(types, t)
			}
			if m.Response != nil && !m.Response.Multiple {
				t := newType(m.Name+"Response", m.Name+" 响应结果", m.Response.Extends, d.createFields(m.Response.Fields))
				types = append(types, t)
			}
			if len(m.Errors) > 0 {
				et := &ErrorType{
					Name:        m.Name + "Errors",
					Description: m.Description + "错误码",
					Errors:      m.Errors,
				}
				errTypes = append(errTypes, et)
			}
		}
	}

	if len(types) > 0 || len(errTypes) > 0 {
		return newDtoModel(d.JavaPackage, d.Service.Name, d.Imports, types, errTypes, d.Enums)
	}
	return nil
}

func (d *definitionV1) validate() error {
	if d.JavaPackage == "" {
		return errors.New("javaPackage not set")
	}
	// 验证普通类型
	for _, t := range d.Types {
		if err := d.validateFields(t.Name, t.Fields); err != nil {
			return err
		}
	}
	// 验证方法参数和返回值
	if d.Service != nil {
		for _, m := range d.Service.Methods {
			if m.Request != nil && !m.Request.Multiple {
				if err := d.validateFields(m.Name+"Request", m.Request.Fields); err != nil {
					return err
				}
			}
			if m.Response != nil && !m.Response.Multiple {
				if err := d.validateFields(m.Name+"Response", m.Response.Fields); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (d *definitionV1) validateFields(typeName string, fields []*fieldInfoV1) error {
	m := map[string]string{}
	for _, f := range fields {
		if name, ok := m[f.Order]; ok {
			return fmt.Errorf("类型 %s 的字段 [%s, %s] 序号重复", typeName, f.Name, name)
		}
		m[f.Order] = f.Name
	}
	return nil
}

func (d *definitionV1) createFields(infos []*fieldInfoV1) (fields []*Field) {
	for _, fi := range infos {
		if !fi.Inherits {
			f := newField(fi.FieldInfo, fi.Kind)
			fields = append(fields, f)
		}
	}
	return
}

func (d *definitionV1) getParamType(f *fieldInfoV1) string {
	var t string
	if f.JavaType == "" {
		switch f.Type {
		case "bytes":
			t = "byte[]"
		case "string":
			t = "String"
		case "int32":
			t = "int"
		case "int64":
			t = "long"
		case "bool":
			t = "boolean"
		case "float":
			t = "float"
		case "double":
			t = "double"
		default:
			t = f.Type
		}
	} else {
		t = f.JavaType
	}
	if f.Modifier == "repeated" {
		t += "[]"
	}
	return t
}
