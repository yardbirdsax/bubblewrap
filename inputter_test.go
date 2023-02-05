package bubblewrap

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

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

func TestInput(t *testing.T) {
	t.Run("with input", func(t *testing.T) {
		i := NewInputter()
		expectedPrompt := "prompt"
		expectedValue := "expected"

		testContext, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelFunc()
		i.ctx = testContext
		in := bytes.NewBuffer([]byte{})
		in.WriteString(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(expectedValue + "\r")}.String())
		i.stdin = in
		out := bytes.NewBuffer([]byte{})
		i.stdout = out

		value, err := i.input(expectedPrompt)

		assert.Equal(t, expectedValue, value, "return value is not as expected")
		assert.Equal(t, expectedPrompt+expectedValue+"\x1b[7m \x1b[0m", i.View(), "prompt is not as expected")
		assert.NoError(t, err, "unexpected error returned")
	})

	t.Run("with cancel", func(t *testing.T) {
		i := NewInputter()
    in := bytes.NewBuffer([]byte{})
    in.WriteString("nothing\u001B")
    i.stdin = in
    expectedError := CancelError(fmt.Errorf("an error"))

    value, err := i.input("prompt")

    assert.Equal(t, "", value, "value was not empty string")
    assert.True(t, i.aborted, "aborted was not set properly")
    assert.IsType(t, expectedError, err, "error returned was not of expected type")
	})

}
