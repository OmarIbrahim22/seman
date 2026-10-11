package components

import (
	"fmt"
	"strings"
	"time"

	"github.com/OmarIbrahim22/seman/internal/models"
	"github.com/OmarIbrahim22/seman/internal/style"
)

func RenderSubjects(items []models.SubjectItem, selected int, t style.Theme) string {
	if len(items) == 0 {
		return t.Dim.Render("No subjects yet")
	}

	var b strings.Builder
	for i, subj := range items {
		if i > 0 {
			b.WriteString("\n\n")
		}
		prefix := "  "
		codeStyle := t.SubjectDim
		nameStyle := t.SubjectDim
		if i == selected {
			prefix = "> "
			codeStyle = t.SubjectActive
			nameStyle = t.SubjectActive
		}
		b.WriteString(codeStyle.Render(prefix + subj.Code))
		if subj.Name != "" {
			b.WriteString("\n")
			b.WriteString(nameStyle.Render("  " + subj.Name))
		}
	}
	return b.String()
}

func RenderFlatExams(exams []models.FlatExam, cursor int, filter string, t style.Theme) string {
	if len(exams) == 0 {
		if filter != "" {
			return t.Dim.Render("No exams for " + filter + " in this period")
		}
		return t.Dim.Render("No exams — press [A] to add one")
	}

	now := time.Now()
	var b strings.Builder
	for i, flat := range exams {
		if i > 0 {
			b.WriteString("\n\n")
		}
		selected := i == cursor
		nameStyle := t.Text
		prefix := "  "
		if selected {
			nameStyle = t.RowActive
			prefix = "> "
		}

		// Line 1: prefix + name + subject tag (when not filtered)
		nameLine := prefix + flat.Exam.Name
		if filter == "" && flat.SubjectCode != "" {
			nameLine += "  (" + flat.SubjectCode + ")"
		}
		if flat.Exam.Priority != "" {
			b.WriteString(nameStyle.Render(nameLine) + "  ")
			b.WriteString(RenderPriority(flat.Exam.Priority))
		} else {
			b.WriteString(nameStyle.Render(nameLine))
		}

		// Line 2: date + countdown
		if flat.Exam.Date != "" {
			b.WriteString("\n")
			dateInfo := "  " + flat.Exam.Date
			if date, ok := parseExamDate(flat.Exam.Date); ok {
				days := int(date.Sub(now).Hours() / 24)
				var countdown string
				switch {
				case days < -1:
					countdown = fmt.Sprintf(" · %d days ago", -days)
				case days == -1:
					countdown = " · yesterday"
				case days == 0:
					countdown = " · today"
				case days == 1:
					countdown = " · tomorrow"
				case days < 14:
					countdown = fmt.Sprintf(" · in %d days", days)
				default:
					countdown = fmt.Sprintf(" · in %d weeks", days/7)
				}
				dateInfo += countdown
			}
			b.WriteString(t.Dim.Render(dateInfo))
		}

		// Retakes (compact, single line)
		if len(flat.Exam.Retakes) > 0 {
			b.WriteString("\n")
			b.WriteString(t.Dim.Render("  Retakes: " + strings.Join(flat.Exam.Retakes, ", ")))
		}
	}
	return b.String()
}
