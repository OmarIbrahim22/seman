package screens

import (
	"strings"

	"github.com/OmarIbrahim22/seman/internal/style"
	"github.com/OmarIbrahim22/seman/internal/ui/components"
	"github.com/charmbracelet/lipgloss"
)

func RenderWeeklyFocus(state State, width, height int, t style.Theme) string {
	layout := components.ComputeWeeklyLayout(width, height)

	title := "Todos"
	var header string
	if state.SubjectFilter != "" {
		title = "Todos · " + state.SubjectFilter
		header = t.Title.Render("Subject: "+state.SubjectFilter) + "  " + t.Dim.Render("[R] clear filter") + "\n\n"
	}

	body := header + state.ChecklistView
	var parts []string
	rightPanel := components.RenderPanel(layout.RightWidth, layout.PanelsHeight, title, body, t)
	parts = append(parts, rightPanel)

	if layout.SpacerHeight > 0 {
		parts = append(parts, strings.Repeat("\n", layout.SpacerHeight-1))
	}

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}
