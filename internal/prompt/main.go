package prompt

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func PromptInput(promptMsg string) string {
	p := tea.NewProgram(Input{Prompt: promptMsg})
	input, err := p.Run()
	if err != nil {
		fmt.Printf("Prompt failed: %v\n", err)
		os.Exit(1)
	}

	inputValue := input.(Input).Input
	if inputValue == "" {
		return ""
	}

	return inputValue
}
