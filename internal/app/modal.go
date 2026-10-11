package app

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/OmarIbrahim22/seman/internal/models"
	"github.com/OmarIbrahim22/seman/internal/style"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type modalKind int

const (
	modalNone modalKind = iota
	modalAddSubject
	modalAddExam
	modalAddProject
	modalAddTodo
	modalBulkAddTodo
	modalEditSubject
	modalEditExam
	modalEditProject
	modalEditTodo
	modalEditLofiURL
	modalFilterExam
	modalSubjectFilter
	modalConfirm
)

type confirmKind int

const (
	confirmDeleteSubject confirmKind = iota
	confirmDeleteExam
	confirmDeleteProject
	confirmDeleteTodo
	confirmClearAll
)

type confirmAction struct {
	kind       confirmKind
	subjectIdx int
	examIdx    int
	projectIdx int
}

type formField struct {
	label    string
	input    textinput.Model
	required bool
}

// isSubjectField reports whether the currently focused form field is a Subject field.
func (m *Model) isSubjectField() bool {
	if m.formFocus < 0 || m.formFocus >= len(m.formFields) {
		return false
	}
	return strings.EqualFold(m.formFields[m.formFocus].label, "Subject")
}

// refreshDropdown recomputes dropdown matches for the active Subject field.
// Matches are sorted: starts-with first, then contains; both code and name are searched.
func (m *Model) refreshDropdown() {
	if !m.isSubjectField() {
		m.dropdownMatches = nil
		m.dropdownCursor = -1
		return
	}
	query := strings.ToLower(strings.TrimSpace(m.formFields[m.formFocus].input.Value()))
	var starts, contains []string
	for _, s := range m.subjects {
		code := strings.ToLower(s.Code)
		name := strings.ToLower(s.Name)
		if query == "" {
			starts = append(starts, s.Code)
		} else if strings.HasPrefix(code, query) || strings.HasPrefix(name, query) {
			starts = append(starts, s.Code)
		} else if strings.Contains(code, query) || strings.Contains(name, query) {
			contains = append(contains, s.Code)
		}
	}
	m.dropdownMatches = append(starts, contains...)
	if len(m.dropdownMatches) == 0 {
		m.dropdownCursor = -1
	} else if m.dropdownCursor >= len(m.dropdownMatches) {
		m.dropdownCursor = len(m.dropdownMatches) - 1
	} else if m.dropdownCursor < 0 {
		m.dropdownCursor = 0
	}
}

// applyDropdownSelection fills the active field with the currently highlighted match.
func (m *Model) applyDropdownSelection() {
	if m.dropdownCursor < 0 || m.dropdownCursor >= len(m.dropdownMatches) {
		return
	}
	m.formFields[m.formFocus].input.SetValue(m.dropdownMatches[m.dropdownCursor])
	m.formFields[m.formFocus].input.CursorEnd()
}

