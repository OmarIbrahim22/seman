package components

import (
	"fmt"
	"strings"

	"github.com/OmarIbrahim22/seman/internal/style"
	"github.com/charmbracelet/lipgloss"
)

type ModalMode int

const (
	ModalHidden ModalMode = iota
	ModalForm
	ModalConfirm
	ModalSubjectSelect
)

type ModalField struct {
	Label string
	Value string
}

type ModalState struct {
	Mode             ModalMode
	Title            string
	Hint             string
	Message          string
	Error            string
	Fields           []ModalField
	DropdownItems    []string
	DropdownCursor   int
	DropdownFieldIdx int // index of the field that owns the dropdown; -1 if none
	SelectItems      []string
	SelectActive     []bool
	SelectCursor     int
}

const dropdownPanelWidth = 18
const dropdownMaxVisible = 8

// RenderModalContent returns just the modal box (+ dropdown panel if active),
// without padding it to fill the full area. Use with PlaceOverlay.
func RenderModalContent(state ModalState, width int, t style.Theme) string {
	modal := RenderModal(state, width, t)
	if state.DropdownFieldIdx >= 0 && len(state.DropdownItems) > 0 {
		dropdown := renderDropdownPanel(state, t)
		return lipgloss.JoinHorizontal(lipgloss.Top, modal, " ", dropdown)
	}
	return modal
}

func RenderModalArea(state ModalState, width, height int, t style.Theme) string {
	modal := RenderModal(state, width, t)

	if state.DropdownFieldIdx >= 0 && len(state.DropdownItems) > 0 {
		dropdown := renderDropdownPanel(state, t)
		combined := lipgloss.JoinHorizontal(lipgloss.Top, modal, " ", dropdown)

		// Left-pad so the modal itself (not the combined block) is centered.
		modalW := lipgloss.Width(modal)
		leftPad := (width - modalW) / 2
		if leftPad < 0 {
			leftPad = 0
		}
		pad := strings.Repeat(" ", leftPad)

		// Vertically center the combined block.
		combinedH := lipgloss.Height(combined)
		topLines := (height - combinedH) / 2
		if topLines < 0 {
			topLines = 0
		}

		var out strings.Builder
		for i := 0; i < topLines; i++ {
			out.WriteString("\n")
		}
		lines := strings.Split(combined, "\n")
		for i, line := range lines {
			if i > 0 {
				out.WriteString("\n")
			}
			out.WriteString(pad + line)
		}
		return out.String()
	}

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, modal)
}

func RenderModal(state ModalState, width int, t style.Theme) string {
	modalW := minInt(70, width-6)
	if modalW < 42 {
		modalW = 42
	}
	box := t.ModalBorder.Copy().Padding(1, 2).Width(modalW)

	var b strings.Builder
	title := t.ModalTitle.Render(state.Title)
	if title != "" {
		b.WriteString(title)
		b.WriteString("\n")
	}
	if state.Mode == ModalConfirm {
		b.WriteString("\n")
		b.WriteString(t.Text.Render(state.Message))
		b.WriteString("\n\n")
		b.WriteString(t.ModalHint.Render(state.Hint))
		return box.Render(b.String())
	}

	if state.Mode == ModalSubjectSelect {
		b.WriteString("\n")
		if len(state.SelectItems) == 0 {
			b.WriteString(t.Dim.Render("No subjects available."))
		} else {
			for i, item := range state.SelectItems {
				if i > 0 {
					b.WriteString("\n")
				}
				checked := "[ ]"
				if i < len(state.SelectActive) && state.SelectActive[i] {
					checked = "[x]"
				}
				line := checked + " " + item
				if i == state.SelectCursor {
					b.WriteString(t.RowActive.Render(line))
				} else {
					b.WriteString(t.Text.Render(line))
				}
			}
		}
		if state.Hint != "" {
			b.WriteString("\n\n")
			b.WriteString(t.ModalHint.Render(state.Hint))
		}
		return box.Render(b.String())
	}

	b.WriteString("\n")
	for i, field := range state.Fields {
		label := fmt.Sprintf("%-12s", field.Label+":")
		b.WriteString(t.Dim.Render(label))
		b.WriteString(" ")
		b.WriteString(field.Value)
		if i < len(state.Fields)-1 {
			b.WriteString("\n")
		}
	}
	if state.Error != "" {
		b.WriteString("\n\n")
		b.WriteString(t.ModalError.Render(state.Error))
	}
	if state.Hint != "" {
		b.WriteString("\n\n")
		b.WriteString(t.ModalHint.Render(state.Hint))
	}

	return box.Render(b.String())
}

func renderDropdownPanel(state ModalState, t style.Theme) string {
	items := state.DropdownItems
	cursor := state.DropdownCursor

	// Sliding window centred on the cursor.
	start := cursor - dropdownMaxVisible/2
	if start < 0 {
		start = 0
	}
	end := start + dropdownMaxVisible
	if end > len(items) {
		end = len(items)
		start = end - dropdownMaxVisible
		if start < 0 {
			start = 0
		}
	}

	var b strings.Builder
	b.WriteString(t.Dim.Render("Subjects"))
	for i := start; i < end; i++ {
		b.WriteString("\n")
		if i == cursor {
			b.WriteString(t.SubjectActive.Render("> " + items[i]))
		} else {
			b.WriteString(t.Dim.Render("  " + items[i]))
		}
	}
	if len(items) > dropdownMaxVisible {
		b.WriteString("\n")
		remaining := len(items) - end
		if remaining > 0 {
			b.WriteString(t.Dim.Render(fmt.Sprintf("  +%d more", remaining)))
		}
	}

	box := t.ModalBorder.Copy().Padding(1, 2).Width(dropdownPanelWidth)
	return box.Render(b.String())
}
