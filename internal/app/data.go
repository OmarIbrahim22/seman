package app

import (
	"github.com/OmarIbrahim22/seman/internal/models"
	"github.com/OmarIbrahim22/seman/internal/storage"
)

func DefaultData() storage.SemesterData {
	return storage.SemesterData{
		Subjects: []models.SubjectItem{
			{
				Code: "PHY150",
				Name: "Physics",
				Exams: []models.ExamItem{
					{
						Name:     "Test",
						Date:     "Jan 14, 2026 @ 10:00",
						Retakes:  []string{"Jan 21, 2026 @ 10:00"},
						Priority: "MED",
					},
					{
						Name:     "Final Exam",
						Date:     "Jan 17, 2026 @ 09:00",
						Retakes:  []string{"Jan 24, 2026 @ 09:00"},
						Priority: "HIGH",
					},
				},
			},
			{
				Code: "HIST210",
				Name: "History",
				Exams: []models.ExamItem{
					{
						Name:     "Oral Exam",
						Date:     "Jan 15, 2026 @ 14:00",
						Retakes:  []string{"Jan 22, 2026 @ 14:00"},
						Priority: "MED",
					},
				},
			},
			{
				Code: "ENG102",
				Name: "English",
				Exams: []models.ExamItem{
					{
						Name:     "Listening Test",
						Date:     "Jan 16, 2026 @ 11:30",
						Retakes:  []string{"Jan 23, 2026 @ 11:30"},
						Priority: "LOW",
					},
				},
			},
			{
				Code: "CS101",
				Name: "Computer Science I",
				Exams: []models.ExamItem{
					{
						Name:     "Final Exam",
						Date:     "Jan 13, 2026 @ 13:00",
						Retakes:  []string{"Jan 20, 2026 @ 13:00"},
						Priority: "HIGH",
					},
				},
			},
			{Code: "MATH201", Name: "Calculus II"},
		},
		Projects: []models.ProjectItem{
			{Name: "Database Design Project", Subject: "CS101", Due: "May 20, 2025", Status: "IN PROGRESS"},
			{Name: "Renaissance Art Essay", Subject: "HIST210", Due: "May 28, 2025", Status: "NOT STARTED"},
			{Name: "Physics Lab Simulation", Subject: "PHY150", Due: "Jun 3, 2025", Status: "IN PROGRESS"},
			{Name: "Calculus Portfolio", Subject: "MATH201", Due: "Jun 10, 2025", Status: "NOT STARTED"},
			{Name: "English Literature Review", Subject: "ENG102", Due: "Jun 15, 2025", Status: "DONE"},
			{Name: "Machine Learning Assignment", Subject: "CS101", Due: "Jun 18, 2025", Status: "NOT STARTED"},
		},
		Checklist: []models.ChecklistItem{
			{Text: "Create flashcards for history", Done: true, Subject: "HIST210"},
			{Text: "Attend study group session", Done: true},
			{Text: "Mock exam practice", Done: false, Subject: "CS101"},
			{Text: "Summarize key formulas", Done: false, Subject: "MATH201"},
			{Text: "Review lectures 8-10", Done: false, Subject: "CS101"},
			{Text: "Complete problem set 5", Done: true, Subject: "MATH201"},
			{Text: "Read chapter 12-13", Done: false, Subject: "PHY150"},
			{Text: "Outline essay", Done: false, Subject: "HIST210"},
			{Text: "Practice pronunciation", Done: true, Subject: "ENG102"},
			{Text: "Work on database project", Done: false, Subject: "CS101"},
			{Text: "Group meeting for CS project", Done: false, Subject: "CS101"},
			{Text: "Submit lab report", Done: true, Subject: "PHY150"},
			{Text: "Study calculus theorems", Done: false, Subject: "MATH201"},
			{Text: "Review weak areas", Done: false, Subject: "PHY150"},
			{Text: "Prepare oral exam notes", Done: false, Subject: "HIST210"},
		},
		WeeklyExams: []string{},
		ConfirmOn:   true,
		WeekStart:   "2026-01-12",
		WeekSpan:    1,
		LofiEnabled: false,
		LofiURL:     "",
	}
}