func (m Model) updateModal(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.modal == modalSubjectFilter {
		if key, ok := msg.(tea.KeyMsg); ok {
			return m.updateSubjectFilterModal(key)
		}
		return m, nil
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		switch key {
		case "esc":
			m.closeModal()
			return m, nil
		case "enter":
			if m.modal == modalConfirm {
				m.pushUndo()
				m.applyConfirmAction()
				m.closeModal()
				return m, nil
			}
			if m.formFocus == len(m.formFields)-1 {
				m.pushUndo()
				if err := m.submitForm(); err != nil {
					m.undoStack = m.undoStack[:len(m.undoStack)-1]
					m.modalError = err.Error()
					return m, nil
				}
				cmd := m.consumeLofiReload()
				m.closeModal()
				return m, cmd
			}
			m.setFormFocus(m.formFocus + 1)
			m.refreshDropdown()
			return m, nil
		case "up":
			if m.isSubjectField() && len(m.dropdownMatches) > 0 {
				if m.dropdownCursor > 0 {
					m.dropdownCursor--
				}
				m.applyDropdownSelection()
				return m, nil
			}
		case "down":
			if m.isSubjectField() && len(m.dropdownMatches) > 0 {
				if m.dropdownCursor < len(m.dropdownMatches)-1 {
					m.dropdownCursor++
				}
				m.applyDropdownSelection()
				return m, nil
			}
		case "tab":
			if m.modal != modalConfirm {
				if m.isSubjectField() && len(m.dropdownMatches) > 0 {
					m.dropdownCursor = (m.dropdownCursor + 1) % len(m.dropdownMatches)
					m.applyDropdownSelection()
					return m, nil
				}
				m.setFormFocus((m.formFocus + 1) % len(m.formFields))
				m.refreshDropdown()
				return m, nil
			}
		case "shift+tab":
			if m.modal != modalConfirm {
				if m.isSubjectField() && len(m.dropdownMatches) > 0 {
					m.dropdownCursor = (m.dropdownCursor - 1 + len(m.dropdownMatches)) % len(m.dropdownMatches)
					m.applyDropdownSelection()
					return m, nil
				}
				next := m.formFocus - 1
				if next < 0 {
					next = len(m.formFields) - 1
				}
				m.setFormFocus(next)
				m.refreshDropdown()
				return m, nil
			}
		case "y", "Y":
			if m.modal == modalConfirm {
				m.applyConfirmAction()
				m.closeModal()
				return m, nil
			}
		case "n", "N":
			if m.modal == modalConfirm {
				m.closeModal()
				return m, nil
			}
		}
	}

	if m.modal == modalConfirm || len(m.formFields) == 0 {
		return m, nil
	}

	field := m.formFields[m.formFocus]
	var cmd tea.Cmd
	field.input, cmd = field.input.Update(msg)
	m.formFields[m.formFocus] = field
	m.refreshDropdown()
	return m, cmd
}

func (m *Model) closeModal() {
	m.modal = modalNone
	m.formFields = nil
	m.formFocus = 0
	m.modalTitle = ""
	m.modalHint = ""
	m.modalError = ""
	m.editSubjectIdx = -1
	m.editExamIdx = -1
	m.editProjectIdx = -1
	m.editTodoIdx = -1
	m.dropdownMatches = nil
	m.dropdownCursor = -1
	m.filterModalActive = nil
	m.filterModalCursor = 0
}

func (m Model) updateSubjectFilterModal(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case "esc":
		m.closeModal()
	case "enter":
		m.applySubjectFilter()
		m.closeModal()
	case "c", "C":
		m.filterModalActive = make([]bool, len(m.subjects))
	case "j", "down":
		if m.filterModalCursor < len(m.subjects)-1 {
			m.filterModalCursor++
		}
	case "k", "up":
		if m.filterModalCursor > 0 {
			m.filterModalCursor--
		}
	case " ":
		if m.filterModalCursor >= 0 && m.filterModalCursor < len(m.filterModalActive) {
			m.filterModalActive[m.filterModalCursor] = !m.filterModalActive[m.filterModalCursor]
		}
	}
	return m, nil
}

func (m *Model) applySubjectFilter() {
	var filters []string
	for i, active := range m.filterModalActive {
		if active && i < len(m.subjects) {
			filters = append(filters, m.subjects[i].Code)
		}
	}
	m.subjectFilters = filters
	m.refreshAllFilters()
}

func (m *Model) setFormFocus(idx int) {
	if len(m.formFields) == 0 {
		return
	}
	for i := range m.formFields {
		if i == idx {
			m.formFields[i].input.Focus()
		} else {
			m.formFields[i].input.Blur()
		}
	}
	m.formFocus = idx
}

func (m *Model) openAddSubject() {
	inputWidth := m.modalInputWidth()
	fields := []formField{
		newFormField("Code", inputWidth, true),
		newFormField("Name", inputWidth, true),
	}
	m.openFormModal(modalAddSubject, "Add Subject", fields)
}

func (m *Model) openAddExam() {
	m.openAddExamWithFilter()
}

func (m *Model) openAddExamWithFilter() {
	inputWidth := m.modalInputWidth()
	fields := []formField{
		newFormField("Subject", inputWidth, true),
		newFormField("Exam Name", inputWidth, true),
		newFormField("Date", inputWidth, true),
		newFormField("Retakes", inputWidth, false),
		newFormField("Priority", inputWidth, false),
	}
	if m.activeTab == tabExams && len(m.subjectFilters) == 1 {
		fields[0].input.SetValue(m.subjectFilters[0])
	}
	fields[2].input.Placeholder = "Jan 2, 2006 @ 15:04"
	fields[3].input.Placeholder = "Jan 5, 2006, Jan 8, 2006"
	fields[4].input.Placeholder = "HIGH / MED / LOW"
	m.openFormModal(modalAddExam, "Add Exam", fields)
}

