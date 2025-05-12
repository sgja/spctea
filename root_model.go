package main

import (
	"sgja/spctea/backend"

	tea "github.com/charmbracelet/bubbletea"
)

type RootModel struct {
	model tea.Model
	app   *backend.App
}

func RootScreen(app *backend.App) RootModel {
	var rootModel tea.Model

	main_menu := MainMenu()
	rootModel = &main_menu
	return RootModel{rootModel, app}
}

func (m RootModel) Init() tea.Cmd {
	return m.model.Init()
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.model.Update(msg)
}

func (m RootModel) View() string {
	return m.model.View()
}

func (m RootModel) SwitchScreen(model tea.Model) (tea.Model, tea.Cmd) {
	m.model = model
	return m.model, m.model.Init()
}
