package tpl

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/template"

	"github.com/cuigh/auxo/ext/texts"
	"github.com/gobuffalo/packr"
)

// app
const (
	// AppPomXML service/msg/task/web 模块 pom.xml 文件
	AppPomXML = "app:pom.xml"
	// BuildJSON 构建信息文件
	BuildJSON = "build.json"
	// AssemblyXML 打包配置文件
	AssemblyXML = "assembly.xml"
	// PomProperties 程序信息文件模板
	PomProperties = "pom.properties"
	// AppConfig service/msg/task/web 模块的 app.conf 文件
	AppConfig    = "app.conf"
	TestDao      = "TestDao.java"
	TestBiz      = "TestBiz.java"
	TestEntity   = "TestObject.java"
	AbstractTest = "AbstractTest.java"
)

// service
const (
	// ServiceBootstrap RPC 服务程序启动类模板
	TestServiceImpl          = "service:TestServiceImpl.java"
	TestServiceTest          = "service:TestServiceTest.java"
	ServiceTestAppProperties = "service:test:app.properties"
)

// msg
const (
	TestHandler     = "msg:TestHandler.java"
	TestHandlerTest = "msg:TestHandlerTest.java"
)

// task
const (
	TestTask      = "task:TestTask.java"
	TestTaskTest  = "task:TestTaskTest.java"
)

// web
const (
	TestController   = "web:TestController.java"
	WebAppConfig     = "web:app.conf"
	WebAppProperties = "web:app.properties"
)

// contract
const (
	// ContractPomXML RPC 服务契约模块 pom 文件
	ContractPomXML = "contract:pom.xml"
	// TestService 测试服务接口
	TestService = "TestService.java"
	// TestType 测试枚举
	TestType = "TestType.java"
	// TestDto 测试服务数据模型
	TestDto = "TestDto.java"
	// TestServiceXML 测试服务描述文件
	TestServiceXML = "TestService.xml"
	// SpringFactories SpringBoot 配置文件
	SpringFactories = "spring.factories"
	// ServiceAutoConfig 服务自动配置类
	ServiceAutoConfig = "ServiceAutoConfig.java"
)

var box packr.Box

func SetBox(b packr.Box) {
	box = b
}

// Execute 执行模板并输出到文件
func Execute(files map[string]string, data interface{}) error {
	for fn, tn := range files {
		if err := execute(fn, tn, data); err != nil {
			return err
		}
	}
	return nil
}

// ExecuteWriter 执行模板并输出到 Writer
func ExecuteWriter(w io.Writer, files map[string]string, data interface{}) error {
	for _, tn := range files {
		if err := executeWriter(w, tn, data); err != nil {
			return err
		}
	}
	return nil
}

func execute(fn, tn string, data interface{}) error {
	t := getTemplate(tn)

	// execute template
	buf := &bytes.Buffer{}
	err := t.Execute(buf, data)
	if err != nil {
		return fmt.Errorf("execute template [%v] failed: %v", tn, err)
	}

	// create dir
	dir := filepath.Dir(fn)
	if _, err = os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("create dir [%v] failed: %v", dir, err)
		}
	}

	// open file
	file, err := os.OpenFile(fn, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("create file [%s] failed: %v", fn, err)
	}

	// save to file
	// err = ioutil.WriteFile(path, buf.Bytes(), 0666)
	_, err = file.Write(buf.Bytes())
	if err != nil {
		return fmt.Errorf("create %s failed: %v", fn, err)
	}

	fmt.Println("C > " + fn)
	return nil
}

func executeWriter(w io.Writer, tn string, data interface{}) error {
	t := getTemplate(tn)

	// execute template
	err := t.Execute(w, data)
	if err != nil {
		return fmt.Errorf("execute template [%v] failed: %v", tn, err)
	}

	return nil
}

func getTemplate(tplName string) *template.Template {
	s := box.String(tplName)
	if s == "" {
		panic("template not found: " + tplName)
	}

	fm := template.FuncMap{
		"camel":  func(s string) string { return texts.Rename(s, texts.Camel) },
		"pascal": func(s string) string { return texts.Rename(s, texts.Pascal) },
		"upper":  func(s string) string { return texts.Rename(s, texts.Upper) },
		"lower":  func(s string) string { return texts.Rename(s, texts.Lower) },
	}
	return template.Must(template.New("T").Funcs(fm).Parse(s))
}

// type Container struct {
// 	dir string
// }

// func (c *Container) Add(parts ...string) {

// }
