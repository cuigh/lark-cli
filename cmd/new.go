package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cuigh/auxo/app"
	"github.com/cuigh/auxo/app/flag"
	"github.com/cuigh/auxo/config"
	"github.com/cuigh/auxo/errors"
	"github.com/cuigh/lark-cli/tpl"
	"github.com/cuigh/lark-cli/util/file"
	"github.com/cuigh/lark-cli/util/pom"
)

func New() *app.Command {
	cmd := app.NewCommand("new", "Create a project or module.", func(ctx *app.Context) error {
		fmt.Println("Usage: lark-cli new project|module")
		return nil
	})
	cmd.Flags.Register(flag.Help)
	cmd.AddCommand(NewProject())
	cmd.AddCommand(NewModule("service"))
	cmd.AddCommand(NewModule("task"))
	cmd.AddCommand(NewModule("web"))
	cmd.AddCommand(NewModule("msg"))
	cmd.AddCommand(NewContract())
	return cmd
}

// NewProject create `project` sub command
func NewProject() *app.Command {
	desc := "Create a project."
	cmd := app.NewCommand("project", desc, func(ctx *app.Context) error {
		args := &struct {
			Group    string `option:"group"`
			Artifact string `option:"artifact"`
		}{}
		if err := config.Unmarshal(args); err != nil {
			return err
		}

		wd, err := os.Getwd()
		if err != nil {
			return errors.Wrap(err, "acquire work directory failed")
		}

		if len(ctx.Args()) == 0 {
			return errors.New("project name is missing")
		}

		name := ctx.Args()[0]
		if args.Artifact == "" {
			args.Artifact = name
		}

		// check dir exist
		dir := filepath.Join(wd, name)
		if file.Exist(dir) {
			return errors.New("directory already exist: " + dir)
		}

		data := map[string]string{
			"GroupID":    args.Group,
			"ArtifactID": args.Artifact,
		}

		// create files
		files := make(map[string]string)
		files[filepath.Join(dir, "pom.xml")] = tpl.ProjectPomXML
		files[filepath.Join(dir, "README.md")] = tpl.ReadMe
		files[filepath.Join(dir, ".gitignore")] = tpl.GitIgnore
		if err = tpl.Execute(files, data); err != nil {
			return err
		}

		fmt.Println("finished.")
		return nil
	})
	cmd.Flags.Register(flag.Help)
	cmd.Flags.String("group", "g", "", "group id")
	cmd.Flags.String("artifact", "a", "", "artifact id")
	return cmd
}

