package components

import (
	"strings"

	"github.com/OmarIbrahim22/seman/internal/style"
)

func RenderDivider(width int, t style.Theme) string {
	if width < 1 {
		return ""
	}
	line := strings.Repeat("-", width)
	return t.Dim.Render(line)
}