func (m *Model) openSubjectFilter() {
	m.filterModalActive = make([]bool, len(m.subjects))
	for i, s := range m.subjects {
		m.filterModalActive[i] = m.isFiltered(s.Code)
	}
	m.filterModalCursor = 0
	m.modal = modalSubjectFilter
	m.modalTitle = "Filter by Subject"
	m.modalHint = "↑↓ navigate · Space toggle · Enter apply · C clear all · Esc cancel"
	m.modalError = ""
	m.formFields = nil
}

func (m *Model) openAddProject() {
	inputWidth := m.modalInputWidth()
	fields := []formField{
		newFormField("Name", inputWidth, true),
		newFormField("Subject", inputWidth, true),
		newFormField("Deadline", inputWidth, true),
		newFormField("Status", inputWidth, false),
	}
	m.openFormModal(modalAddProject, "Add Project", fields)
}

func (m *Model) openAddTodo() {
	inputWidth := m.modalInputWidth()
	fields := []formField{
		newFormField("Task", inputWidth, true),
		newFormField("Subject", inputWidth, false),
	}
	m.openFormModal(modalAddTodo, "Add Todo", fields)
}

func (m *Model) openBulkAddTodo() {
	w := m.modalInputWidth()
	fields := []formField{
		newFormField("Task Name", w, true),
		newFormField("Subject", w, false),
		newFormField("Weeks", w, true),
		newFormField("Start #", w, false),
	}
	fields[2].input.Placeholder = "e.g. 12"
	fields[3].input.Placeholder = "1"
	m.openFormModal(modalBulkAddTodo, "Bulk Add Todos", fields)
}

func (m *Model) openEditCurrent() {
	switch m.activeTab {
	case tabSubjects:
		m.openEditSubject()
	case tabExams:
		m.openEditExam()
	case tabTodos:
		m.openEditTodo()
	case tabProjects:
		m.openEditProject()
	}
}

func (m *Model) openEditSubject() {
	if m.selectedSubj < 0 || m.selectedSubj >= len(m.subjects) {
		return
	}
	subj := m.subjects[m.selectedSubj]
	inputWidth := m.modalInputWidth()
	fields := []formField{
		newFormField("Code", inputWidth, true),
		newFormField("Name", inputWidth, true),
	}
	fields[0].input.SetValue(subj.Code)
	fields[1].input.SetValue(subj.Name)
	m.editSubjectIdx = m.selectedSubj
	m.openFormModal(modalEditSubject, "Edit Subject", fields)
}

func (m *Model) openEditExam() {
	if m.examCursor < 0 || m.examCursor >= len(m.flatExams) {
		return
	}
	flat := m.flatExams[m.examCursor]
	exam := flat.Exam
	inputWidth := m.modalInputWidth()
	fields := []formField{
		newFormField("Exam Name", inputWidth, true),
		newFormField("Date", inputWidth, true),
		newFormField("Retakes", inputWidth, false),
		newFormField("Priority", inputWidth, false),
	}
	fields[0].input.SetValue(exam.Name)
	fields[1].input.SetValue(exam.Date)
	fields[2].input.SetValue(strings.Join(exam.Retakes, ", "))
	fields[3].input.SetValue(exam.Priority)
	m.editSubjectIdx = flat.SubjectIdx
	m.editExamIdx = flat.ExamIdx
	m.openFormModal(modalEditExam, "Edit Exam ("+flat.SubjectCode+")", fields)
}

func (m *Model) openEditProject() {
	if m.projectCursor < 0 || m.projectCursor >= len(m.projects) {
		return
	}
	project := m.projects[m.projectCursor]
	inputWidth := m.modalInputWidth()
	fields := []formField{
		newFormField("Name", inputWidth, true),
		newFormField("Subject", inputWidth, true),
		newFormField("Deadline", inputWidth, true),
		newFormField("Status", inputWidth, false),
	}
	fields[0].input.SetValue(project.Name)
	fields[1].input.SetValue(project.Subject)
	fields[2].input.SetValue(project.Due)
	fields[3].input.SetValue(project.Status)
	m.editProjectIdx = m.projectCursor
	m.openFormModal(modalEditProject, "Edit Project", fields)
}

