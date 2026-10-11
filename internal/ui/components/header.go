package components

import (
	"time"

	"github.com/OmarIbrahim22/seman/internal/style"
	"github.com/charmbracelet/lipgloss"
)

func RenderHeader(width int, t style.Theme) string {
	contentWidth := width - barBorderX - barPaddingX*2
	if contentWidth < 1 {
		contentWidth = 1
	}

	title := t.Title.Render("Student Manager")
	date := t.Dim.Render(time.Now().Format("01-02-2006"))
	content := AlignLine(contentWidth, title, date)

	styleWidth := width - barBorderX
	if styleWidth < 1 {
		styleWidth = 1
	}

	box := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(t.Border).
		Padding(barPaddingY, barPaddingX).
		Width(styleWidth)

	return box.Render(content)
}