// NewModule create `new service/msg/task/web` sub command
func NewModule(moduleType string) *app.Command {
	desc := fmt.Sprintf("Create a %s module.", moduleType)
	cmd := app.NewCommand(moduleType, desc, func(ctx *app.Context) error {
		args := &struct {
			Group    string `option:"group"`
			Artifact string `option:"artifact"`
			Package  string `option:"package"`
		}{}
		if err := config.Unmarshal(args); err != nil {
			return app.Fatal(1, err)
		}

		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("acquire work directory failed: %v", err)
		}

		// load parent pom
		p, err := pom.NewPom(filepath.Join(wd, "pom.xml"))
		if err != nil {
			return err
		}

		// check args
		var name string
		if len(ctx.Args()) == 0 {
			if p == nil {
				return fmt.Errorf("module name is missing")
			}
			name = fmt.Sprintf("%v-%v", p.GetArtifactID(), moduleType)
		} else {
			name = ctx.Args()[0]
		}

		// build template data
		if args.Group == "" && p != nil {
			args.Group = p.GetGroupID()
		}
		if args.Group == "" {
			return fmt.Errorf("group arg is missing")
		}
		if args.Artifact == "" {
			args.Artifact = name
		}
		if args.Package == "" {
			args.Package = args.Group + "." + strings.Replace(args.Artifact, "-", ".", -1)
		}
		data := map[string]string{
			"Type":       moduleType,
			"GroupID":    args.Group,
			"ArtifactID": args.Artifact,
			"Package":    args.Package,
		}

		// check dir exist
		moduleDir := filepath.Join(wd, name)
		_, err = os.Stat(moduleDir)
		if err == nil {
			return fmt.Errorf("directory already exist: %v", moduleDir)
		}

		// create empty dirs
		var dirs []string
		if moduleType == "web" {
			dirs = append(dirs, filepath.Join(wd, name, "src", "main", "resources", "view"))
			dirs = append(dirs, filepath.Join(wd, name, "src", "main", "resources", "static", "js"))
			dirs = append(dirs, filepath.Join(wd, name, "src", "main", "resources", "static", "css"))
		}
		dirs = append(dirs, filepath.Join(wd, name, "src", "test", "java"))
		file.CreateDir(dirs...)

		// create files
		files := make(map[string]string)
		files[filepath.Join(moduleDir, "pom.xml")] = tpl.AppPomXML
		files[filepath.Join(moduleDir, "build", "build.json")] = tpl.BuildJSON
		files[filepath.Join(moduleDir, "build", "assembly.xml")] = tpl.AssemblyXML
		files[filepath.Join(moduleDir, "build", "pom.properties")] = tpl.PomProperties
		if moduleType == "web" {
			files[filepath.Join(moduleDir, "src", "main", "resources", "etc", "app.conf")] = tpl.WebAppConfig
		} else {
			files[filepath.Join(moduleDir, "src", "main", "resources", "etc", "app.conf")] = tpl.AppConfig
		}
		files[file.NewPath(moduleDir, "src", "main", "java").Join(strings.Split(args.Package, ".")...).Join("Bootstrap.java").String()] = fmt.Sprintf("%v:Bootstrap.java", moduleType)
		files[file.NewPath(moduleDir, "src", "main", "java").Join(strings.Split(args.Package, ".")...).Join("dao", "TestDao.java").String()] = tpl.TestDao
		files[file.NewPath(moduleDir, "src", "main", "java").Join(strings.Split(args.Package, ".")...).Join("biz", "TestBiz.java").String()] = tpl.TestBiz
		files[file.NewPath(moduleDir, "src", "main", "java").Join(strings.Split(args.Package, ".")...).Join("entity", "TestObject.java").String()] = tpl.TestEntity
		files[file.NewPath(moduleDir, "src", "test", "java").Join(strings.Split(args.Package, ".")...).Join("AbstractTest.java").String()] = tpl.AbstractTest
		switch moduleType {
		case "service":
			files[file.NewPath(moduleDir, "src", "main", "java").Join(strings.Split(args.Package, ".")...).Join("impl", "TestServiceImpl.java").String()] = tpl.TestServiceImpl
			files[file.NewPath(moduleDir, "src", "test", "java").Join(strings.Split(args.Package, ".")...).Join("TestServiceTest.java").String()] = tpl.TestServiceTest
			files[filepath.Join(moduleDir, "src", "test", "resources", "etc", "app.properties")] = tpl.ServiceTestAppProperties
		case "msg":
			files[file.NewPath(moduleDir, "src", "main", "java").Join(strings.Split(args.Package, ".")...).Join("handler", "TestHandler.java").String()] = tpl.TestHandler
			files[file.NewPath(moduleDir, "src", "test", "java").Join(strings.Split(args.Package, ".")...).Join("handler", "TestHandlerTest.java").String()] = tpl.TestHandlerTest
		case "task":
			files[file.NewPath(moduleDir, "src", "main", "java").Join(strings.Split(args.Package, ".")...).Join("executor", "TestTask.java").String()] = tpl.TestTask
			files[file.NewPath(moduleDir, "src", "test", "java").Join(strings.Split(args.Package, ".")...).Join("executor", "TestTaskTest.java").String()] = tpl.TestTaskTest
		case "web":
			files[file.NewPath(moduleDir, "src", "main", "java").Join(strings.Split(args.Package, ".")...).Join("controller", "TestController.java").String()] = tpl.TestController
			files[filepath.Join(moduleDir, "src", "main", "resources", "etc", "app.properties")] = tpl.WebAppProperties
		}
		if err = tpl.Execute(files, data); err != nil {
			return err
		}

		// modify files
		if p != nil {
			p.AddModule(args.Artifact)
		}

		fmt.Println("finished.")
		return nil
	})
	cmd.Flags.Register(flag.Help)
	cmd.Flags.String("group", "g", "", "group id")
	cmd.Flags.String("artifact", "a", "", "artifact id")
	cmd.Flags.String("package", "p", "", "package")
	return cmd
}

