//go:build js && wasm

package main

import (
	"time"

	"github.com/OmarIbrahim22/seman/internal/models"
	"github.com/OmarIbrahim22/seman/internal/storage"
)

// demoData is a sample semester laid out around today's date, so a first-time
// visitor sees a populated week instead of an empty app.
func demoData() storage.SemesterData {
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	monday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, 1-weekday)

	exam := func(days int, clock string) string {
		return monday.AddDate(0, 0, days).Format("Jan 2, 2006") + " @ " + clock
	}
	due := func(days int) string { return monday.AddDate(0, 0, days).Format("2006-01-02") }
	day := func(days int) string { return monday.AddDate(0, 0, days).Format("Jan 2, 2006") }

	return storage.SemesterData{
		Subjects: []models.SubjectItem{
			{Code: "IB002", Name: "Algorithms and Data Structures", Exams: []models.ExamItem{
				{Name: "Midterm", Date: exam(2, "10:00"), Retakes: []string{exam(9, "10:00")}, Priority: "HIGH"},
				{Name: "Final Exam", Date: exam(23, "09:00"), Retakes: []string{exam(30, "09:00")}, Priority: "HIGH"},
			}},
			{Code: "PB138", Name: "Web Development", Exams: []models.ExamItem{
				{Name: "Project Defence", Date: exam(4, "14:00"), Priority: "MED"},
			}},
			{Code: "PV080", Name: "Information Security", Exams: []models.ExamItem{
				{Name: "Written Test", Date: exam(8, "12:00"), Retakes: []string{exam(15, "12:00")}, Priority: "MED"},
			}},
			{Code: "MB143", Name: "Statistics", Exams: []models.ExamItem{
				{Name: "Quiz 3", Date: exam(1, "08:00"), Priority: "LOW"},
				{Name: "Final Exam", Date: exam(17, "13:00"), Retakes: []string{exam(24, "13:00")}, Priority: "HIGH"},
			}},
			{Code: "VB000", Name: "Academic Writing"},
		},
		Projects: []models.ProjectItem{
			{Name: "Study planner web app", Subject: "PB138", Due: day(4), Status: "IN PROGRESS"},
			{Name: "Graph algorithms homework", Subject: "IB002", Due: day(6), Status: "NOT STARTED"},
			{Name: "Buggy code security review", Subject: "PV080", Due: day(11), Status: "IN PROGRESS"},
			{Name: "Regression analysis report", Subject: "MB143", Due: day(13), Status: "NOT STARTED"},
			{Name: "Essay draft", Subject: "VB000", Due: day(-3), Status: "DONE"},
		},
		Checklist: []models.ChecklistItem{
			{Text: "Revise sorting and heaps", Done: true, Due: due(0), Subject: "IB002"},
			{Text: "Solve last year's midterm", Done: false, Due: due(1), Subject: "IB002"},
			{Text: "Finish the email verification flow", Done: false, Due: due(2), Subject: "PB138"},
			{Text: "Write Storybook docs for the note list", Done: true, Due: due(1), Subject: "PB138"},
			{Text: "Prepare slides for the defence", Done: false, Due: due(3), Subject: "PB138"},
			{Text: "Read the chapter on hypothesis testing", Done: false, Due: due(0), Subject: "MB143"},
			{Text: "Run the static analyser on the buggy code", Done: false, Due: due(4), Subject: "PV080"},
			{Text: "Send the essay draft for feedback", Done: true, Due: due(2), Subject: "VB000"},
			{Text: "Book a study room for Friday", Done: false, Due: due(3)},
			{Text: "Go through dynamic programming notes", Done: false, Due: due(8), Subject: "IB002"},
			{Text: "Start the regression report", Done: false, Due: due(9), Subject: "MB143"},
		},
		WeeklyExams: []string{},
		ConfirmOn:   true,
		WeekStart:   monday.Format("2006-01-02"),
		WeekSpan:    1,
		Theme:       "catppuccin",
	}
}
