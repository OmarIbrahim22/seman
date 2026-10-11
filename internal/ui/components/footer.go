package components

import (
	"fmt"

	"github.com/OmarIbrahim22/seman/internal/style"
	"github.com/charmbracelet/lipgloss"
)

// Tab ID constants kept in sync with app/model.go.
const (
	footerTabDashboard = 0
	footerTabExams     = 1
	footerTabTodos     = 2
	footerTabProjects  = 3
	footerTabSettings  = 4
	footerTabLofi      = 5
	footerTabSubjects  = 6
)

func tabHint(activeTab int) string {
	switch activeTab {
	case footerTabSubjects:
		return "[S] Add subject  [E] Edit  [D] Delete  [T] Today  [Q] Quit"
	case footerTabExams:
		return "[A] Add exam  [E] Edit  [D] Delete  [F] Filter  [R] Clear  [G] Global  [T] Today  [Q] Quit"
	case footerTabTodos:
		return "[N] Add  [B] Bulk add  [E] Edit  [D] Delete  [Space] Toggle  [F] Filter  [R] Clear  [G] Global  [T] Today  [Q] Quit"
	case footerTabProjects:
		return "[P] Add project  [E] Edit  [D] Delete  [F] Filter  [R] Clear  [T] Today  [Q] Quit"
	case footerTabSettings:
		return "[T] Theme  [O] Confirm  [W] Week span  [L] Lofi  [U] Lofi URL  [Q] Quit"
	case footerTabLofi:
		return "[Enter] Play  [Space] Pause  [N] Next  [B] Prev  [X] Stop  [Q] Quit"
	default: // Dashboard
		return "[S] Subject  [A] Exam  [P] Project  [F] Filter  [R] Clear  [G] Global  [T] Today  [Q] Quit"
	}
}

func RenderFooter(width, tabCount, activeTab int, t style.Theme) string {
	contentWidth := width - barBorderX - barPaddingX*2
	if contentWidth < 1 {
		contentWidth = 1
	}

	left := t.FooterHint.Render(tabHint(activeTab))
	right := t.FooterHint.Render(fmt.Sprintf("[1-%d] Switch tabs", tabCount))
	content := AlignLine(contentWidth, left, right)

	styleWidth := width - barBorderX
	if styleWidth < 1 {
		styleWidth = 1
	}

	box := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(t.Border).
		Padding(barPaddingY, barPaddingX).
		Width(styleWidth)

	return box.Render(content)
}
