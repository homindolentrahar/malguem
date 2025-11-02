package command

import (
	"fmt"
	"malguem/internal/model"
	"malguem/internal/remote"
	"os"
	"sync"

	"github.com/spf13/cobra"
	"go.yaml.in/yaml/v3"
)

var get = &cobra.Command{
	Use:   "get",
	Short: "Fetch template from remote source",
	Run: func(cmd *cobra.Command, args []string) {
		// Read malguem.yaml
		malguem, err := readMalguem()
		if err != nil {
			fmt.Printf("Failed to get templates: %v\n", err)
			os.Exit(1)
		}

		var templateUrls []string
		for _, info := range malguem.Templates {
			templateUrls = append(templateUrls, info.Github.Url)
		}

		// Define the wait group
		var wg sync.WaitGroup
		worker := 10 // Maximal worker unit
		// Channel to async-communicate between goroutines
		jobs := make(chan string, len(templateUrls))

		for range worker {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for url := range jobs {
					// Clone repo
					path, err := remote.Clone(url)
					if err != nil || path == "" {
						fmt.Printf("⚠️  Failed to get template from: %s\n", url)
						continue
					}

					fmt.Printf("🗂️  Template location: %s\n", path)
				}
			}()
		}

		// Feed source into jobs channel
		for _, source := range templateUrls {
			jobs <- source
		}
		close(jobs)

		wg.Wait()
	},
}

func readMalguem() (*model.Malguem, error) {
	file, err := os.ReadFile("malguem.yaml")
	if err != nil {
		return nil, err
	}

	var malguem model.Malguem
	err = yaml.Unmarshal(file, &malguem)
	if err != nil {
		return nil, err
	}

	return &malguem, nil
}
