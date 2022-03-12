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
	"github.com/netsensei/hector/internal/ui"
	"github.com/spf13/cobra"
)

const useHighPerformanceRenderer = false

func init() {
	rootCmd.AddCommand(startCmd)
}

type App struct {
	Tabs      *ui.Tabs
	ActiveTab int
	viewport  viewport.Model
	ready     bool
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
			tab := ui.Tab{
				URL:     "https://cucumber.com",
				Status:  "Done.",
				Content: rendered,
			}

			a.Tabs.Add(tab)
		}

		if k := msg.String(); k == "ctrl+x" {
			a.Tabs.Remove()
		}
	case tea.WindowSizeMsg:
		footerHeight := lipgloss.Height(a.FooterView(a.Tabs.Current()))
		verticalMarginHeight := footerHeight

		if !a.ready {
			a.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			a.viewport.YPosition = 0
			a.viewport.HighPerformanceRendering = useHighPerformanceRenderer
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

	tab, activeTab := a.Tabs.Current()

	return fmt.Sprintf("%s\n%s",
		a.CanvasView(tab),
		a.FooterView(tab, activeTab),
	)
}

func (a App) CanvasView(tab *ui.Tab) string {
	a.viewport.SetContent(tab.Content)
	return a.viewport.View()
}

func (a App) FooterView(tab *ui.Tab, activeTab int) string {
	var statusStyle = lipgloss.NewStyle().Background(lipgloss.Color("205")).PaddingRight(2).PaddingLeft(2)
	var tabStyle = lipgloss.NewStyle().Background(lipgloss.Color("205")).PaddingRight(2).PaddingLeft(2)
	var urlStyle = lipgloss.NewStyle().Background(lipgloss.Color("237")).PaddingLeft(2)

	tabIndicatorStr := fmt.Sprintf("tab %d", activeTab)
	url := tab.URL + strings.Repeat(" ", max(0, a.viewport.Width-lipgloss.Width(tab.URL)-lipgloss.Width(tabIndicatorStr)))

	status := statusStyle.Render(tab.Status)
	tabIndicator := tabStyle.Render(tabIndicatorStr)
	url = urlStyle.Render(url)

	return lipgloss.JoinHorizontal(lipgloss.Center, status, tabIndicator, url)
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

	tabs := ui.NewTabs()

	tab := ui.Tab{
		URL:     "http://artichoke.com",
		Status:  "Done.",
		Content: rendered,
	}

	tabs.Add(tab)

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
