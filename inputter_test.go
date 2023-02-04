package bubblewrap

import (
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestInputterUpdate(t *testing.T) {
    i := &inputter{
    textinput: textinput.New(),
  }
  i.textinput.Focus()

  t.Run("input text", func(t *testing.T) {
    expectedInput := "this is input"
    keyMsg := tea.KeyMsg{
      Runes: []rune(expectedInput),
    }

    _, _ = i.Update(keyMsg)

    assert.Equal(t, expectedInput, i.textinput.Value(), "returned value is not expected")
  })

  t.Run("enter", func(t *testing.T) {
    keyMsg := tea.KeyMsg{
      Type: tea.KeyEnter,
    }
    _, cmd := i.Update(keyMsg)

    assert.Equal(t, tea.Quit(), cmd(), "did not return the Quite command")
    assert.True(t, i.quitting, "quitting not set properly")
    assert.False(t, i.aborted, "aborted not set properly")
  })


}