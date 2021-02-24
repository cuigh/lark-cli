package tpl

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/template"

	"github.com/cuigh/auxo/ext/texts"
)

//go:embed assets
var fs embed.FS

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
	b, err := fs.ReadFile("assets/" + tplName)
	if err != nil {
		panic(fmt.Sprintf("load template [%s] failed: %s", tplName, err))
	}

	fm := template.FuncMap{
		"camel":  func(s string) string { return texts.Rename(s, texts.Camel) },
		"pascal": func(s string) string { return texts.Rename(s, texts.Pascal) },
		"upper":  func(s string) string { return texts.Rename(s, texts.Upper) },
		"lower":  func(s string) string { return texts.Rename(s, texts.Lower) },
	}
	return template.Must(template.New("T").Funcs(fm).Parse(string(b)))
}
