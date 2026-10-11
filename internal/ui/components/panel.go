package components

import (
	"fmt"
	"strings"

	"github.com/OmarIbrahim22/seman/internal/style"
	"github.com/charmbracelet/lipgloss"
)

func RenderPanel(width, height int, title, body string, t style.Theme) string {
	if width <= 0 || height <= 0 {
		return ""
	}

	contentW, _ := PanelContentSize(width, height)
	styleW := width - panelBorderX
	styleH := height - panelBorderY
	if styleW < 1 {
		styleW = 1
	}
	if styleH < 1 {
		styleH = 1
	}

	box := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(t.Border).
		Padding(panelPaddingY, panelPaddingX).
		Width(styleW).
		Height(styleH)

	content := body
	if title != "" {
		titleLine := RenderPanelTitle(title, contentW, t)
		if body != "" {
			content = titleLine + "\n" + body
		} else {
			content = titleLine
		}
	}

	return box.Render(content)
}

func RenderPanelTitle(title string, width int, t style.Theme) string {
	if width <= 0 {
		return ""
	}

	prefix := fmt.Sprintf("- %s ", title)
	if lipgloss.Width(prefix) > width {
		prefix = TruncateString(prefix, width)
	}
	if lipgloss.Width(prefix) < width {
		prefix += strings.Repeat("-", width-lipgloss.Width(prefix))
	}
	return t.Title.Render(prefix)
}
