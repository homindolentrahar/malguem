package command

import (
	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "malguem",
	Short: "Generate boilerplate code with your own template",
	Long:  "Tools for generating boilerplate code for mulitple language using your own template",
}

func Run() {
	root.CompletionOptions = cobra.CompletionOptions{
		HiddenDefaultCmd: true,
	}
	root.AddCommand(initialize)
	root.Execute()
}
