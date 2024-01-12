package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

var finalCmd string = "tailscale up "

func (i flag) Title() string       { return i.name }
func (i flag) Description() string { return i.description }
func (i flag) FilterValue() string { return i.name }

type model struct {
	list      list.Model
	choice    flag
	answer    string
	inputting bool
	quitting  bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				// finalCmd += m.choice.name + " "

			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) booleanInput() string {
	choices := []string{"true", "false"}
	s := strings.Builder{}

	s.WriteString("What value do you want to pass to?\n\n")
	for i := 0; i < len(choices); i++ {
		if m.answer == string(choices[i]) {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(string(choices[i]))
		s.WriteString("\n")

	}
	s.WriteString("\n(press q to quit)\n")

	return s.String()

}

func (m model) View() string {
	if m.inputting && m.choice.isBooleanFlag {
		m.inputting = false
		return m.booleanInput()
	}

	if m.choice.name != "" {
		return quitTextStyle.Render(finalCmd) + "\n" +
			m.list.View()
	}

	// if m.quitting {
	// 	return quitTextStyle.Render("Not hungry? That’s cool.")
	// }

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

	l := list.New(items, list.NewDefaultDelegate(), 0, 35)
	l.Title = "Select flags to pass to tailscale up"

	m := model{list: l}

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
