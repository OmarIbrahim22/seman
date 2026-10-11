package screens

import (
	"time"

	"github.com/OmarIbrahim22/seman/internal/models"
	"github.com/OmarIbrahim22/seman/internal/ui/components"
)

type State struct {
	ActiveTab     int
	ConfirmOn     bool
	ThemeName     string
	SubjectFilter string
	ChecklistView string
	Projects      []models.ProjectItem
	Subjects      []models.SubjectItem
	SelectedSubj  int
	ExamCursor    int
	FlatExams     []models.FlatExam
	WeekLabel     string
	FilterStart   time.Time
	FilterEnd     time.Time
	FilterAll     bool
	WeekSpan      int
	WeeklyExams   []string
	ProjectCursor int
	LofiEnabled   bool
	LofiURL       string
	LofiStatus    string
	LofiError     string
	LofiPlaylist  []models.LofiTrack
	LofiCursor    int
	LofiOffset    int
	LofiNow       int
	Modal         components.ModalState
}