func (m *Model) openEditTodo() {
	if m.checklistCursor < 0 || m.checklistCursor >= len(m.checklistItems) {
		return
	}
	item := m.checklistItems[m.checklistCursor]
	inputWidth := m.modalInputWidth()
	fields := []formField{
		newFormField("Task", inputWidth, true),
		newFormField("Subject", inputWidth, false),
	}
	fields[0].input.SetValue(item.Text)
	fields[1].input.SetValue(item.Subject)
	m.editTodoIdx = m.checklistCursor
	m.openFormModal(modalEditTodo, "Edit Todo", fields)
}

func (m *Model) openEditLofiURL() {
	inputWidth := m.modalInputWidth()
	fields := []formField{
		newFormField("Playlist URL", inputWidth, true),
	}
	fields[0].input.SetValue(m.lofi.url)
	m.openFormModal(modalEditLofiURL, "Edit Lofi Playlist", fields)
}

func (m *Model) openFormModal(kind modalKind, title string, fields []formField) {
	t := style.ThemeOf(m.themeName)
	for i := range fields {
		fields[i].input.TextStyle = t.InputText
		fields[i].input.PlaceholderStyle = t.InputHint
		fields[i].input.CursorStyle = t.InputCursor
	}
	m.modal = kind
	m.formFields = fields
	m.modalTitle = title
	hasSubjectField := false
	for _, f := range fields {
		if strings.EqualFold(f.label, "Subject") {
			hasSubjectField = true
			break
		}
	}
	if hasSubjectField {
		m.modalHint = "Tab/↑↓ select subject · Enter to advance · Esc to cancel"
	} else {
		m.modalHint = "Tab to switch · Enter to save · Esc to cancel"
	}
	m.modalError = ""
	m.dropdownCursor = 0
	m.setFormFocus(0)
	m.refreshDropdown()
}

func newFormField(label string, width int, required bool) formField {
	input := textinput.New()
	input.Prompt = ""
	input.Width = width
	input.CharLimit = 200
	t := style.NewTheme()
	input.TextStyle = t.InputText
	input.PlaceholderStyle = t.InputHint
	input.CursorStyle = t.InputCursor
	return formField{
		label:    label,
		input:    input,
		required: required,
	}
}

