package command

import (
	"fmt"
	"malguem/internal/config"
	"malguem/internal/remote"
	"malguem/internal/util"
	"os"
	"sync"

	"github.com/spf13/cobra"
)

var get = &cobra.Command{
	Use:   "get",
	Short: "Fetch template from remote source",
	Run: func(cmd *cobra.Command, args []string) {
		// Read malguem.yaml
		malguem, err := config.ReadMalguem()
		if err != nil {
			fmt.Printf("Failed to get templates: %v\n", err)
			os.Exit(1)
		}

		var templateUrls []string
		templateOutputs := make(map[string]string)
		for _, info := range malguem.Templates {
			url := info.Github.Url
			templateUrls = append(templateUrls, url)
			templateOutputs[url] = info.Output
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
					// Copy fetched template into output directory
					output := templateOutputs[url]
					copyTemplateIntoOutput(output, path)
				}
			}()
		}

		// Feed source into jobs channel
		for _, source := range templateUrls {
			jobs <- source
		}
		close(jobs)

		wg.Wait()

		fmt.Printf("\n✅  All templates updated. Happy coding!\n")
	},
}

func copyTemplateIntoOutput(outputPath, cachePath string) {
	// Make sure the output path exists
	os.MkdirAll(outputPath, os.ModePerm)

	// Copy template from cachePath into outputPath
	util.CopyDir(cachePath, outputPath)
}
