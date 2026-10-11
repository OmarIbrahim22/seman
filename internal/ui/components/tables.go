package components

import (
	"fmt"
	"strings"

	"github.com/OmarIbrahim22/seman/internal/models"
	"github.com/OmarIbrahim22/seman/internal/style"
)

func RenderProjectsTable(items []models.ProjectItem, selected, width int, t style.Theme) string {
	if width <= 0 {
		return ""
	}
	header := []string{"Project Name", "Subject", "Deadline", "Status"}

	colSubject := 10
	colDeadline := 16
	colStatus := 14
	colName := width - (colSubject + colDeadline + colStatus + 3)
	if colName < 16 {
		colName = 16
	}

	line := fmt.Sprintf(
		"%-*s %-*s %-*s %-*s",
		colName, header[0],
		colSubject, header[1],
		colDeadline, header[2],
		colStatus, header[3],
	)

	var b strings.Builder
	b.WriteString(t.Dim.Render(line))

	for i, item := range items {
		b.WriteString("\n")
		name := TruncateString(item.Name, colName)
		subject := TruncateString(item.Subject, colSubject)
		due := TruncateString(item.Due, colDeadline)
		status := RenderStatusBadge(item.Status, colStatus, t)
		row := fmt.Sprintf(
			"%-*s %-*s %-*s %-*s",
			colName, name,
			colSubject, subject,
			colDeadline, due,
			colStatus, status,
		)
		if i == selected {
			b.WriteString(t.RowActive.Render(row))
		} else {
			b.WriteString(t.Text.Render(row))
		}
	}
	return b.String()
}
