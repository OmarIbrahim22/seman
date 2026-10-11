package components

import (
	"fmt"
	"strings"

	"github.com/OmarIbrahim22/seman/internal/style"
	"github.com/charmbracelet/lipgloss"
)

type TabItem struct {
	ID    int
	Label string
}

func RenderTabs(active, width int, weekLabel string, items []TabItem, t style.Theme) string {
	tabs := make([]string, 0, len(items))
	for i, item := range items {
		text := fmt.Sprintf("[%d] %s", i+1, item.Label)
		if item.ID == active {
			tabs = append(tabs, t.TabActive.Render(text))
		} else {
			tabs = append(tabs, t.TabInactive.Render(text))
		}
	}
	left := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
	left = lipgloss.NewStyle().PaddingLeft(1).Render(left)

	if weekLabel == "" || width <= lipgloss.Width(left) {
		return left
	}

	rightLabel := "< " + weekLabel + " >"
	right := t.Dim.Render(rightLabel)
	rightWidth := width - lipgloss.Width(left)
	if rightWidth < 1 {
		return left
	}
	right = lipgloss.NewStyle().Width(rightWidth).Height(3).Align(lipgloss.Right, lipgloss.Center).Render(right)
	if lipgloss.Width(right) > rightWidth {
		right = strings.TrimSpace(lipgloss.NewStyle().Width(rightWidth).Render(right))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}
