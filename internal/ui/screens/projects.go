package screens

import (
	"github.com/OmarIbrahim22/seman/internal/style"
	"github.com/OmarIbrahim22/seman/internal/ui/components"
)

func RenderProjectsTab(state State, width, height int, t style.Theme) string {
	layout := components.ComputeProjectsLayout(width, height)

	title := "Projects Overview"
	var header string
	if state.SubjectFilter != "" {
		title = "Projects · " + state.SubjectFilter
		header = t.Title.Render("Subject: "+state.SubjectFilter) + "  " + t.Dim.Render("[R] clear filter") + "\n\n"
	}

	body := header + components.RenderProjectsTable(state.Projects, state.ProjectCursor, layout.TableWidth, t)
	return components.RenderPanel(width, height, title, body, t)
}
