package bubblewrap

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Input prompts for, well, input!
func Input(prompt string) (string, error) {
  var input string

  return input, nil
}

type inputter struct {
  modelbase
  textinput textinput.Model
}

func (i *inputter) Init() tea.Cmd { return textinput.Blink }

func (i *inputter) View() string {
  if i.quitting {
    return ""
  }

  return i.textinput.View()
}

func (i *inputter) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  switch msg := msg.(type) {
  case tea.KeyMsg:
    switch msg.String() {
    case "ctrl+c", "esc":
      i.quitting = true
      i.aborted = true
      return i, tea.Quit
    case "enter":
      i.quitting = true
      return i, tea.Quit
    }
  }

  var cmd tea.Cmd
  i.textinput, cmd = i.textinput.Update(msg)
  return i, cmd
}

func (i *inputter) input (prompt string) (string, error) {

  i.textinput = textinput.New()
  i.textinput.Prompt = prompt

  return i.textinput.Value(), nil
}
