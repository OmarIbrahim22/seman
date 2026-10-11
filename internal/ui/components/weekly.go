package components

import (
	"fmt"
	"strings"

	"github.com/OmarIbrahim22/seman/internal/style"
	"github.com/charmbracelet/lipgloss"
)

func RenderWeekHeader(label string, width int, t style.Theme) string {
	if label == "" {
		label = "Week 1 - TBD"
	}
	text := fmt.Sprintf("< %s >", label)
	box := lipgloss.NewStyle().Width(width).Align(lipgloss.Center)
	return t.Dim.Render(box.Render(text))
}

func RenderWeeklyExams(exams []string, t style.Theme) string {
	if len(exams) == 0 {
		return t.Dim.Render("No exams this week")
	}
	var b strings.Builder
	for i, exam := range exams {
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString(t.Text.Render("- " + exam))
	}
	return b.String()
}