func (m *Model) submitForm() error {
	switch m.modal {
	case modalAddSubject:
		code := strings.TrimSpace(m.formFields[0].input.Value())
		name := strings.TrimSpace(m.formFields[1].input.Value())
		if code == "" || name == "" {
			return fmt.Errorf("Code and Name are required.")
		}
		m.subjects = append(m.subjects, models.SubjectItem{Code: code, Name: name})
		m.selectedSubj = len(m.subjects) - 1
		m.persist()
	case modalAddExam:
		subjectCode := strings.TrimSpace(m.formFields[0].input.Value())
		examName := strings.TrimSpace(m.formFields[1].input.Value())
		date := strings.TrimSpace(m.formFields[2].input.Value())
		retakesRaw := strings.TrimSpace(m.formFields[3].input.Value())
		priority := strings.TrimSpace(m.formFields[4].input.Value())
		if subjectCode == "" || examName == "" || date == "" {
			return fmt.Errorf("Subject, Exam Name, and Date are required.")
		}
		idx := findSubjectIndex(m.subjects, subjectCode)
		if idx < 0 {
			return fmt.Errorf("Subject code not found.")
		}
		retakes := splitCSV(retakesRaw)
		m.subjects[idx].Exams = append(m.subjects[idx].Exams, models.ExamItem{
			Name:     examName,
			Date:     date,
			Retakes:  retakes,
			Priority: strings.ToUpper(priority),
		})
		newExamIdx := len(m.subjects[idx].Exams) - 1
		m.refreshFlatExams()
		for i, flat := range m.flatExams {
			if flat.SubjectIdx == idx && flat.ExamIdx == newExamIdx {
				m.examCursor = i
				break
			}
		}
		m.persist()
	case modalAddProject:
		name := strings.TrimSpace(m.formFields[0].input.Value())
		subject := strings.TrimSpace(m.formFields[1].input.Value())
		deadline := strings.TrimSpace(m.formFields[2].input.Value())
		status := strings.TrimSpace(m.formFields[3].input.Value())
		if name == "" || subject == "" || deadline == "" {
			return fmt.Errorf("Name, Subject, and Deadline are required.")
		}
		if status == "" {
			status = "NOT STARTED"
		}
		m.projects = append(m.projects, models.ProjectItem{
			Name:    name,
			Subject: subject,
			Due:     deadline,
			Status:  strings.ToUpper(status),
		})
		m.projectCursor = len(m.projects) - 1
		m.sortProjectsByStatus()
		m.refreshProjectFilter()
		m.persist()
	case modalAddTodo:
		task := strings.TrimSpace(m.formFields[0].input.Value())
		subject := strings.TrimSpace(m.formFields[1].input.Value())
		if task == "" {
			return fmt.Errorf("Task is required.")
		}
		if subject != "" && findSubjectIndex(m.subjects, subject) < 0 {
			return fmt.Errorf("Subject not found.")
		}
		subjectCode := strings.ToUpper(subject)
		m.checklistItems = append(m.checklistItems, models.ChecklistItem{
			Text:    task,
			Done:    false,
			Due:     m.weekStart.Format("2006-01-02"),
			Subject: subjectCode,
		})
		m.checklistCursor = len(m.checklistItems) - 1
		m.sortChecklistByDone()
		m.persist()
		m.refreshChecklistView()
	case modalBulkAddTodo:
		baseName := strings.TrimSpace(m.formFields[0].input.Value())
		subject := strings.TrimSpace(m.formFields[1].input.Value())
		weeksStr := strings.TrimSpace(m.formFields[2].input.Value())
		startStr := strings.TrimSpace(m.formFields[3].input.Value())
		if baseName == "" {
			return fmt.Errorf("Task Name is required.")
		}
		weeks, err := strconv.Atoi(weeksStr)
		if err != nil || weeks < 1 {
			return fmt.Errorf("Weeks must be a positive number.")
		}
		if weeks > 52 {
			return fmt.Errorf("Weeks must be 52 or fewer.")
		}
		startNum := 1
		if startStr != "" {
			n, err := strconv.Atoi(startStr)
			if err != nil || n < 1 {
				return fmt.Errorf("Start # must be a positive number.")
			}
			startNum = n
		}
		if subject != "" && findSubjectIndex(m.subjects, subject) < 0 {
			return fmt.Errorf("Subject not found.")
		}
		subjectCode := strings.ToUpper(subject)
		for i := 0; i < weeks; i++ {
			due := m.weekStart.AddDate(0, 0, i*7)
			m.checklistItems = append(m.checklistItems, models.ChecklistItem{
				Text:    fmt.Sprintf("%s - Week %d", baseName, startNum+i),
				Done:    false,
				Due:     due.Format("2006-01-02"),
				Subject: subjectCode,
			})
		}
		m.checklistCursor = len(m.checklistItems) - 1
		m.sortChecklistByDone()
		m.persist()
		m.refreshChecklistView()
	case modalEditSubject:
		if m.editSubjectIdx < 0 || m.editSubjectIdx >= len(m.subjects) {
			return nil
		}
		code := strings.TrimSpace(m.formFields[0].input.Value())
		name := strings.TrimSpace(m.formFields[1].input.Value())
		if code == "" || name == "" {
			return fmt.Errorf("Code and Name are required.")
		}
		m.subjects[m.editSubjectIdx].Code = code
		m.subjects[m.editSubjectIdx].Name = name
		m.persist()
	case modalEditExam:
		if m.editSubjectIdx < 0 || m.editSubjectIdx >= len(m.subjects) {
			return nil
		}
		exams := m.subjects[m.editSubjectIdx].Exams
		if m.editExamIdx < 0 || m.editExamIdx >= len(exams) {
			return nil
		}
		examName := strings.TrimSpace(m.formFields[0].input.Value())
		date := strings.TrimSpace(m.formFields[1].input.Value())
		retakesRaw := strings.TrimSpace(m.formFields[2].input.Value())
		priority := strings.TrimSpace(m.formFields[3].input.Value())
		if examName == "" || date == "" {
			return fmt.Errorf("Exam Name and Date are required.")
		}
		exams[m.editExamIdx].Name = examName
		exams[m.editExamIdx].Date = date
		exams[m.editExamIdx].Retakes = splitCSV(retakesRaw)
		exams[m.editExamIdx].Priority = strings.ToUpper(priority)
		m.subjects[m.editSubjectIdx].Exams = exams
		editedSI := m.editSubjectIdx
		editedEI := m.editExamIdx
		m.refreshFlatExams()
		for i, flat := range m.flatExams {
			if flat.SubjectIdx == editedSI && flat.ExamIdx == editedEI {
				m.examCursor = i
				break
			}
		}
		m.persist()
	case modalEditProject:
		if m.editProjectIdx < 0 || m.editProjectIdx >= len(m.projects) {
			return nil
		}
		name := strings.TrimSpace(m.formFields[0].input.Value())
		subject := strings.TrimSpace(m.formFields[1].input.Value())
		deadline := strings.TrimSpace(m.formFields[2].input.Value())
		status := strings.TrimSpace(m.formFields[3].input.Value())
		if name == "" || subject == "" || deadline == "" {
			return fmt.Errorf("Name, Subject, and Deadline are required.")
		}
		if status == "" {
			status = "NOT STARTED"
		}
		m.projects[m.editProjectIdx].Name = name
		m.projects[m.editProjectIdx].Subject = subject
		m.projects[m.editProjectIdx].Due = deadline
		m.projects[m.editProjectIdx].Status = strings.ToUpper(status)
		m.projectCursor = m.editProjectIdx
		m.sortProjectsByStatus()
		m.refreshProjectFilter()
		m.persist()
	case modalEditTodo:
		if m.editTodoIdx < 0 || m.editTodoIdx >= len(m.checklistItems) {
			return nil
		}
		task := strings.TrimSpace(m.formFields[0].input.Value())
		subject := strings.TrimSpace(m.formFields[1].input.Value())
		if task == "" {
			return fmt.Errorf("Task is required.")
		}
		if subject != "" && findSubjectIndex(m.subjects, subject) < 0 {
			return fmt.Errorf("Subject not found.")
		}
		m.checklistItems[m.editTodoIdx].Text = task
		m.checklistItems[m.editTodoIdx].Subject = strings.ToUpper(subject)
		m.sortChecklistByDone()
		m.persist()
		m.refreshChecklistView()
	case modalFilterExam:
		// legacy: handled by modalSubjectFilter now; kept for safety
	case modalEditLofiURL:
		url := strings.TrimSpace(m.formFields[0].input.Value())
		if url == "" {
			return fmt.Errorf("Playlist URL is required.")
		}
		m.lofi.url = url
		m.lofiReload = true
		m.persist()
		if m.lofi.cmd != nil && m.lofi.status != lofiStatusStopped {
			_ = m.sendLofiCommand("loadfile", m.lofi.url, "replace")
			m.lofi.status = lofiStatusPlaying
		}
	}
	return nil
}

