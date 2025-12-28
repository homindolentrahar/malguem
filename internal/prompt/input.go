package prompt

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type Input struct {
	Prompt string
	Input  string
	IsDone bool
}

func (input Input) Init() tea.Cmd {
	return nil
}

func (input Input) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			input.IsDone = true
			return input, tea.Quit
		case "backspace":
			if len(input.Input) > 0 {
				input.Input = input.Input[:len(input.Input)-1]
			}
		case "ctrl+c", ":q":
			return input, tea.Quit
		default:
			input.Input += msg.String()
		}
	}

	return input, nil
}

func (input Input) View() string {
	if input.IsDone {
		return fmt.Sprintf("%s: %s\n", input.Prompt, input.Input)
	}

	return fmt.Sprintf("%s: %s", input.Prompt, input.Input)
}
