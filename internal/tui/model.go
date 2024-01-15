package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/veirt/tailscale-tui/internal/tailscale"
)

type model struct {
	list      list.Model
	choice    tailscale.Flag
	TextInput textinput.Model
	answer    string
	inputting bool
	quitting  bool
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func InitialModel() model {
	flags := tailscale.GetTailscaleUpFlags()
	items := []list.Item{}

	for _, flag := range flags {
		items = append(items, flag)
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 30)
	l.Title = "Select flags to pass to tailscale up"

	ti := textinput.New()
	ti.Placeholder = "Enter here..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 30

	m := model{list: l, TextInput: ti}

	return m
}