func findSubjectIndex(items []models.SubjectItem, code string) int {
	code = strings.ToUpper(strings.TrimSpace(code))
	for i, item := range items {
		if strings.ToUpper(item.Code) == code {
			return i
		}
	}
	return -1
}

func splitCSV(raw string) []string {
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		val := strings.TrimSpace(part)
		if val != "" {
			out = append(out, val)
		}
	}
	return out
}

func (m *Model) queueDelete() {
	switch m.activeTab {
	case tabSubjects:
		if len(m.subjects) == 0 {
			return
		}
		action := confirmAction{kind: confirmDeleteSubject, subjectIdx: m.selectedSubj}
		message := fmt.Sprintf("Delete subject %s and its exams?", m.subjects[m.selectedSubj].Code)
		m.confirmOrApply(action, message)
	case tabExams:
		if len(m.flatExams) == 0 || m.examCursor < 0 || m.examCursor >= len(m.flatExams) {
			return
		}
		flat := m.flatExams[m.examCursor]
		action := confirmAction{kind: confirmDeleteExam, subjectIdx: flat.SubjectIdx, examIdx: flat.ExamIdx}
		message := fmt.Sprintf("Delete exam \"%s\" (%s)?", flat.Exam.Name, flat.SubjectCode)
		m.confirmOrApply(action, message)
	case tabProjects:
		if len(m.projects) == 0 {
			return
		}
		action := confirmAction{kind: confirmDeleteProject, projectIdx: m.projectCursor}
		message := fmt.Sprintf("Delete project %s?", m.projects[m.projectCursor].Name)
		m.confirmOrApply(action, message)
	case tabTodos:
		if len(m.checklistItems) == 0 || m.checklistCursor < 0 || m.checklistCursor >= len(m.checklistItems) {
			return
		}
		action := confirmAction{kind: confirmDeleteTodo, projectIdx: m.checklistCursor}
		message := fmt.Sprintf("Delete task \"%s\"?", m.checklistItems[m.checklistCursor].Text)
		m.confirmOrApply(action, message)
	}
}

