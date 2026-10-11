package app

import "github.com/OmarIbrahim22/seman/internal/models"

const maxUndoStack = 50

type undoSnapshot struct {
	subjects       []models.SubjectItem
	projects       []models.ProjectItem
	checklistItems []models.ChecklistItem
	weeklyExams    []string
}

func (m *Model) pushUndo() {
	snap := undoSnapshot{
		checklistItems: make([]models.ChecklistItem, len(m.checklistItems)),
		projects:       make([]models.ProjectItem, len(m.projects)),
		weeklyExams:    make([]string, len(m.weeklyExams)),
		subjects:       make([]models.SubjectItem, len(m.subjects)),
	}
	copy(snap.checklistItems, m.checklistItems)
	copy(snap.projects, m.projects)
	copy(snap.weeklyExams, m.weeklyExams)
	for i, s := range m.subjects {
		sc := s
		sc.Exams = make([]models.ExamItem, len(s.Exams))
		for j, e := range s.Exams {
			ec := e
			if len(e.Retakes) > 0 {
				ec.Retakes = make([]string, len(e.Retakes))
				copy(ec.Retakes, e.Retakes)
			}
			sc.Exams[j] = ec
		}
		snap.subjects[i] = sc
	}
	if len(m.undoStack) >= maxUndoStack {
		m.undoStack = m.undoStack[1:]
	}
	m.undoStack = append(m.undoStack, snap)
}

func (m *Model) undo() {
	if len(m.undoStack) == 0 {
		return
	}
	snap := m.undoStack[len(m.undoStack)-1]
	m.undoStack = m.undoStack[:len(m.undoStack)-1]

	m.subjects = snap.subjects
	m.projects = snap.projects
	m.checklistItems = snap.checklistItems
	m.weeklyExams = snap.weeklyExams

	if m.selectedSubj >= len(m.subjects) {
		m.selectedSubj = len(m.subjects) - 1
	}
	if m.selectedSubj < 0 {
		m.selectedSubj = 0
	}
	if m.projectCursor >= len(m.projects) {
		m.projectCursor = len(m.projects) - 1
	}
	if m.projectCursor < 0 {
		m.projectCursor = 0
	}
	if m.checklistCursor >= len(m.checklistItems) {
		m.checklistCursor = len(m.checklistItems) - 1
	}
	if m.checklistCursor < 0 {
		m.checklistCursor = -1
	}

	m.sortChecklistByDone()
	m.sortProjectsByStatus()
	m.refreshFlatExams()
	m.refreshAllFilters()
	m.persist()
}
