package app

import (
	"sort"
	"strings"
	"time"

	"github.com/OmarIbrahim22/seman/internal/models"
	"github.com/OmarIbrahim22/seman/internal/storage"
	"github.com/OmarIbrahim22/seman/internal/style"
	"github.com/OmarIbrahim22/seman/internal/ui/components"
	"github.com/OmarIbrahim22/seman/internal/ui/screens"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	width             int
	height            int
	activeTab         int
	checklist         viewport.Model
	checklistItems    []models.ChecklistItem
	checklistCursor   int
	todoVisible       []int
	projects          []models.ProjectItem
	projectVisible    []int
	subjects          []models.SubjectItem
	selectedSubj      int
	examCursor        int
	subjectFilters    []string
	filterModalActive []bool
	filterModalCursor int
	flatExams         []models.FlatExam
	weekLabel         string
	weekStart         time.Time
	weeklyExams       []string
	projectCursor     int
	confirmOn         bool
	weekSpan          int
	previousWeekSpan  int
	lofi              lofiState
	lofiPlaylist      []models.LofiTrack
	lofiCursor        int
	lofiOffset        int
	lofiListHeight    int
	lofiNow           int
	lofiReload        bool
	themeName         string
	modal             modalKind
	formFields        []formField
	formFocus         int
	modalTitle        string
	modalHint         string
	modalError        string
	dropdownMatches   []string
	dropdownCursor    int
	confirmAction     confirmAction
	editSubjectIdx    int
	editExamIdx       int
	editProjectIdx    int
	editTodoIdx       int
	store             storage.Store
	undoStack         []undoSnapshot
}

const (
	tabDashboard = iota
	tabExams
	tabTodos
	tabProjects
	tabSettings
	tabLofi
	tabSubjects
)

func NewModel(store storage.Store, data storage.SemesterData, hasData bool) Model {
	if !hasData {
		data = DefaultData()
	}
	vp := viewport.New(0, 0)
	m := Model{
		activeTab:      tabDashboard,
		checklist:      vp,
		selectedSubj:   0,
		examCursor:     0,
		projectCursor:  0,
		editSubjectIdx: -1,
		editExamIdx:    -1,
		editProjectIdx: -1,
		editTodoIdx:    -1,
		lofiNow:        -1,
		store:          store,
	}
	m.lofiPlaylist = defaultLofiPlaylist()
	m.applyData(data)
	return m
}

