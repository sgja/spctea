package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/go-resty/resty/v2"
)

const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1).Foreground(lipgloss.Color("241"))
	quitTextStyle     = lipgloss.NewStyle()
)

func (i Faction) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Faction)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.Name)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type model struct {
	list      list.Model
	choice    string
	quitting  bool
	viewport  viewport.Model
	step      int
	callsign  string
	textInput textinput.Model
	client    *resty.Client
	faction   Faction
	token     string
}

func newModel(client *resty.Client, factions []Faction) (*model, error) {

	items := []list.Item{}
	for _, f := range factions {
		items = append(items, f)
	}

	const defaultWidth = 20

	ti := textinput.New()
	ti.Placeholder = "Agent X"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Choose a faction"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	vp := viewport.New(64, 20)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingRight(2)
	vp.MouseWheelEnabled = true

	fp := filepicker.New()
	fp.CurrentDirectory, _ = os.UserHomeDir()
	return &model{viewport: vp, list: l, step: 0, textInput: ti, client: client}, nil
}

func (m model) helpView() string {
	return helpStyle.PaddingLeft(0).PaddingTop(0).Render("\n  ↑/↓: Navigate • backspace: Go back • q: Quit\n")
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.viewport.Width = msg.Width
		return m, nil
	default:
		switch m.step {
		case 0:
			return m.updateText(msg)
		case 1:
			return m.updateList(msg)
		case 2:
			return m.updateDetails(msg)
		case 3:
			return m.updateSave(msg)
		default:
			return m, nil

		}
	}
}

func (m model) updateText(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			m.callsign = m.textInput.Value()
			m.step = 1
			return m, nil
		default:
			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		}
	default:
		return m, nil
	}
}

func (m model) updateList(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "backspace":
			m.step = 0
			return m, nil

		case "enter":
			i, ok := m.list.SelectedItem().(Faction)
			if ok {
				m.step = 2
				m.faction = i
				m.choice = format_faction(i, m.viewport.Width-4)
				renderer, err := glamour.NewTermRenderer(
					glamour.WithAutoStyle(),
					glamour.WithWordWrap(m.viewport.Width-4),
				)
				if err != nil {
					return m, tea.Quit
				}

				str, err := renderer.Render(m.choice)
				if err != nil {
					return m, tea.Quit
				}
				m.viewport.SetContent(str)
			}
			return m, nil
		default:
			var cmd tea.Cmd
			m.list, cmd = m.list.Update(msg)
			return m, cmd
		}

	default:
		return m, nil
	}
}

func (m model) updateDetails(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q", "esc":
			m.quitting = true
			return m, tea.Quit

		case "backspace":
			m.step = 1
			return m, nil

		case "enter":
			m.step = 3
			token, agent, err := register(m.client, m.callsign, m.faction.Symbol)
			if err != nil {
				return m, tea.Quit
			}
			homedir, err := os.UserHomeDir()
			if err != nil {
				return m, tea.Quit
			}
			agentpath := path.Join(homedir, ".spctea", "agents", fmt.Sprintf("%s.json", m.callsign))
			os.MkdirAll(path.Dir(agentpath), 0777)
			os.WriteFile(agentpath, []byte(agent), 0666)
			m.token = token
			return m, tea.Quit

		default:
			var cmd tea.Cmd
			m.viewport, cmd = m.viewport.Update(msg)
			return m, cmd
		}
	case tea.MouseMsg:
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	default:
		return m, nil
	}
}

func (m model) updateSave(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q", "esc":
			m.quitting = true
			return m, tea.Quit

		case "backspace":
			m.step = 2
			return m, nil

		case "enter":
			m.step = 4
			return m, tea.Quit

		default:
			return m, nil
		}
	default:
		return m, nil
	}
}

func format_faction(faction Faction, width int) string {
	text := fmt.Sprintf("# %s\n\n", faction.Name) +
		fmt.Sprintf("HQ: %s\n\n", faction.Headquarters) +
		fmt.Sprintf("Traits:\n\n")
	for _, t := range faction.Traits {
		text += fmt.Sprintf("## %s\n\n", t.Name)
		text += fmt.Sprintf("%s\n\n", t.Description)
	}
	return text
}

func (m model) View() string {
	switch m.step {
	case 0:
		return m.textInput.View()
	case 1:
		return "\n" + m.list.View()
	case 2:
		return m.viewport.View() + m.helpView()
	case 3:
		return quitTextStyle.Width(m.viewport.Width - 8).Render(m.token)
	}

	return ""
}

func show_factions(client *resty.Client, factions []Faction) {

	m, _ := newModel(client, factions)

	if _, err := tea.NewProgram(m, tea.WithMouseAllMotion()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
