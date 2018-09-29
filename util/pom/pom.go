package pom

import (
	"fmt"

	"github.com/beevik/etree"
	"github.com/cuigh/lark/util/file"
)

// Pom Maven pom 文件
type Pom struct {
	filename string
	doc      *etree.Document
}

// NewPom 创建 Pom 对象
func NewPom(filename string) (*Pom, error) {
	if file.NotExist(filename) {
		return nil, nil
	}

	doc := etree.NewDocument()
	if err := doc.ReadFromFile(filename); err != nil {
		return nil, err
	}
	doc.Indent(4)

	return &Pom{
		filename: filename,
		doc:      doc,
	}, nil
}

// GetGroupID 获取 groupId 属性
func (p *Pom) GetGroupID() string {
	elem := p.doc.FindElement("project/groupId")
	return elem.Text()
}

// GetArtifactID 获取 artifactId 属性
func (p *Pom) GetArtifactID() string {
	elem := p.doc.FindElement("project/artifactId")
	return elem.Text()
}

// AddModule 添加模块
func (p *Pom) AddModule(name string) {
	elem := p.doc.FindElement("project/modules")
	elem.CreateElement("module").SetText(name)
	if err := p.doc.WriteToFile(p.filename); err != nil {
		panic(err)
	}
	fmt.Println("M > " + p.filename)
}
