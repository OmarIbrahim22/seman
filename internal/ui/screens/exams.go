package screens

import (
	"github.com/OmarIbrahim22/seman/internal/style"
	"github.com/OmarIbrahim22/seman/internal/ui/components"
)

func RenderExams(state State, width, height int, t style.Theme) string {
	var header string
	if state.SubjectFilter != "" {
		header = t.Title.Render("Subject: "+state.SubjectFilter) + "  " + t.Dim.Render("[R] clear filter")
		header += "\n\n"
	}

	body := header + components.RenderFlatExams(state.FlatExams, state.ExamCursor, state.SubjectFilter, t)
	return components.RenderPanel(width, height, "Exams", body, t)
}
