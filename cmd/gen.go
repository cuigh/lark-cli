package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cuigh/auxo/app"
	"github.com/cuigh/auxo/app/flag"
	"github.com/cuigh/lark/tpl"
	"github.com/cuigh/lark/util/rsd"
)

// Gen create `gen` sub command
func Gen() *app.Command {
	cmd := app.NewCommand("gen", "Generate codes.", func(ctx *app.Context) error {
		fmt.Println("Usage: lark-cli gen rpc <NAME>")
		return nil
	})
	cmd.Flags.Register(flag.Help)
	cmd.AddCommand(GenRPC())
	return cmd
}

// GenRPC generate RPC service contract files from service definition file
func GenRPC() *app.Command {
	cmd := app.NewCommand("rpc", "Generate RPC service contract files for rpc contract project.", func(ctx *app.Context) error {
		filenames := ctx.Args()
		if len(filenames) == 0 {
			return fmt.Errorf("service definition files is missing")
		}

		defs := make(map[string]rsd.Definition)
		for _, filename := range filenames {
			if !strings.EqualFold(filepath.Ext(filename), ".xml") {
				return fmt.Errorf("not a valid definition file: %s", filename)
			}

			d, err := rsd.Load(filename)
			if err != nil {
				return err
			}
			defs[filename] = d
		}

		for filename, d := range defs {
			dir := filepath.Dir(filepath.Dir(filename))
			dir = filepath.Join(dir, "java", strings.Replace(d.GetPackage(), ".", string(filepath.Separator), -1))

			// service
			if sm := d.GetServiceModel(); sm != nil {
				files := make(map[string]string)
				files[filepath.Join(dir, "iface", sm.Name+"Service.java")] = "rpc/Service.java"
				err := tpl.Execute(files, sm)
				if err != nil {
					return err
				}
			}

			// enum
			enums := d.GetEnumModels()
			for _, em := range enums {
				files := make(map[string]string)
				files[filepath.Join(dir, "constant", em.Name+".java")] = "rpc/Enum.java"

				// err = tpl.ExecuteWriter(os.Stdout, files, e)
				err := tpl.Execute(files, em)
				if err != nil {
					return err
				}
			}

			// dto
			if dm := d.GetDtoModel(); dm != nil {
				files := make(map[string]string)
				files[filepath.Join(dir, "dto", dm.Name+"Dto.java")] = "rpc/Dto.java"
				err := tpl.Execute(files, dm)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
	cmd.Flags.Register(flag.Help)
	return cmd
}
