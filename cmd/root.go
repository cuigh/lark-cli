package cmd

import (
	"fmt"

	"github.com/cuigh/auxo/app"
)

// Root creates root command
func Root(ctx *app.Context) error {
	fmt.Println("Try `lark-cli -h` for more information")
	return nil
}
