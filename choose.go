package bubblewrap

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Choose lets you prompt the user with a list of options
// and allows them to select one or more of them.
func Choose(options []string, optFns ...func(*chooser) error) ([]string, error) {
	c, err := NewChooser(options, optFns...)
	if err != nil {
		return []string{}, err
	}
	return c.choose()
}

type option struct {
	value    string
	selected bool
}

func newOptionsFromStrings(strings []string) []option {
	options := []option{}
	for _, s := range strings {
		options = append(options, option{value: s})
	}
	return options
}

// Chooser allows you to prompt a user with a list of options and let them select from them.
type chooser struct {
	options      []option
	currentIndex int
	stdin        io.Reader
	stdout       io.Writer
	ctx          context.Context
	aborted      bool
	paginator    paginator.Model
	cursor       string

	cursorStyle       lipgloss.Style
	itemStyle         lipgloss.Style
	selectedItemStyle lipgloss.Style
}

func NewChooser(options []string, opts ...func(*chooser) error) (*chooser, error) {
	c := &chooser{
		options:      newOptionsFromStrings(options),
		currentIndex: 0,
		stdin:        os.Stdin,
		stdout:       os.Stdout,
		paginator:    paginator.New(),
		cursor:       ">",
	}
	c.paginator.PerPage = 10

	for _, f := range opts {
		err := f(c)
		if err != nil {
			return c, err
		}
	}

	return c, nil
}

func (c *chooser) Init() tea.Cmd {
	return nil
}

func (c *chooser) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return c, nil
	case tea.KeyMsg:
		switch keyPressed := msg.String(); keyPressed {
		case " ":
			switch c.options[c.currentIndex].selected {
			case true:
				c.options[c.currentIndex].selected = false
			case false:
				c.options[c.currentIndex].selected = true
			}
		case "down":
			if c.currentIndex < len(c.options)-1 {
				c.currentIndex++
				if c.currentIndex >= (c.paginator.Page + 1) * c.paginator.PerPage {
					c.paginator.Page++
				}
			}
		case "up":
			if c.currentIndex > 0 {
				c.currentIndex--
				if c.currentIndex < (c.paginator.Page + 1) * c.paginator.PerPage {
					c.paginator.Page--
				}
			}
		case "enter":
			return c, tea.Quit
		case "ctrl+c", "esc":
			c.aborted = true
			return c, tea.Quit
		}
	}
	return c, nil
}

func (c *chooser) View() string {
	viewBuilder := strings.Builder{}

	first, last := c.paginator.GetSliceBounds(len(c.options))
	for i, opt := range c.options[first:last] {
		if i == c.currentIndex%c.paginator.PerPage {
			viewBuilder.WriteString(c.cursorStyle.Render(c.cursor))
		} else {
			viewBuilder.WriteString(strings.Repeat(" ", len(c.cursor)))
		}

		if opt.selected {
			viewBuilder.WriteString(c.itemStyle.Render("[X] "))
		} else {
			viewBuilder.WriteString(c.itemStyle.Render("[ ] "))
		}

		viewBuilder.WriteString(c.itemStyle.Render(opt.value))
		viewBuilder.WriteString("\n")
	}

	return viewBuilder.String()
}

func (c *chooser) selected() []string {
	selected := []string{}
	for _, o := range c.options {
		if o.selected {
			selected = append(selected, o.value)
		}
	}
	return selected
}

func (c *chooser) choose() ([]string, error) {
	p := tea.NewProgram(c, tea.WithContext(c.ctx), tea.WithInput(c.stdin), tea.WithOutput(c.stdout))
	if _, err := p.Run(); err != nil {
		return []string{}, err
	}
	if c.aborted {
		return []string{}, CancelError(fmt.Errorf("user cancelled operation"))
	}
	return c.selected(), nil
}
