package main

import (
	"bufio"
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func checkTailscale() error {
	_, err := exec.LookPath("tailscale")

	return err
}

type flag struct {
	name          string
	description   string
	isBooleanFlag bool
}

var (
	quitTextStyle = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

var boolChoices = []string{"true", "false"}

type command struct {
	name  string
	flags map[flag]string
}

func (c command) String() string {
	result := c.name + " "

	// c.flags is a map of flag -> value
	// turn it into a string, sorted by flag name

	// sort the flags by name
	var keys []flag
	for k := range c.flags {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].name < keys[j].name
	})

	for _, k := range keys {
		if k.isBooleanFlag {
			result += k.name + "=" + c.flags[k] + " "
		} else {
			result += k.name + " \"" + c.flags[k] + "\" "
		}
	}

	return result
}

var finalCmd command = command{name: "tailscale up", flags: map[flag]string{}}

func (i flag) Title() string       { return i.name }
func (i flag) Description() string { return i.description }
func (i flag) FilterValue() string { return i.name }

type model struct {
	list      list.Model
	choice    flag
	TextInput textinput.Model
	answer    string
	inputting bool
	quitting  bool
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.inputting && m.choice.isBooleanFlag {
		return updateBooleanChoice(msg, m)
	}

	if m.inputting && !m.choice.isBooleanFlag {
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
			i, ok := m.list.SelectedItem().(flag)
			if ok {
				m.choice = i

				m.inputting = true
				if m.choice.isBooleanFlag {
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
			finalCmd.flags[m.choice] = m.TextInput.Value()
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
			finalCmd.flags[m.choice] = m.answer
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

func (m model) View() string {
	if m.quitting {
		return quitTextStyle.Render(finalCmd.String())
	}

	if m.inputting && m.choice.isBooleanFlag {
		m.inputting = false
		return booleanInputView(m)
	}

	if m.inputting && !m.choice.isBooleanFlag {
		m.inputting = false
		return regularInputView(m)
	}

	if m.choice.name != "" {
		return quitTextStyle.Render(finalCmd.String()) + "\n" +
			m.list.View()
	}

	return "\n" + m.list.View()
}

func getTailscaleUpFlags() []flag {
	out, err := exec.Command("tailscale", "up", "--help").CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		if scanner.Text() == "FLAGS" {
			break
		}
	}

	cmds := []flag{}
	for scanner.Scan() {
		out := scanner.Text()
		out = strings.TrimSpace(out)

		// if it starts with --, it's a flag
		if strings.HasPrefix(out, "--") {

			isBooleanFlag := strings.Contains(out, "false")
			name := ""
			if isBooleanFlag {
				name = out[:strings.IndexByte(out, ',')]
			} else {
				name = out[:strings.IndexByte(out, ' ')]
			}

			scanner.Scan()
			desc := strings.TrimSpace(scanner.Text())

			fl := flag{name: name, isBooleanFlag: isBooleanFlag, description: desc}
			cmds = append(cmds, fl)

		}
	}

	return cmds

}

func launchTui(flags []flag) {
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

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

}

func main() {
	if err := checkTailscale(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	flags := getTailscaleUpFlags()

	launchTui(flags)

}