func (m *Model) queueClearAll() {
	action := confirmAction{kind: confirmClearAll}
	message := "Clear all data? This cannot be undone."
	m.confirmOrApply(action, message)
}

func (m *Model) confirmOrApply(action confirmAction, message string) {
	if m.confirmOn {
		m.modal = modalConfirm
		m.confirmAction = action
		m.modalTitle = "Confirm"
		m.modalHint = "[Y] Confirm  [N] Cancel"
		m.modalError = message
		return
	}
	m.confirmAction = action
	m.applyConfirmAction()
}

func (m *Model) applyConfirmAction() {
	switch m.confirmAction.kind {
	case confirmDeleteExam:
		si := m.confirmAction.subjectIdx
		ei := m.confirmAction.examIdx
		if si >= 0 && si < len(m.subjects) {
			exams := m.subjects[si].Exams
			if ei >= 0 && ei < len(exams) {
				m.subjects[si].Exams = append(exams[:ei], exams[ei+1:]...)
			}
		}
		m.refreshFlatExams()
	case confirmDeleteSubject:
		if m.confirmAction.subjectIdx >= 0 && m.confirmAction.subjectIdx < len(m.subjects) {
			m.subjects = append(m.subjects[:m.confirmAction.subjectIdx], m.subjects[m.confirmAction.subjectIdx+1:]...)
			if m.selectedSubj >= len(m.subjects) {
				m.selectedSubj = len(m.subjects) - 1
			}
			if m.selectedSubj < 0 {
				m.selectedSubj = 0
			}
		}
	case confirmDeleteProject:
		if m.confirmAction.projectIdx >= 0 && m.confirmAction.projectIdx < len(m.projects) {
			m.projects = append(m.projects[:m.confirmAction.projectIdx], m.projects[m.confirmAction.projectIdx+1:]...)
			if m.projectCursor >= len(m.projects) {
				m.projectCursor = len(m.projects) - 1
			}
			if m.projectCursor < 0 {
				m.projectCursor = 0
			}
			m.refreshProjectFilter()
		}
	case confirmDeleteTodo:
		if m.confirmAction.projectIdx >= 0 && m.confirmAction.projectIdx < len(m.checklistItems) {
			idx := m.confirmAction.projectIdx
			m.checklistItems = append(m.checklistItems[:idx], m.checklistItems[idx+1:]...)
			if m.checklistCursor >= len(m.checklistItems) {
				m.checklistCursor = len(m.checklistItems) - 1
			}
			if m.checklistCursor < 0 {
				m.checklistCursor = 0
			}
			m.refreshChecklistView()
		}
	case confirmClearAll:
		m.subjects = nil
		m.projects = nil
		m.checklistItems = nil
		m.weeklyExams = nil
		m.selectedSubj = 0
		m.projectCursor = 0
		m.projectVisible = nil
		m.refreshAllFilters()
	}
	m.persist()
}

func (m Model) modalInputWidth() int {
	modalW := minInt(70, m.width-6)
	if modalW < 42 {
		modalW = 42
	}
	return maxInt(12, modalW-18)
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
