package commands

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

const useHighPerformanceRenderer = false

func init() {
	rootCmd.AddCommand(startCmd)
}

type App struct {
	Tabs      []Tab
	ActiveTab int
	viewport  viewport.Model
	ready     bool
}

type Tab struct {
	URL     string
	Status  string
	Content string
}

func (a App) Init() tea.Cmd {
	return nil
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return a, tea.Quit
		}

		if k := msg.String(); k == "ctrl+t" {
			content, err := ioutil.ReadFile("cucumber.md")
			if err != nil {
				fmt.Println("could not load file:", err)
				os.Exit(1)
			}

			rendered, _ := glamour.Render(string(content), "dark")

			a.ActiveTab++
			tab := Tab{
				URL:     "https://cucumber.com",
				Status:  "Done.",
				Content: rendered,
			}

			a.Tabs = append(a.Tabs, tab)
			a.viewport.SetContent(tab.Content)
		}
	case tea.WindowSizeMsg:
		footerHeight := lipgloss.Height(a.FooterView(&a.Tabs[a.ActiveTab]))
		verticalMarginHeight := footerHeight

		if !a.ready {
			a.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			a.viewport.YPosition = 0
			a.viewport.HighPerformanceRendering = useHighPerformanceRenderer
			a.viewport.SetContent(a.Tabs[a.ActiveTab].Content)
			a.ready = true
		} else {
			a.viewport.Width = msg.Width
			a.viewport.Height = msg.Height - verticalMarginHeight
		}
	}

	// Handle keyboard and mouse events in the viewport
	a.viewport, cmd = a.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return a, tea.Batch(cmds...)
}

func (a App) View() string {
	if !a.ready {
		return "\n  Initializing..."
	}

	return fmt.Sprintf("%s\n%s",
		a.viewport.View(),
		a.FooterView(&a.Tabs[a.ActiveTab]),
	)
}

func (a App) FooterView(tab *Tab) string {
	var statusStyle = lipgloss.NewStyle().Background(lipgloss.Color("205")).PaddingRight(2).PaddingLeft(2)
	var urlStyle = lipgloss.NewStyle().Background(lipgloss.Color("237")).PaddingLeft(2)

	url := tab.URL + strings.Repeat(" ", max(0, a.viewport.Width-lipgloss.Width(tab.URL)))

	status := statusStyle.Render(tab.Status)
	url = urlStyle.Render(url)

	return lipgloss.JoinHorizontal(lipgloss.Center, status, url)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func boot(cmd *cobra.Command, args []string) {
	content, err := ioutil.ReadFile("artichoke.md")
	if err != nil {
		fmt.Println("could not load file:", err)
		os.Exit(1)
	}

	// Replace with a custom renderer for gophertext / gemtext
	rendered, _ := glamour.Render(string(content), "dark")

	tabs := []Tab{}

	tab := Tab{
		URL:     "http://artichoke.com",
		Status:  "Done.",
		Content: rendered,
	}

	tabs = append(tabs, tab)

	app := App{
		Tabs:      tabs,
		ActiveTab: 0,
	}

	errs := make(chan error, 1)
	prog := tea.NewProgram(
		app,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
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