// NewContract create `contract` sub command
func NewContract() *app.Command {
	desc := "Create a service contract module."
	cmd := app.NewCommand("contract", desc, func(ctx *app.Context) error {
		args := &struct {
			Group    string `option:"group"`
			Artifact string `option:"artifact"`
			Package  string `option:"package"`
		}{}
		if err := config.Unmarshal(args); err != nil {
			return app.Fatal(1, err)
		}

		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("acquire work directory failed: %v", err)
		}

		// load parent pom
		p, err := pom.NewPom(filepath.Join(wd, "pom.xml"))
		if err != nil {
			return err
		}

		// check args
		var name string
		if len(ctx.Args()) == 0 {
			if p == nil {
				return fmt.Errorf("module name is missing")
			}
			name = p.GetArtifactID() + "-service-contract"
		} else {
			name = ctx.Args()[0]
		}

		// build template data
		if args.Group == "" && p != nil {
			args.Group = p.GetGroupID()
		}
		if args.Group == "" {
			return fmt.Errorf("group arg is missing")
		}
		if args.Artifact == "" {
			args.Artifact = name
		}
		if args.Package == "" {
			args.Package = args.Group + "." + strings.Replace(strings.TrimSuffix(args.Artifact, "-contract"), "-", ".", -1)
		}
		data := map[string]string{
			"GroupID":    args.Group,
			"ArtifactID": args.Artifact,
			"Package":    args.Package,
		}

		// check dir exist
		moduleDir := filepath.Join(wd, name)
		_, err = os.Stat(moduleDir)
		if err == nil {
			return fmt.Errorf("directory already exist: %v", moduleDir)
		}

		// todo: create files
		files := make(map[string]string)
		files[filepath.Join(moduleDir, "pom.xml")] = tpl.ContractPomXML
		files[file.NewPath(moduleDir, "src", "main", "java").Join(strings.Split(args.Package, ".")...).Join("constant", "TestType.java").String()] = tpl.TestType
		files[file.NewPath(moduleDir, "src", "main", "java").Join(strings.Split(args.Package, ".")...).Join("dto", "TestDto.java").String()] = tpl.TestDto
		files[file.NewPath(moduleDir, "src", "main", "java").Join(strings.Split(args.Package, ".")...).Join("iface", "TestService.java").String()] = tpl.TestService
		files[file.NewPath(moduleDir, "src", "main", "java").Join(strings.Split(args.Package, ".")...).Join("spring", "ServiceAutoConfig.java").String()] = tpl.ServiceAutoConfig
		files[filepath.Join(moduleDir, "src", "main", "resources", "META-INF", "spring.factories")] = tpl.SpringFactories
		files[filepath.Join(moduleDir, "src", "main", "rsd", "TestService.xml")] = tpl.TestServiceXML
		if err = tpl.Execute(files, data); err != nil {
			return err
		}

		// modify files
		if p != nil {
			p.AddModule(args.Artifact)
		}

		fmt.Println("finished.")
		return nil
	})
	cmd.Flags.Register(flag.Help)
	cmd.Flags.String("group", "g", "", "group id")
	cmd.Flags.String("artifact", "a", "", "artifact id")
	cmd.Flags.String("package", "p", "", "package")
	return cmd
}
