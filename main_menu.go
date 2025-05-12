package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	list_style = lipgloss.NewStyle()
	docStyle   = lipgloss.NewStyle().
			Width(50).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center)
)

func GetDocStyle(vertical int, horizontal int) lipgloss.Style {
	style := lipgloss.NewStyle().Width(horizontal).Height(vertical).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center)
	return style
}

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func MainMenu() MainMenuModel {
	items := []list.Item{
		item{title: "Login", desc: "Login with an existing agent"},
		item{title: "Register", desc: "Register a new agent"},
	}
	delegate := list.NewDefaultDelegate()
	delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.AlignHorizontal(lipgloss.Center)
	list := list.New(items, delegate, 20, 20)
	list.Styles.Title = list.Styles.Title.AlignHorizontal(lipgloss.Center)
	list.SetFilteringEnabled(false)
	list.SetShowFilter(false)
	list.SetShowTitle(false)
	list.SetShowStatusBar(false)
	list.SetShowHelp(false)

	return MainMenuModel{list: list}
}

type MainMenuModel struct {
	list   list.Model
	width  int
	height int
}

func (m MainMenuModel) Init() tea.Cmd {
	return nil
}

func (m MainMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width-5, msg.Height-5)
		m.width = msg.Width
		m.height = msg.Height
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m MainMenuModel) View() string {
	return GetDocStyle(m.height, m.width).Render(m.list.View())
}
