package components

import (
	"strings"

	"github.com/OmarIbrahim22/seman/internal/style"
	"github.com/charmbracelet/lipgloss"
)

func RenderStatusBadge(status string, width int, t style.Theme) string {
	label := strings.ToUpper(status)
	style := t.StatusNotStr
	switch label {
	case "DONE":
		style = t.StatusDone
	case "IN PROGRESS":
		style = t.StatusInProg
	case "NOT STARTED":
		style = t.StatusNotStr
	}
	badge := style.Render(label)
	if lipgloss.Width(badge) < width {
		return badge + strings.Repeat(" ", width-lipgloss.Width(badge))
	}
	return TruncateString(badge, width)
}

func RenderPriority(priority string) string {
	style := lipgloss.NewStyle().Padding(0, 1).Bold(true)
	switch strings.ToUpper(priority) {
	case "HIGH":
		style = style.Foreground(lipgloss.Color("#0c0c0c")).Background(lipgloss.Color("#ff5f5f"))
	case "MED":
		style = style.Foreground(lipgloss.Color("#0c0c0c")).Background(lipgloss.Color("#d4a017"))
	case "LOW":
		style = style.Foreground(lipgloss.Color("#0c0c0c")).Background(lipgloss.Color("#3a7f3a"))
	default:
		style = style.Foreground(lipgloss.Color("#0c0c0c")).Background(lipgloss.Color("#39ff14"))
	}
	return style.Render(priority)
}
