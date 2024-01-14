package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/veirt/tailscale-tui/internal/tailscale"
	"github.com/veirt/tailscale-tui/internal/tui"
	"os"
)

func launchTUI(flags []tailscale.Flag) {
	m := tui.InitialModel()

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

}

func version() string {
	return "tailscale-tui v0.1.0"
}

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "-v" || os.Args[1] == "--version" {
			fmt.Println(version())
			os.Exit(0)
		}
	}

	if err := tailscale.CheckTailscale(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	flags := tailscale.GetTailscaleUpFlags()

	launchTUI(flags)

}
