package screens

import (
	"github.com/OmarIbrahim22/seman/internal/style"
	"github.com/OmarIbrahim22/seman/internal/ui/components"
	"github.com/charmbracelet/lipgloss"
)

func RenderDashboard(state State, width, height int, t style.Theme) string {
	layout := components.ComputeDashboardLayout(width, height)

	examTitle := "Upcoming Exams (Priority)"
	todoTitle := "Todos"
	projectTitle := "Projects (Soonest Deadlines)"
	if state.SubjectFilter != "" {
		examTitle = "Exams · " + state.SubjectFilter
		todoTitle = "Todos · " + state.SubjectFilter
		projectTitle = "Projects · " + state.SubjectFilter
	}

	leftBody := components.RenderUpcomingExams(state.Subjects, 5, state.FilterStart, state.FilterEnd, state.FilterAll, state.SubjectFilter, t)
	leftPanel := components.RenderPanel(layout.LeftWidth, layout.PanelHeight, examTitle, leftBody, t)

	if layout.MiddleWidth <= 0 {
		return leftPanel
	}

	middlePanel := components.RenderPanel(layout.MiddleWidth, layout.PanelHeight, todoTitle, state.ChecklistView, t)
	if layout.RightWidth <= 0 {
		return lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, " ", middlePanel)
	}

	rightPanel := components.RenderPanel(layout.RightWidth, layout.PanelHeight, projectTitle, components.RenderProjects(state.Projects, t), t)
	return lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, " ", middlePanel, " ", rightPanel)
}
