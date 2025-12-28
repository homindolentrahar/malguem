package command

import (
	"fmt"
	"malguem/internal/config"
	"malguem/internal/remote"
	"malguem/internal/util"
	"os"

	"github.com/spf13/cobra"
)

var makeCommand = &cobra.Command{
	Use:   "make",
	Short: "Generate code from template",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var templateName string

		templateName = args[0]

		// Read malguem.yaml config file
		malguemConfig, err := config.ReadMalguem()
		if err != nil {
			fmt.Println("Malguem config is not found")
			os.Exit(1)
		}

		if _, ok := malguemConfig.Templates[templateName]; !ok {
			fmt.Printf("🌧️  `%s` template not found\n", templateName)
			os.Exit(1)
		}

		// Read template config file on each template
		template := malguemConfig.Templates[templateName]

		// Read from output `template` directory
		templatePath := template.Path
		if _, err := os.Stat(templatePath); err != nil {
			// If not found in the output `template` directory, then read form cache
			templateUrl := template.Remote.Url
			cachePath, err := remote.CachePath(templateUrl)
			if err != nil {
				fmt.Printf("Error: %v", err)
				os.Exit(1)
			}

			// Check if the cache path exists or not
			if _, err := os.Stat(cachePath); err != nil {
				// Clone the repo first
				_, err := remote.Clone(templateUrl)
				if err != nil {
					fmt.Printf("Error: %v", err)
					os.Exit(1)
				}
			}

			templatePath = cachePath
		}

		// Generate code from template
		err = util.RenderTemplate(templatePath, template.Output)
		if err != nil {
			fmt.Printf("Error: %v", err)
			os.Exit(1)
		}
	},
}