func (m Model) Init() tea.Cmd {
	if m.lofi.enabled && strings.TrimSpace(m.lofi.url) != "" {
		return loadLofiPlaylist(m.lofi.url)
	}
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.resize(msg.Width, msg.Height)
		return m, nil
	case lofiExitMsg:
		m.handleLofiExit(msg)
		return m, nil
	case lofiPlaylistMsg:
		m.applyLofiPlaylist(msg)
		return m, nil
	case lofiIndexMsg:
		cmd := m.applyLofiIndex(msg)
		return m, cmd
	case lofiPlaybackMsg:
		cmd := m.applyLofiPlayback(msg)
		return m, cmd
	case tea.KeyMsg:
		key := msg.String()
		if m.modal != modalNone {
			return m.updateModal(msg)
		}
		if len(key) == 1 && key[0] >= '1' && key[0] <= '9' {
			if m.switchToTab(int(key[0] - '1')) {
				return m, nil
			}
		}
		switch key {
		case "ctrl+c", "q":
			m.shutdownLofi()
			return m, tea.Quit
		case "ctrl+z":
			m.undo()
			return m, nil
		case "g", "G":
			m.toggleGlobalView()
			return m, nil
		case "left":
			m.shiftWeek(-1)
			return m, nil
		case "right":
			m.shiftWeek(1)
			return m, nil
		case "a", "A":
			m.openAddExamWithFilter()
			return m, nil
		case "s", "S":
			m.openAddSubject()
			return m, nil
		case "p", "P":
			m.openAddProject()
			return m, nil
		case "e", "E":
			m.openEditCurrent()
			return m, nil
		case "f", "F":
			switch m.activeTab {
			case tabDashboard, tabExams, tabTodos, tabProjects:
				m.openSubjectFilter()
				return m, nil
			}
		case "r", "R":
			switch m.activeTab {
			case tabDashboard, tabExams, tabTodos, tabProjects:
				m.subjectFilters = nil
				m.refreshAllFilters()
				return m, nil
			}
		case "n", "N":
			if m.activeTab == tabTodos {
				m.openAddTodo()
				return m, nil
			}
		case "d", "D":
			m.queueDelete()
			return m, nil
		case "c", "C":
			if m.activeTab == tabSettings {
				m.queueClearAll()
				return m, nil
			}
		case "o", "O":
			if m.activeTab == tabSettings {
				m.confirmOn = !m.confirmOn
				m.persist()
				return m, nil
			}
		case "t", "T":
			if m.activeTab == tabSettings {
				m.cycleTheme()
				return m, nil
			}
			m.jumpToCurrentWeek()
			return m, nil
		case "w", "W":
			if m.activeTab == tabSettings {
				m.cycleWeekSpan()
				return m, nil
			}
		case "l", "L":
			if m.activeTab == tabSettings {
				cmd := m.toggleLofiEnabled()
				return m, cmd
			}
		case "u", "U":
			if m.activeTab == tabSettings || m.activeTab == tabLofi {
				m.openEditLofiURL()
				return m, nil
			}
		}

		if m.activeTab == tabTodos {
			switch key {
			case "j", "down":
				m.moveChecklistCursor(1)
				return m, nil
			case "k", "up":
				m.moveChecklistCursor(-1)
				return m, nil
			case "pgdown":
				m.moveChecklistCursor(m.checklist.Height)
				return m, nil
			case "pgup":
				m.moveChecklistCursor(-m.checklist.Height)
				return m, nil
			case " ", "enter", "x", "X":
				m.toggleChecklistItem()
				return m, nil
			case "b", "B":
				m.openBulkAddTodo()
				return m, nil
			}

			var cmd tea.Cmd
			m.checklist, cmd = m.checklist.Update(msg)
			return m, cmd
		}

		if m.activeTab == tabDashboard {
			switch key {
			case "j", "down":
				m.checklist.LineDown(1)
				return m, nil
			case "k", "up":
				m.checklist.LineUp(1)
				return m, nil
			case "pgdown":
				m.checklist.ViewDown()
				return m, nil
			case "pgup":
				m.checklist.ViewUp()
				return m, nil
			}

			var cmd tea.Cmd
			m.checklist, cmd = m.checklist.Update(msg)
			return m, cmd
		}

		if m.activeTab == tabSubjects {
			switch key {
			case "j", "down":
				if m.selectedSubj < len(m.subjects)-1 {
					m.selectedSubj++
				}
				return m, nil
			case "k", "up":
				if m.selectedSubj > 0 {
					m.selectedSubj--
				}
				return m, nil
			case "enter":
				if len(m.subjects) == 0 {
					m.openAddSubject()
				} else {
					m.openEditSubject()
				}
				return m, nil
			}
		}

		if m.activeTab == tabExams {
			switch key {
			case "j", "down":
				m.moveExamCursor(1)
				return m, nil
			case "k", "up":
				m.moveExamCursor(-1)
				return m, nil
			}
		}

		if m.activeTab == tabProjects {
			switch key {
			case "j", "down":
				m.moveProjectCursor(1)
				return m, nil
			case "k", "up":
				m.moveProjectCursor(-1)
				return m, nil
			}
		}

		if m.activeTab == tabLofi {
			switch key {
			case "j", "down":
				m.moveLofiCursor(1)
				return m, nil
			case "k", "up":
				m.moveLofiCursor(-1)
				return m, nil
			case "enter":
				cmd := m.playLofiAt(m.lofiCursor)
				return m, cmd
			case " ":
				cmd := m.toggleLofiPlayPause()
				return m, cmd
			case "n", "N":
				cmd := m.lofiNext()
				return m, cmd
			case "b", "B":
				cmd := m.lofiPrev()
				return m, cmd
			case "x", "X":
				m.lofiStop()
				return m, nil
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	t := style.ThemeOf(m.themeName)
	header := components.RenderHeader(m.width, t)
	tabs := components.RenderTabs(m.activeTab, m.width, m.weekLabel, m.tabItems(), t)
	divider := components.RenderDivider(m.width, t)
	mainHeight := components.MainAreaHeight(m.height)
	state := m.viewState()
	main := screens.RenderMain(state, m.width, mainHeight, t)
	if state.Modal.Mode != components.ModalHidden {
		modal := components.RenderModalContent(state.Modal, m.width, t)
		main = components.PlaceOverlay(main, modal)
	}
	footer := components.RenderFooter(m.width, len(m.tabItems()), m.activeTab, t)

	return strings.Join([]string{header, tabs, divider, main, divider, footer}, "\n")
}

func (m *Model) resize(width, height int) {
	m.width = width
	m.height = height

	mainHeight := components.MainAreaHeight(height)
	switch m.activeTab {
	case tabDashboard:
		layout := components.ComputeDashboardLayout(width, mainHeight)
		if layout.ChecklistWidth > 0 && layout.ChecklistHeight > 0 {
			m.checklist.Width = layout.ChecklistWidth
			m.checklist.Height = layout.ChecklistHeight
		}
	case tabTodos:
		layout := components.ComputeWeeklyLayout(width, mainHeight)
		if layout.ActionsWidth > 0 && layout.ActionsHeight > 0 {
			m.checklist.Width = layout.ActionsWidth
			m.checklist.Height = layout.ActionsHeight
		}
	case tabLofi:
		m.resizeLofi(width, mainHeight)
	}
	m.refreshChecklistView()
}

func (m *Model) applyData(data storage.SemesterData) {
	m.subjects = data.Subjects
	m.projects = data.Projects
	m.checklistItems = data.Checklist
	m.weeklyExams = data.WeeklyExams
	m.confirmOn = data.ConfirmOn
	m.themeName = data.Theme
	if m.themeName == "" {
		m.themeName = "green"
	}
	m.setWeekSpanFromData(data.WeekSpan)
	m.setWeekStartFromData(data.WeekStart)
	m.lofi.enabled = data.LofiEnabled
	m.lofi.url = strings.TrimSpace(data.LofiURL)
	m.lofi.status = lofiStatusStopped
	m.lofi.err = ""
	m.ensureTodoDueDates()
	m.sortExamsByPriority()
	m.sortProjectsByStatus()
	m.sortChecklistByDone()
	m.refreshAllFilters()

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
	m.refreshFlatExams()
	if m.lofiNow >= len(m.lofiPlaylist) {
		m.lofiNow = len(m.lofiPlaylist) - 1
	}
	if m.lofiNow < 0 && len(m.lofiPlaylist) > 0 {
		m.lofiNow = 0
	}
	if m.lofiCursor >= len(m.lofiPlaylist) {
		m.lofiCursor = len(m.lofiPlaylist) - 1
	}
	if m.lofiCursor < 0 && len(m.lofiPlaylist) > 0 {
		m.lofiCursor = 0
	}
	m.ensureLofiVisible()
}

func (m Model) exportData() storage.SemesterData {
	return storage.SemesterData{
		Subjects:    m.subjects,
		Projects:    m.projects,
		Checklist:   m.checklistItems,
		WeeklyExams: m.weeklyExams,
		ConfirmOn:   m.confirmOn,
		WeekStart:   m.weekStart.Format("2006-01-02"),
		WeekSpan:    m.weekSpan,
		LofiEnabled: m.lofi.enabled,
		LofiURL:     m.lofi.url,
		Theme:       m.themeName,
	}
}

func (m *Model) cycleTheme() {
	names := style.ThemeNames
	cur := 0
	for i, n := range names {
		if n == m.themeName {
			cur = i
			break
		}
	}
	m.themeName = names[(cur+1)%len(names)]
	m.persist()
}

func (m *Model) persist() {
	if m.store == nil {
		return
	}
	_ = m.store.Save(m.exportData())
}

func (m *Model) refreshChecklistView() {
	m.refreshTodoFilter()
	t := style.ThemeOf(m.themeName)
	filtered := make([]models.ChecklistItem, 0, len(m.todoVisible))
	for _, idx := range m.todoVisible {
		if idx >= 0 && idx < len(m.checklistItems) {
			filtered = append(filtered, m.checklistItems[idx])
		}
	}
	visibleCursor := m.visibleIndex(m.todoVisible, m.checklistCursor)
	m.checklist.SetContent(components.RenderChecklist(filtered, visibleCursor, m.activeTab == tabTodos, t))
	m.ensureChecklistVisible(visibleCursor, len(filtered))
}

func (m *Model) sortChecklistByDone() {
	if len(m.checklistItems) < 2 {
		return
	}

	oldCursor := m.checklistCursor
	if oldCursor < 0 {
		oldCursor = -1
	}
	indices := make([]int, 0, len(m.checklistItems))
	for i, item := range m.checklistItems {
		if !item.Done {
			indices = append(indices, i)
		}
	}
	for i, item := range m.checklistItems {
		if item.Done {
			indices = append(indices, i)
		}
	}

	newItems := make([]models.ChecklistItem, len(m.checklistItems))
	newCursor := -1
	for newIdx, oldIdx := range indices {
		newItems[newIdx] = m.checklistItems[oldIdx]
		if oldIdx == oldCursor {
			newCursor = newIdx
		}
	}

	m.checklistItems = newItems
	m.checklistCursor = newCursor
}

func (m *Model) ensureChecklistVisible(cursor, count int) {
	if count == 0 || m.checklist.Height <= 0 {
		return
	}
	if cursor < m.checklist.YOffset {
		m.checklist.SetYOffset(cursor)
		return
	}
	if cursor >= m.checklist.YOffset+m.checklist.Height {
		m.checklist.SetYOffset(cursor - m.checklist.Height + 1)
	}
}

func (m *Model) moveChecklistCursor(delta int) {
	if len(m.todoVisible) == 0 {
		return
	}
	pos := m.visibleIndex(m.todoVisible, m.checklistCursor)
	if pos < 0 {
		pos = 0
	}
	pos += delta
	if pos < 0 {
		pos = 0
	}
	if pos >= len(m.todoVisible) {
		pos = len(m.todoVisible) - 1
	}
	m.checklistCursor = m.todoVisible[pos]
	m.refreshChecklistView()
}

func (m *Model) toggleChecklistItem() {
	if len(m.todoVisible) == 0 {
		return
	}
	idx := m.checklistCursor
	if idx < 0 || idx >= len(m.checklistItems) {
		return
	}
	m.pushUndo()
	m.checklistItems[idx].Done = !m.checklistItems[idx].Done
	m.sortChecklistByDone()
	m.persist()
	m.refreshChecklistView()
}

func (m *Model) setWeekStartFromData(value string) {
	if value != "" {
		if t, err := time.ParseInLocation("2006-01-02", value, time.Local); err == nil {
			m.weekStart = weekStartOf(t)
			m.updateWeekLabel()
			return
		}
	}
	m.weekStart = weekStartOf(time.Now())
	m.updateWeekLabel()
}

func (m *Model) shiftWeek(delta int) {
	m.weekStart = m.weekStart.AddDate(0, 0, delta*7)
	m.updateWeekLabel()
	m.refreshFlatExams()
	m.refreshChecklistView()
	m.persist()
}

func (m *Model) jumpToCurrentWeek() {
	m.weekStart = weekStartOf(time.Now())
	m.updateWeekLabel()
	m.refreshFlatExams()
	m.refreshChecklistView()
	m.persist()
}

func (m *Model) updateWeekLabel() {
	m.weekLabel = weekLabel(m.weekStart, m.weekSpan)
}

func (m *Model) setWeekSpanFromData(value int) {
	switch value {
	case -1:
		m.weekSpan = -1
	case 1, 2, 3, 4:
		m.weekSpan = value
	default:
		m.weekSpan = 1
	}
}

func (m *Model) toggleGlobalView() {
	if m.weekSpan < 0 {
		if m.previousWeekSpan <= 0 {
			m.previousWeekSpan = 1
		}
		m.weekSpan = m.previousWeekSpan
	} else {
		m.previousWeekSpan = m.weekSpan
		m.weekSpan = -1
	}
	m.updateWeekLabel()
	m.refreshFlatExams()
	m.refreshChecklistView()
	m.persist()
}

func (m *Model) cycleWeekSpan() {
	switch m.weekSpan {
	case 1:
		m.weekSpan = 2
	case 2:
		m.weekSpan = 3
	case 3:
		m.weekSpan = 4
	case 4:
		m.weekSpan = -1
	default:
		m.weekSpan = 1
	}
	m.updateWeekLabel()
	m.refreshFlatExams()
	m.refreshChecklistView()
	m.persist()
}

func (m *Model) weekRange() (time.Time, time.Time, bool) {
	if m.weekSpan < 0 {
		return time.Time{}, time.Time{}, true
	}
	start := m.weekStart
	end := start.AddDate(0, 0, m.weekSpan*7)
	return start, end, false
}

func (m *Model) refreshTodoFilter() {
	start, end, all := m.weekRange()
	visible := make([]int, 0, len(m.checklistItems))
	for i, item := range m.checklistItems {
		if len(m.subjectFilters) > 0 && !m.isFiltered(item.Subject) {
			continue
		}
		if !all {
			due, ok := parseTodoDate(item.Due)
			if !ok {
				continue
			}
			if due.Before(start) || !due.Before(end) {
				continue
			}
		}
		visible = append(visible, i)
	}
	m.todoVisible = visible
	if len(m.todoVisible) == 0 {
		m.checklistCursor = -1
		return
	}
	if m.visibleIndex(m.todoVisible, m.checklistCursor) == -1 {
		m.checklistCursor = m.todoVisible[0]
	}
}

func (m *Model) isFiltered(subject string) bool {
	for _, f := range m.subjectFilters {
		if strings.EqualFold(f, subject) {
			return true
		}
	}
	return false
}

func (m *Model) subjectFilterLabel() string {
	if len(m.subjectFilters) == 0 {
		return ""
	}
	return strings.Join(m.subjectFilters, ", ")
}

func (m *Model) refreshProjectFilter() {
	visible := make([]int, 0, len(m.projects))
	for i, p := range m.projects {
		if len(m.subjectFilters) > 0 && !m.isFiltered(p.Subject) {
			continue
		}
		visible = append(visible, i)
	}
	m.projectVisible = visible
	if len(m.projectVisible) == 0 {
		m.projectCursor = -1
		return
	}
	if m.visibleIndex(m.projectVisible, m.projectCursor) == -1 {
		m.projectCursor = m.projectVisible[0]
	}
}

func (m *Model) refreshAllFilters() {
	m.refreshFlatExams()
	m.refreshProjectFilter()
	m.refreshChecklistView() // internally calls refreshTodoFilter
}

func (m *Model) refreshFlatExams() {
	start, end, all := m.weekRange()
	var list []models.FlatExam
	for si, subject := range m.subjects {
		if len(m.subjectFilters) > 0 && !m.isFiltered(subject.Code) {
			continue
		}
		for ei, exam := range subject.Exams {
			if !all {
				date, ok := parseExamDate(exam.Date)
				if !ok {
					continue
				}
				if date.Before(start) || !date.Before(end) {
					continue
				}
			}
			list = append(list, models.FlatExam{
				SubjectCode: subject.Code,
				SubjectIdx:  si,
				ExamIdx:     ei,
				Exam:        exam,
			})
		}
	}
	sort.Slice(list, func(i, j int) bool {
		di, oki := parseExamDate(list[i].Exam.Date)
		dj, okj := parseExamDate(list[j].Exam.Date)
		if !oki {
			return false
		}
		if !okj {
			return true
		}
		return di.Before(dj)
	})
	m.flatExams = list
	if len(m.flatExams) == 0 {
		m.examCursor = -1
	} else if m.examCursor < 0 || m.examCursor >= len(m.flatExams) {
		m.examCursor = 0
	}
}

func (m *Model) visibleIndex(items []int, value int) int {
	for i, idx := range items {
		if idx == value {
			return i
		}
	}
	return -1
}

func (m *Model) moveExamCursor(delta int) {
	if len(m.flatExams) == 0 {
		return
	}
	m.examCursor += delta
	if m.examCursor < 0 {
		m.examCursor = 0
	}
	if m.examCursor >= len(m.flatExams) {
		m.examCursor = len(m.flatExams) - 1
	}
}

func (m *Model) moveProjectCursor(delta int) {
	if len(m.projectVisible) == 0 {
		return
	}
	pos := m.visibleIndex(m.projectVisible, m.projectCursor)
	if pos < 0 {
		pos = 0
	} else {
		pos += delta
	}
	if pos < 0 {
		pos = 0
	}
	if pos >= len(m.projectVisible) {
		pos = len(m.projectVisible) - 1
	}
	m.projectCursor = m.projectVisible[pos]
}

func (m *Model) ensureTodoDueDates() {
	if len(m.checklistItems) == 0 {
		return
	}
	changed := false
	defaultDue := m.weekStart.Format("2006-01-02")
	for i := range m.checklistItems {
		if strings.TrimSpace(m.checklistItems[i].Due) == "" {
			m.checklistItems[i].Due = defaultDue
			changed = true
		}
	}
	if changed {
		m.persist()
	}
}

func (m Model) viewState() screens.State {
	start, end, all := m.weekRange()

	visibleProjects := make([]models.ProjectItem, 0, len(m.projectVisible))
	for _, idx := range m.projectVisible {
		if idx >= 0 && idx < len(m.projects) {
			visibleProjects = append(visibleProjects, m.projects[idx])
		}
	}
	visibleProjectCursor := m.visibleIndex(m.projectVisible, m.projectCursor)

	state := screens.State{
		ActiveTab:     m.activeTab,
		ConfirmOn:     m.confirmOn,
		ThemeName:     m.themeName,
		SubjectFilter: m.subjectFilterLabel(),
		ChecklistView: m.checklist.View(),
		Projects:      visibleProjects,
		Subjects:      m.subjects,
		SelectedSubj:  m.selectedSubj,
		ExamCursor:    m.examCursor,
		FlatExams:     m.flatExams,
		WeekLabel:     m.weekLabel,
		FilterStart:   start,
		FilterEnd:     end,
		FilterAll:     all,
		WeekSpan:      m.weekSpan,
		WeeklyExams:   m.weeklyExams,
		ProjectCursor: visibleProjectCursor,
		LofiEnabled:   m.lofi.enabled,
		LofiURL:       m.lofi.url,
		LofiStatus:    m.lofi.status,
		LofiError:     m.lofi.err,
		LofiPlaylist:  m.lofiPlaylist,
		LofiCursor:    m.lofiCursor,
		LofiOffset:    m.lofiOffset,
		LofiNow:       m.lofiNow,
	}

	if m.modal == modalNone {
		return state
	}

	fields := make([]components.ModalField, 0, len(m.formFields))
	for _, field := range m.formFields {
		fields = append(fields, components.ModalField{
			Label: field.label,
			Value: field.input.View(),
		})
	}

	dropdownFieldIdx := -1
	if m.isSubjectField() && len(m.dropdownMatches) > 0 {
		dropdownFieldIdx = m.formFocus
	}
	modalState := components.ModalState{
		Title:            m.modalTitle,
		Hint:             m.modalHint,
		Fields:           fields,
		DropdownItems:    m.dropdownMatches,
		DropdownCursor:   m.dropdownCursor,
		DropdownFieldIdx: dropdownFieldIdx,
	}
	switch m.modal {
	case modalConfirm:
		modalState.Mode = components.ModalConfirm
		modalState.Message = m.modalError
	case modalSubjectFilter:
		modalState.Mode = components.ModalSubjectSelect
		items := make([]string, len(m.subjects))
		for i, s := range m.subjects {
			items[i] = s.Code
		}
		modalState.SelectItems = items
		modalState.SelectActive = m.filterModalActive
		modalState.SelectCursor = m.filterModalCursor
	default:
		modalState.Mode = components.ModalForm
		modalState.Error = m.modalError
	}
	state.Modal = modalState
	return state
}

func (m *Model) examsForSelected() []models.ExamItem {
	if m.selectedSubj < 0 || m.selectedSubj >= len(m.subjects) {
		return nil
	}
	return m.subjects[m.selectedSubj].Exams
}

func (m *Model) sortProjectsByStatus() {
	if len(m.projects) < 2 {
		return
	}
	selectedKey := projectKey(m.projects, m.projectCursor)
	sort.SliceStable(m.projects, func(i, j int) bool {
		return statusRank(m.projects[i].Status) < statusRank(m.projects[j].Status)
	})
	if selectedKey != "" {
		if idx := findProjectIndex(m.projects, selectedKey); idx >= 0 {
			m.projectCursor = idx
		}
	}
}

func (m *Model) sortExamsByPriority() {
	if len(m.subjects) == 0 {
		return
	}
	var selectedKey string
	selectedSubject := m.selectedSubj
	if selectedSubject >= 0 && selectedSubject < len(m.subjects) {
		selectedKey = examKey(m.subjects[selectedSubject].Exams, m.examCursor)
	}
	for i := range m.subjects {
		exams := m.subjects[i].Exams
		if len(exams) < 2 {
			continue
		}
		sort.SliceStable(exams, func(a, b int) bool {
			return priorityRank(exams[a].Priority) < priorityRank(exams[b].Priority)
		})
		m.subjects[i].Exams = exams
	}
	if selectedKey != "" && selectedSubject >= 0 && selectedSubject < len(m.subjects) {
		if idx := findExamIndex(m.subjects[selectedSubject].Exams, selectedKey); idx >= 0 {
			m.examCursor = idx
		}
	}
	m.refreshFlatExams()
}

func projectKey(items []models.ProjectItem, idx int) string {
	if idx < 0 || idx >= len(items) {
		return ""
	}
	item := items[idx]
	return strings.ToUpper(item.Name) + "|" + strings.ToUpper(item.Subject) + "|" + item.Due
}

func examKey(items []models.ExamItem, idx int) string {
	if idx < 0 || idx >= len(items) {
		return ""
	}
	item := items[idx]
	return strings.ToUpper(item.Name) + "|" + item.Date + "|" + strings.ToUpper(item.Priority)
}

func findProjectIndex(items []models.ProjectItem, key string) int {
	for i := range items {
		if projectKey(items, i) == key {
			return i
		}
	}
	return -1
}

func findExamIndex(items []models.ExamItem, key string) int {
	for i := range items {
		if examKey(items, i) == key {
			return i
		}
	}
	return -1
}

func statusRank(status string) int {
	switch strings.ToUpper(strings.TrimSpace(status)) {
	case "NOT STARTED":
		return 0
	case "IN PROGRESS":
		return 1
	case "DONE":
		return 2
	default:
		return 3
	}
}

func priorityRank(priority string) int {
	switch strings.ToUpper(strings.TrimSpace(priority)) {
	case "HIGH":
		return 0
	case "MED":
		return 1
	case "LOW":
		return 2
	default:
		return 3
	}
}

func (m Model) tabItems() []components.TabItem {
	items := []components.TabItem{
		{ID: tabDashboard, Label: "Dashboard"},
		{ID: tabSubjects, Label: "Subjects"},
		{ID: tabExams, Label: "Exams"},
		{ID: tabTodos, Label: "Todos"},
		{ID: tabProjects, Label: "Projects"},
		{ID: tabSettings, Label: "Settings"},
	}
	if m.lofi.enabled {
		items = append(items, components.TabItem{ID: tabLofi, Label: "Lofi"})
	}
	return items
}

func (m *Model) switchToTab(index int) bool {
	items := m.tabItems()
	if index < 0 || index >= len(items) {
		return false
	}
	m.activeTab = items[index].ID
	m.resize(m.width, m.height)
	return true
}

func (m *Model) consumeLofiReload() tea.Cmd {
	if !m.lofiReload {
		return nil
	}
	m.lofiReload = false
	if strings.TrimSpace(m.lofi.url) == "" {
		return nil
	}
	return loadLofiPlaylist(m.lofi.url)
}

func (m *Model) resizeLofi(width, height int) {
	gap := 1
	available := width - gap
	if available < 1 {
		available = width
		gap = 0
	}
	leftWidth := available / 2
	rightWidth := available - leftWidth
	if leftWidth < 1 {
		leftWidth = available
		rightWidth = 0
	}
	playlistWidth := rightWidth
	if rightWidth == 0 {
		m.lofiListHeight = 0
		m.lofiOffset = 0
		return
	}

	_, contentH := components.PanelContentSize(playlistWidth, height)
	if contentH < 1 {
		contentH = 1
	}
	m.lofiListHeight = contentH
	m.ensureLofiVisible()
}
