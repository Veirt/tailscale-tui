package main

import (
	"fmt"
	"os"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/veirt/tailscale-tui/internal/tailscale"
	"github.com/veirt/tailscale-tui/internal/tui"
)

func launchTUI(flags []tailscale.Flag) {
	m := tui.InitialModel()

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

}

func main() {
	if err := tailscale.CheckTailscale(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	flags := tailscale.GetTailscaleUpFlags()

	launchTUI(flags)

}
