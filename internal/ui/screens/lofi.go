package screens

import (
	"fmt"
	"strings"

	"github.com/OmarIbrahim22/seman/internal/models"
	"github.com/OmarIbrahim22/seman/internal/style"
	"github.com/OmarIbrahim22/seman/internal/ui/components"
	"github.com/charmbracelet/lipgloss"
)

func RenderLofi(state State, width, height int, t style.Theme) string {
	if !state.LofiEnabled {
		body := t.Dim.Render("Enable the Lofi tab in Settings to use the player.")
		return components.RenderPanel(width, height, "Lofi Player", body, t)
	}

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

	playerWidth := leftWidth
	playlistWidth := rightWidth
	if rightWidth == 0 {
		playerWidth = width
	}

	title, note := currentLofiInfo(state)
	if title == "" {
		title = "No track selected"
	}
	if note == "" {
		note = "Press Enter to play the selected track"
	}

	status := state.LofiStatus
	if status == "" {
		status = "Stopped"
	}

	topLines := []string{
		"Now Playing",
		title,
		note,
		"Status: " + status,
		"",
		buildLofiControls(state.LofiStatus, components.PanelContentWidth(playerWidth), t),
	}
	if state.LofiError != "" {
		topLines = append(topLines, "", state.LofiError)
	}

	topContentW := components.PanelContentWidth(playerWidth)
	centered := make([]string, 0, len(topLines))
	for i, line := range topLines {
		if i == 5 {
			centered = append(centered, line)
			continue
		}
		style := lipgloss.NewStyle().Width(topContentW).Align(lipgloss.Center)
		if i == 0 {
			centered = append(centered, t.Dim.Render(style.Render(line)))
		} else if i == 2 || i == 3 {
			centered = append(centered, t.Dim.Render(style.Render(line)))
		} else if i == len(topLines)-1 && state.LofiError != "" {
			centered = append(centered, t.ModalError.Render(style.Render(line)))
		} else {
			centered = append(centered, t.Text.Render(style.Render(line)))
		}
	}
	topBody := strings.Join(centered, "\n")
	playerPanel := components.RenderPanel(playerWidth, height, "Lofi Music Player", topBody, t)
	if rightWidth == 0 {
		return playerPanel
	}

	_, playlistContentH := components.PanelContentSize(playlistWidth, height)
	if playlistContentH < 1 {
		playlistContentH = 1
	}
	maxLines := models.LofiVisibleCap*2 - 1
	if playlistContentH > maxLines {
		playlistContentH = maxLines
	}
	playlistBody := renderLofiPlaylist(state.LofiPlaylist, state.LofiCursor, state.LofiOffset, playlistContentH, playlistWidth, t)
	playlistPanel := components.RenderPanel(playlistWidth, height, "Playlist", playlistBody, t)

	if gap > 0 {
		return lipgloss.JoinHorizontal(lipgloss.Top, playerPanel, " ", playlistPanel)
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, playerPanel, playlistPanel)
}

func buildLofiControls(status string, width int, t style.Theme) string {
	label := "Play"
	if strings.EqualFold(status, "Playing") {
		label = "Pause"
	}
	button := func(text string) string {
		style := lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(t.Border).
			Padding(0, 1).
			Foreground(t.Accent)
		return style.Render(text)
	}
	controls := []string{
		button("Previous [B]"),
		button(label + " [Space]"),
		button("Next [N]"),
	}
	line := lipgloss.JoinHorizontal(lipgloss.Top, controls...)
	return lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(line)
}

func currentLofiInfo(state State) (string, string) {
	idx := state.LofiNow
	if idx < 0 || idx >= len(state.LofiPlaylist) {
		idx = state.LofiCursor
	}
	if idx < 0 || idx >= len(state.LofiPlaylist) {
		return "", ""
	}
	item := state.LofiPlaylist[idx]
	return item.Title, item.Note
}

func renderLofiPlaylist(items []models.LofiTrack, cursor, offset, visibleLines, width int, t style.Theme) string {
	if len(items) == 0 {
		return t.Dim.Render("No playlist items yet.")
	}

	contentW := components.PanelContentWidth(width)
	lines := make([]string, 0, len(items)*2-1)
	for i, item := range items {
		left := t.Text.Render(fmt.Sprintf("[%d] %s", i+1, item.Title))
		right := t.Dim.Render(item.Note)
		line := components.AlignLine(contentW, left, right)
		if i == cursor {
			line = t.RowActive.Render(line)
		}
		lines = append(lines, line)
		if i < len(items)-1 {
			lines = append(lines, t.Dim.Render(strings.Repeat("-", contentW)))
		}
	}

	if visibleLines <= 0 || visibleLines >= len(lines) {
		return strings.Join(lines, "\n")
	}
	if offset < 0 {
		offset = 0
	}
	maxOffset := len(lines) - visibleLines
	if offset > maxOffset {
		offset = maxOffset
	}
	return strings.Join(lines[offset:offset+visibleLines], "\n")
}
