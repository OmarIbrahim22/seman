package storage

import "github.com/OmarIbrahim22/seman/internal/models"

type SemesterData struct {
	Subjects    []models.SubjectItem   `json:"subjects"`
	Projects    []models.ProjectItem   `json:"projects"`
	Checklist   []models.ChecklistItem `json:"checklist"`
	WeeklyExams []string               `json:"weekly_exams"`
	ConfirmOn   bool                   `json:"confirm_on"`
	WeekStart   string                 `json:"week_start"`
	WeekSpan    int                    `json:"week_span"`
	LofiEnabled bool                   `json:"lofi_enabled"`
	LofiURL     string                 `json:"lofi_url"`
	Theme       string                 `json:"theme"`
}

type Store interface {
	Load() (SemesterData, bool, error)
	Save(SemesterData) error
}
