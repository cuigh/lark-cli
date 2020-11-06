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
	"github.com/cuigh/lark/tpl"
	"github.com/cuigh/lark/util/file"
	"github.com/cuigh/lark/util/pom"
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

		if args.Group == "" {
			return errors.New("group is missing")
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
		files[filepath.Join(dir, "pom.xml")] = "project/pom.xml"
		files[filepath.Join(dir, "README.md")] = "project/README.md"
		files[filepath.Join(dir, ".gitignore")] = "project/gitignore"
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

		fp := func(name string) string {
			return fmt.Sprintf("modules/%s/%s", moduleType, name)
		}

		// create files
		files := make(map[string]string)
		files[filepath.Join(moduleDir, "pom.xml")] = fp("pom.xml")
		files[filepath.Join(moduleDir, "src", "main", "resources", "application.yml")] = fp("application.yml")
		files[file.NewPath(moduleDir, "src", "main", "java").Join(strings.Split(args.Package, ".")...).Join("Bootstrap.java").String()] = fp("Bootstrap.java")
		switch moduleType {
		case "service":
			files[file.NewPath(moduleDir, "src", "main", "java").Join(strings.Split(args.Package, ".")...).Join("dao", "TestDao.java").String()] = fp("TestDao.java")
			files[file.NewPath(moduleDir, "src", "main", "java").Join(strings.Split(args.Package, ".")...).Join("biz", "TestBiz.java").String()] = fp("TestBiz.java")
			files[file.NewPath(moduleDir, "src", "main", "java").Join(strings.Split(args.Package, ".")...).Join("entity", "TestObject.java").String()] = fp("TestObject.java")
			files[file.NewPath(moduleDir, "src", "main", "java").Join(strings.Split(args.Package, ".")...).Join("impl", "TestServiceImpl.java").String()] = fp("TestServiceImpl.java")
			files[file.NewPath(moduleDir, "src", "test", "java").Join(strings.Split(args.Package, ".")...).Join("TestServiceTests.java").String()] = fp("TestServiceTests.java")
		case "msg":
			files[file.NewPath(moduleDir, "src", "main", "java").Join(strings.Split(args.Package, ".")...).Join("handler", "TestHandler.java").String()] = fp("TestHandler.java")
			files[file.NewPath(moduleDir, "src", "test", "java").Join(strings.Split(args.Package, ".")...).Join("handler", "TestHandlerTests.java").String()] = fp("TestHandlerTests.java")
		case "task":
			files[file.NewPath(moduleDir, "src", "main", "java").Join(strings.Split(args.Package, ".")...).Join("executor", "TestExecutor.java").String()] = fp("TestExecutor.java")
			files[file.NewPath(moduleDir, "src", "test", "java").Join(strings.Split(args.Package, ".")...).Join("executor", "TestExecutorTests.java").String()] = fp("TestExecutorTests.java")
		case "web":
			files[file.NewPath(moduleDir, "src", "main", "java").Join(strings.Split(args.Package, ".")...).Join("controller", "TestController.java").String()] = fp("TestController.java")
			files[file.NewPath(moduleDir, "src", "test", "java").Join(strings.Split(args.Package, ".")...).Join("executor", "TestControllerTests.java").String()] = fp("TestControllerTests.java")
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

		files := make(map[string]string)
		files[filepath.Join(moduleDir, "pom.xml")] = "modules/contract/pom.xml"
		files[file.NewPath(moduleDir, "src", "main", "java").Join(strings.Split(args.Package, ".")...).Join("constant", "TestType.java").String()] = "modules/contract/TestType.java"
		files[file.NewPath(moduleDir, "src", "main", "java").Join(strings.Split(args.Package, ".")...).Join("dto", "TestDto.java").String()] = "modules/contract/TestDto.java"
		files[file.NewPath(moduleDir, "src", "main", "java").Join(strings.Split(args.Package, ".")...).Join("iface", "TestService.java").String()] = "modules/contract/TestService.java"
		files[file.NewPath(moduleDir, "src", "main", "java").Join(strings.Split(args.Package, ".")...).Join("config", "ProxyConfigurer.java").String()] = "modules/contract/ProxyConfigurer.java"
		files[filepath.Join(moduleDir, "src", "main", "resources", "META-INF", "spring.factories")] = "modules/contract/spring.factories"
		files[filepath.Join(moduleDir, "src", "main", "services", "TestService.xml")] = "modules/contract/TestService.xml"
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
