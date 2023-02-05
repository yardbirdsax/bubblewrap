package bubblewrap

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Input prompts for, well, input!
func Input(prompt string) (string, error) {
  input, err := NewInputter().input(prompt)
  if err != nil {
    return "", err
  }
  return input, nil
}

type inputter struct {
  ctx context.Context
  modelbase
  textinput textinput.Model
  stdin io.Reader
  stdout io.Writer
}

func NewInputter() *inputter {
  return &inputter{
    textinput: textinput.New(),
    ctx: context.TODO(),
    stdin: os.Stdin,
    stdout: os.Stdout,
  }
}

func (i *inputter) Init() tea.Cmd { return textinput.Blink }

func (i *inputter) View() string {
  return i.textinput.View()
}

func (i *inputter) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  switch msg := msg.(type) {
  case tea.KeyMsg:
    switch msg.Type {
    case tea.KeyCtrlC, tea.KeyEscape:
      i.quitting = true
      i.aborted = true
      return i, tea.Quit
    case tea.KeyEnter:
      i.quitting = true
      return i, tea.Quit
    }
  }

  var cmd tea.Cmd
  i.textinput, cmd = i.textinput.Update(msg)
  return i, cmd
}

func (i *inputter) input (prompt string) (string, error) {

  i.textinput.Prompt = prompt
  i.textinput.Focus()

  p := tea.NewProgram(i, tea.WithContext(i.ctx), tea.WithInput(i.stdin), tea.WithOutput(i.stdout))
  if _, err := p.Run(); err != nil {
    return "", err
  }
  if i.aborted {
    return "", CancelError(fmt.Errorf("user canceled operation"))
  }

  return i.textinput.Value(), nil
}
