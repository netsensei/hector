package commands

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(startCmd)
}

type App struct {
	Tabs      []Tab
	ActiveTab int
}

type Tab struct {
	URL    string
	Status string
}

func (a App) Init() tea.Cmd {
	return nil
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return a, tea.Quit
		}
	}
	return a, nil
}

func (a App) View() string {
	return "string"
}

func boot(cmd *cobra.Command, args []string) {
	tabs := make([]Tab, 0)
	app := App{
		Tabs:      tabs,
		ActiveTab: 0,
	}

	errs := make(chan error, 1)
	prog := tea.NewProgram(app)
	go func() {
		err := <-errs
		if err != nil {
			log.Print(err)
		}
	}()

	errs <- prog.Start()
	prog.Kill()
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Hector",
	Run:   boot,
}
