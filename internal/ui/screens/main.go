package screens

import (
	"github.com/OmarIbrahim22/seman/internal/style"
	"github.com/OmarIbrahim22/seman/internal/ui/components"
	"github.com/charmbracelet/lipgloss"
)

func RenderMain(state State, width, height int, t style.Theme) string {
	switch state.ActiveTab {
	case 0: // tabDashboard
		return RenderDashboard(state, width, height, t)
	case 1: // tabExams
		return RenderExams(state, width, height, t)
	case 2: // tabTodos
		return RenderWeeklyFocus(state, width, height, t)
	case 3: // tabProjects
		return RenderProjectsTab(state, width, height, t)
	case 4: // tabSettings
		return RenderSettingsTab(width, height, state.ConfirmOn, state.WeekSpan, state.LofiEnabled, state.LofiURL, state.ThemeName, t)
	case 5: // tabLofi
		return RenderLofi(state, width, height, t)
	case 6: // tabSubjects
		return RenderSubjectsTab(state, width, height, t)
	default:
		return RenderPlaceholder(width, height, t)
	}
}

func RenderPlaceholder(width, height int, t style.Theme) string {
	message := t.Dim.Render("Coming soon")
	box := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Align(lipgloss.Center, lipgloss.Center)
	return box.Render(message)
}

func RenderModal(state State, width, height int, t style.Theme) string {
	return components.RenderModalArea(state.Modal, width, height, t)
}
