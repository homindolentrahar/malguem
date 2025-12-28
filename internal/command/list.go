package command

import (
	"fmt"
	"malguem/internal/config"
	"malguem/internal/util"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func list() *cobra.Command {
	var isGlobal bool

	var cmd = &cobra.Command{
		Use:   "list",
		Short: "List all featured templates",
		Run: func(cmd *cobra.Command, args []string) {
			var templates []string
			var err error
			if isGlobal {
				templates, err = getGlobalTemplates()
			} else {
				templates, err = getLocalTemplates()
			}

			if err != nil {
				fmt.Printf("failed to get featured templates: %s\n", err.Error())
				os.Exit(1)
			}

			if len(templates) < 1 {
				fmt.Println("🚫  No templates found")
				os.Exit(1)
			}

			fmt.Printf("\n🗂️  %d template(s) found:\n\n", len(templates))
			for index, item := range templates {
				fmt.Printf("   %d. %s\n", index+1, item)
			}
			fmt.Println()
		},
	}

	cmd.Flags().BoolVarP(&isGlobal, "global", "G", false, "List all global templates")

	return cmd
}

func getLocalTemplates() ([]string, error) {
	malguem, err := config.ReadMalguem()
	if err != nil {
		return nil, err
	}

	var templates []string
	for name, _ := range malguem.Templates {
		templates = append(templates, name)
	}

	return templates, nil
}

func getGlobalTemplates() ([]string, error) {
	cacheDir, err := util.CacheDir()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(cacheDir)
	if err != nil {
		return nil, err
	}

	if len(entries) < 1 {
		return nil, fmt.Errorf("🚫  No templates found\n")
	}

	var templates []string
	for _, entry := range entries {
		if entry.IsDir() == false {
			continue
		}

		parts := strings.Split(entry.Name(), "-")
		item := strings.Join(parts[:len(parts)-1], "-")
		templates = append(templates, item)
	}

	return templates, nil
}
