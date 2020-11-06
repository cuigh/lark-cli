package rsd

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
)

// Load 解析服务定义配置
func Load(filename string) (Definition, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	header := new(headerInfo)
	err = xml.Unmarshal(bytes, header)
	if err != nil {
		return nil, fmt.Errorf("parse header failed: %s", err)
	}

	var d Definition
	if header.Version == "1.0" {
		d = new(definitionV1)
	} else {
		return nil, fmt.Errorf("unknown file version: %s", header.Version)
	}

	err = xml.Unmarshal(bytes, d)
	if err != nil {
		return nil, err
	}

	err = d.validate()
	if err != nil {
		return nil, err
	}

	return d, nil
}

// LoadAll 解析目录中所有服务定义配置
func LoadAll(dir string) ([]Definition, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	list := make([]Definition, 0)
	for _, f := range files {
		if f.IsDir() {
			continue
		}

		d, err := Load(f.Name())
		if err != nil {
			return nil, err
		}

		list = append(list, d)
	}
	return list, nil
}

type headerInfo struct {
	Version string `xml:"version,attr"`
}

// FieldInfo 基本字段定义
type FieldInfo struct {
	Name            string `xml:"name,attr"`
	Modifier        string `xml:"modifier,attr"`
	Type            string `xml:"type,attr"`
	Order           string `xml:"order,attr"`
	JavaType        string `xml:"javaType,attr"`
	Description     string `xml:"description,attr"`
	JavaAnnotations string `xml:"javaAnnotations,attr"`
}

// Definition 服务定义接口
type Definition interface {
	validate() error
	GetPackage() string
	GetEnumModels() (enums []*EnumModel)
	GetServiceModel() *ServiceModel
	GetDtoModel() *DtoModel
}
