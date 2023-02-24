package bubblewrap

import (
	"bytes"
	"context"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChooseUpdate(t *testing.T) {
	testCases := []struct {
		name     string
		choices  []string
		expected []string
		inputs   []tea.KeyMsg
		aborted  bool
		cmds     []tea.Cmd
	}{
		{
			name:     "select one",
			choices:  []string{"one", "two", "three"},
			expected: []string{"two"},
			inputs:   []tea.KeyMsg{{Type: tea.KeyDown}, {Type: tea.KeyRunes, Runes: []rune(" ")}, {Type: tea.KeyEnter}},
			cmds:     []tea.Cmd{nil, nil, tea.Quit},
		},
		{
			name:     "select two",
			choices:  []string{"one", "two", "three"},
			expected: []string{"two", "three"},
			inputs:   []tea.KeyMsg{{Type: tea.KeyDown}, {Type: tea.KeyRunes, Runes: []rune(" ")}, {Type: tea.KeyDown}, {Type: tea.KeyRunes, Runes: []rune(" ")}, {Type: tea.KeyEnter}},
			cmds:     []tea.Cmd{nil, nil, nil, nil, tea.Quit},
		},
		{
			name:     "cancelled with ctrl+c",
			choices:  []string{"one", "two", "three"},
			expected: []string{},
			inputs:   []tea.KeyMsg{{Type: tea.KeyCtrlC}},
			aborted:  true,
			cmds:     []tea.Cmd{tea.Quit},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			c, err := NewChooser(tc.choices)
			require.NoError(t, err, "NewChooser returned unexpected error")
			testContext, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancelFunc()
			in := bytes.NewBuffer([]byte{})
			out := bytes.NewBuffer([]byte{})
			c.stdin = in
			c.stdout = out
			c.ctx = testContext
			require.NoError(t, err, "could not write expected input")
			actualCmds := []tea.Cmd{}

			for _, m := range tc.inputs {
				_, cmd := c.Update(m)
				actualCmds = append(actualCmds, cmd)
			}
			actual := c.selected()

			assert.Equal(t, tc.expected, actual)
			assert.NoError(t, err, "choose returned unexpected error")
			assert.Equal(t, tc.aborted, c.aborted, "aborted was not the expected value")
			assert.True(t, c.quitting, "quitting was not set")
			assert.Len(t, actualCmds, len(tc.cmds), "number of actual and expected commands not the same")
			// This little snippet is necessary because you can't simply assert that
			// two funcs are equal to each other.
			for i, av := range actualCmds {
				if i-1 > len(tc.cmds) {
					t.Errorf("actual command at position %d (%v) is outside of range of expected commands", i, av)
					continue
				}
				ev := tc.cmds[i]
				if ev == nil && av == nil {
					continue
				}
				if (ev == nil && av != nil) || (ev != nil && av == nil) {
					t.Errorf("one command is nil but other is not: actual (%v), expected (%v)", av, ev)
					continue
				}
				actualVal := av()
				expectedVal := ev()
				if expectedVal != actualVal {
					t.Errorf("result for actual value (%v) does not match result for expected value (%v)", actualVal, expectedVal)
					continue
				}
			}
		})
	}
}

func TestChooseView(t *testing.T) {
	testCases := []struct {
		name         string
		initialModel func() *chooser
		msgs         []tea.Msg
		expectedView string
	}{
		{
			name: "basic list",
			initialModel: func() *chooser {
				c, _ := NewChooser([]string{"one", "two", "three"})
				return c
			},
			expectedView: `>[ ] one
 [ ] two
 [ ] three
`,
		},
		{
			name: "basic list with selected",
			initialModel: func() *chooser {
				c, _ := NewChooser([]string{"one", "two", "three"})
				c.options[0].selected = true
				return c
			},
			expectedView: `>[X] one
 [ ] two
 [ ] three
`,
		},
		{
			name: "basic list with un-selected",
			initialModel: func() *chooser {
				c, _ := NewChooser([]string{"one", "two", "three"})
				return c
			},
			msgs: []tea.Msg{
				tea.KeyMsg{Type: tea.KeySpace},
				tea.KeyMsg{Type: tea.KeySpace},
			},
			expectedView: `>[ ] one
 [ ] two
 [ ] three
`,
		},
		{
			name: "paged list",
			initialModel: func() *chooser {
				c, _ := NewChooser([]string{"one", "two", "three", "four"})
				c.paginator.PerPage = 2
				return c
			},
			msgs: []tea.Msg{
				tea.KeyMsg{Type: tea.KeyDown},
				tea.KeyMsg{Type: tea.KeyDown},
			},
			expectedView: ">[ ] three\n [ ] four\n",
		},
		{
			name: "paged list backwards",
			initialModel: func() *chooser {
				c, _ := NewChooser([]string{"one", "two", "three", "four"})
				c.paginator.PerPage = 2
				return c
			},
			msgs: []tea.Msg{
				tea.KeyMsg{Type: tea.KeyDown},
				tea.KeyMsg{Type: tea.KeyDown},
				tea.KeyMsg{Type: tea.KeyUp},
			},
			expectedView: " [ ] one\n>[ ] two\n",
		},
		{
			name: "paged list backwards with selected",
			initialModel: func() *chooser {
				c, _ := NewChooser([]string{"one", "two", "three", "four"})
				c.paginator.PerPage = 2
				return c
			},
			msgs: []tea.Msg{
				tea.KeyMsg{Type: tea.KeyDown},
				tea.KeyMsg{Type: tea.KeySpace},
				tea.KeyMsg{Type: tea.KeyDown},
				tea.KeyMsg{Type: tea.KeySpace},
				tea.KeyMsg{Type: tea.KeyUp},
			},
			expectedView: " [ ] one\n>[X] two\n",
		},
		{
			name: "paged list backwards with list smaller than page size",
			initialModel: func() *chooser {
				c, _ := NewChooser([]string{"one", "two", "three", "four"})
				c.paginator.PerPage = 10
				return c
			},
			msgs: []tea.Msg{
				tea.KeyMsg{Type: tea.KeyDown},
				tea.KeyMsg{Type: tea.KeyDown},
				tea.KeyMsg{Type: tea.KeyUp},
			},
			expectedView: " [ ] one\n>[ ] two\n [ ] three\n [ ] four\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := tc.initialModel()
			for _, msg := range tc.msgs {
				m.Update(msg)
			}
			actualView := m.View()
			assert.Equal(t, tc.expectedView, actualView, "expected view does not match actual view")
		})
	}
}

func TestWithStyles(t *testing.T) {
	expectedCursorStyle := lipgloss.NewStyle().Bold(true)
	expectedItemStyle := lipgloss.NewStyle().Underline(true)
	expectedSelectedItemStyle := lipgloss.NewStyle().Italic(true)

	c, err := NewChooser(
		[]string{"one"},
		WithCursorStyle(expectedCursorStyle),
		WithItemStyle(expectedItemStyle),
		WithSelectedItemStyle(expectedSelectedItemStyle),
	)

	assert.NoError(t, err, "NewChooser returned unexpected error")
	assert.EqualValues(t, expectedCursorStyle, c.cursorStyle, "cursor style did not match")
	assert.EqualValues(t, expectedItemStyle, c.itemStyle, "item style did not match")
	assert.EqualValues(t, expectedSelectedItemStyle, c.selectedItemStyle, "selected item style did not match")
}
