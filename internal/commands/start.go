package commands

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/netsensei/hector/internal/ui"
	"github.com/spf13/cobra"
)

const useHighPerformanceRenderer = false

const INIT = "initializing"
const EXIT = "exiting"
const READY = "ready"
const CMND = "command"
const INPUT = "input"
const VIEW = "view"

func init() {
	rootCmd.AddCommand(startCmd)
}

type App struct {
	Tabs      *ui.Tabs
	ActiveTab int
	viewport  viewport.Model
	navInput  textinput.Model
	state     string
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
		if k := msg.String(); k == "ctrl+c" || k == "esc" {
			a.state = EXIT
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

		if k := msg.String(); k == "ctrl+p" {
			a.Tabs.Down()
		}

		if k := msg.String(); k == "ctrl+n" {
			a.Tabs.Up()
		}

		if k := msg.String(); k == "ctrl+u" {
			a.state = INPUT
			a.navInput.Focus()
		}

		if k := msg.String(); k == "ctrl+v" {
			a.state = VIEW
		}

		if a.state == INPUT {
			if k := msg.String(); k == "enter" {
				tab, _ := a.Tabs.Current()

				fileName := a.navInput.Value()
				content, err := ioutil.ReadFile(fileName)
				if err != nil {
					tab.Content = fmt.Sprintf("could not load file: %s", fileName)
				} else {
					rendered, _ := glamour.Render(string(content), "dark")
					tab.Content = rendered
				}
			}
		}

	case tea.WindowSizeMsg:
		footerHeight := lipgloss.Height(a.FooterView())
		verticalMarginHeight := footerHeight

		if a.state != READY {
			a.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			a.viewport.YPosition = 0
			a.viewport.HighPerformanceRendering = useHighPerformanceRenderer
			a.state = READY
		} else {
			a.viewport.Width = msg.Width
			a.viewport.Height = msg.Height - verticalMarginHeight
		}
	}

	// Handle keyboard and mouse events in the viewport
	if a.state == INPUT {
		a.navInput, cmd = a.navInput.Update(msg)
	} else {
		a.viewport, cmd = a.viewport.Update(msg)
	}

	cmds = append(cmds, cmd)

	return a, tea.Batch(cmds...)
}

func (a App) View() string {
	if a.state != READY && a.state == INIT {
		return "\n  Initializing..."
	}

	return fmt.Sprintf("%s\n%s",
		a.CanvasView(),
		a.FooterView(),
	)
}

func (a App) CanvasView() string {
	tab, _ := a.Tabs.Current()
	a.viewport.SetContent(tab.Content)
	return a.viewport.View()
}

func (a App) FooterView() string {
	tab, activeTab := a.Tabs.Current()
	count := a.Tabs.Count()

	var statusStyle = lipgloss.NewStyle().Background(lipgloss.Color("205")).PaddingRight(2).PaddingLeft(2)
	var tabStyle = lipgloss.NewStyle().Background(lipgloss.Color("205")).PaddingRight(2).PaddingLeft(2)
	var urlStyle = lipgloss.NewStyle().Background(lipgloss.Color("237")).PaddingLeft(2)

	var tabIndicatorStr string

	if count == 1 {
		tabIndicatorStr = fmt.Sprintf("tab %d", activeTab)
	} else {
		if activeTab == 0 {
			tabIndicatorStr = fmt.Sprintf("tab %d \u00BB", activeTab)
		} else if activeTab > 0 && activeTab < count-1 {
			tabIndicatorStr = fmt.Sprintf("\u00AB tab %d \u00BB", activeTab)
		} else {
			tabIndicatorStr = fmt.Sprintf("\u00AB tab %d", activeTab)
		}
	}

	if a.state == INPUT {
		a.navInput.Focus()
	} else {
		a.navInput.Reset()
	}

	navInput := a.navInput.View()
	url := navInput // + strings.Repeat(" ", max(0, a.viewport.Width-lipgloss.Width(navInput)-lipgloss.Width(tabIndicatorStr)))

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

	url := "http://artichoke.com"

	tab := ui.Tab{
		URL:     url,
		Status:  "Done.",
		Content: rendered,
	}

	tabs.Add(tab)

	navInput := textinput.New()
	navInput.Placeholder = url
	navInput.CharLimit = 510
	navInput.Width = 255

	app := App{
		Tabs:      tabs,
		ActiveTab: 0,
		state:     INIT,
		navInput:  navInput,
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
