package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/veirt/tailscale-tui/internal/tailscale"
	"strings"
)

var (
	quitTextStyle = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

func (m model) View() string {
	if m.quitting {
		return quitTextStyle.Render(tailscale.FinalCmd.String())
	}

	if m.inputting && m.choice.IsBooleanFlag {
		m.inputting = false
		return booleanInputView(m)
	}

	if m.inputting && !m.choice.IsBooleanFlag {
		m.inputting = false
		return regularInputView(m)
	}

	if m.choice.Name != "" {
		return quitTextStyle.Render(tailscale.FinalCmd.String()) + "\n" +
			m.list.View()
	}

	return "\n" + m.list.View()
}

var boolChoices = []string{"true", "false"}

func booleanInputView(m model) string {
	s := strings.Builder{}

	s.WriteString(m.choice.Description())
	s.WriteString("\n")
	s.WriteString("What value do you want to pass to?\n\n")
	for i := 0; i < len(boolChoices); i++ {
		if m.answer == string(boolChoices[i]) {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(string(boolChoices[i]))
		s.WriteString("\n")

	}

	s.WriteString("\n(press q to go back)\n")

	return s.String()
}

func regularInputView(m model) string {
	s := strings.Builder{}

	s.WriteString(m.choice.Description())
	s.WriteString("\n")
	s.WriteString("What value do you want to pass to?\n\n")
	s.WriteString(m.TextInput.View())

	return s.String()
}
