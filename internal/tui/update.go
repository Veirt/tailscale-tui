package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/veirt/tailscale-tui/internal/tailscale"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.inputting && m.choice.IsBooleanFlag {
		return updateBooleanChoice(msg, m)
	}

	if m.inputting && !m.choice.IsBooleanFlag {
		return updateInput(msg, m)
	}

	return updateListSelection(msg, m)

}

func updateListSelection(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(tailscale.Flag)
			if ok {
				m.choice = i

				m.inputting = true
				if m.choice.IsBooleanFlag {
					m.answer = "true" // initialize to true
				}

			}

			return m, nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func updateInput(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			m.inputting = false
			return m, nil

		case "enter":
			m.inputting = false
			tailscale.FinalCmd.Flags[m.choice] = m.TextInput.Value()
			m.TextInput.Reset()
			return m, nil
		}

	}

	m.TextInput, cmd = m.TextInput.Update(msg)
	return m, cmd
}

func updateBooleanChoice(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			m.inputting = false
			return m, nil

		case "enter":
			tailscale.FinalCmd.Flags[m.choice] = m.answer
			m.inputting = false
			return m, nil

		case "down", "j":
			m.answer = "false"

		case "up", "k":
			m.answer = "true"
		}
	}

	return m, nil

}
