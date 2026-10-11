package components

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/OmarIbrahim22/seman/internal/models"
	"github.com/OmarIbrahim22/seman/internal/style"
)

type upcomingExam struct {
	Subject  string
	Name     string
	Priority string
	Date     time.Time
}

func RenderUpcomingExams(subjects []models.SubjectItem, limit int, start, end time.Time, all bool, filter string, t style.Theme) string {
	exams := collectUpcomingExams(subjects, start, end, all, filter)

	if len(exams) == 0 {
		if filter != "" {
			return t.Dim.Render("No exams for " + filter)
		}
		return t.Dim.Render("No exams this week")
	}

	if limit <= 0 || limit > len(exams) {
		limit = len(exams)
	}

	var b strings.Builder
	for i := 0; i < limit; i++ {
		if i > 0 {
			b.WriteString("\n")
		}
		item := exams[i]
		b.WriteString(t.Text.Render(fmt.Sprintf("- %s (%s)", item.Name, item.Subject)))
		b.WriteString("\n")
		dateStr := item.Date.Format("Jan 2, 2006 @ 15:04")
		b.WriteString(t.Dim.Render("  " + dateStr + " "))
		if item.Priority != "" {
			b.WriteString(RenderPriority(item.Priority))
		}
	}

	if len(exams) > limit {
		b.WriteString("\n")
		b.WriteString(t.Dim.Render(fmt.Sprintf("...and %d more", len(exams)-limit)))
	}

	return b.String()
}

func collectUpcomingExams(subjects []models.SubjectItem, start, end time.Time, all bool, filter string) []upcomingExam {
	list := make([]upcomingExam, 0)
	for _, subject := range subjects {
		if filter != "" && !strings.EqualFold(subject.Code, filter) {
			continue
		}
		for _, exam := range subject.Exams {
			date, ok := parseExamDate(exam.Date)
			if !ok {
				continue
			}
			if !all && (date.Before(start) || !date.Before(end)) {
				continue
			}
			list = append(list, upcomingExam{
				Subject:  subject.Code,
				Name:     exam.Name,
				Priority: exam.Priority,
				Date:     date,
			})
		}
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Date.Before(list[j].Date)
	})
	return list
}

func parseExamDate(value string) (time.Time, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, false
	}
	layouts := []string{
		"Jan 2, 2006 @ 15:04",
		"Jan 2, 2006",
		"2006-01-02 15:04",
		"2006-01-02",
		"2.1.2006 15:04",
		"2.1.2006",
		"02.01.2006 15:04",
		"02.01.2006",
		"2/1/2006 15:04",
		"2/1/2006",
		"02/01/2006 15:04",
		"02/01/2006",
	}
	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, value, time.Local); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

func RenderChecklist(items []models.ChecklistItem, selected int, showCursor bool, t style.Theme) string {
	var b strings.Builder
	for i, item := range items {
		if i > 0 {
			b.WriteString("\n")
		}
		box := "[ ]"
		rowStyle := t.CheckboxTodo
		if item.Done {
			box = "[x]"
			rowStyle = t.CheckboxDone
		}
		if showCursor && i == selected {
			rowStyle = t.RowActive
		}
		if item.Subject != "" {
			b.WriteString(rowStyle.Render(box) + " " + t.SubjectTag.Render(item.Subject) + "  " + rowStyle.Render(item.Text))
		} else {
			b.WriteString(rowStyle.Render(fmt.Sprintf("%s %s", box, item.Text)))
		}
	}
	return b.String()
}

func RenderProjects(items []models.ProjectItem, t style.Theme) string {
	var b strings.Builder
	for i, item := range items {
		if i > 0 {
			b.WriteString("\n")
		}
		name := t.Text.Render("- " + item.Name)
		meta := t.ProjectDetail.Render("  " + item.Subject + " - Due: " + item.Due)
		b.WriteString(name)
		b.WriteString("\n")
		b.WriteString(meta)
	}
	return b.String()
}
