package screens

import (
	"github.com/OmarIbrahim22/seman/internal/style"
	"github.com/OmarIbrahim22/seman/internal/ui/components"
)

func RenderSubjectsTab(state State, width, height int, t style.Theme) string {
	var body string
	if len(state.Subjects) == 0 {
		body = t.Dim.Render("No subjects yet — press [A] to add one")
	} else {
		body = components.RenderSubjects(state.Subjects, state.SelectedSubj, t)
	}
	return components.RenderPanel(width, height, "Subjects", body, t)
}
