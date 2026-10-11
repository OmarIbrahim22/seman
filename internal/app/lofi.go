package app

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/OmarIbrahim22/seman/internal/models"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	lofiStatusStopped = "Stopped"
	lofiStatusPlaying = "Playing"
	lofiStatusPaused  = "Paused"
	lofiStatusLoading = "Loading..."
)

type lofiState struct {
	enabled    bool
	url        string
	status     string
	err        string
	socketPath string
	cmd        *exec.Cmd
}

type lofiExitMsg struct {
	err error
}

type lofiPlaylistMsg struct {
	tracks []models.LofiTrack
	err    error
}

type lofiIndexMsg struct {
	index int
}

type lofiPlaybackMsg struct {
	playing  bool
	attempts int
	err      error
}

func defaultLofiPlaylist() []models.LofiTrack {
	return []models.LofiTrack{
		{Title: "Lofi Hip Hop Radio - Beats to Study/Relax", Note: "Classic lofi beats"},
		{Title: "ChillHop Radio - Jazzy & Lofi Hip Hop", Note: "Smooth jazzy lofi"},
		{Title: "Lofi Girl - Sleep/Chill Radio", Note: "Calm beats for sleeping"},
		{Title: "Coffee Shop Ambience", Note: "Cozy cafe vibes"},
		{Title: "Rainy Day Lofi", Note: "Rain sounds with beats"},
		{Title: "Midnight Lofi Mix", Note: "Late night study session"},
		{Title: "Synthwave Lofi", Note: "Retro synth lofi beats"},
		{Title: "Japanese Lofi", Note: "Traditional Japanese inspired"},
		{Title: "Piano Lofi Collection", Note: "Soft piano melodies"},
		{Title: "Chill Beats Vol. 1", Note: "Best chill collection"},
	}
}

func (m *Model) toggleLofiEnabled() tea.Cmd {
	m.lofi.enabled = !m.lofi.enabled
	if !m.lofi.enabled {
		m.shutdownLofi()
		m.lofi.err = ""
		if m.activeTab == tabLofi {
			m.activeTab = tabSettings
		}
	}
	m.persist()
	if m.lofi.enabled && strings.TrimSpace(m.lofi.url) != "" {
		return loadLofiPlaylist(m.lofi.url)
	}
	return nil
}

func (m *Model) moveLofiCursor(delta int) {
	if len(m.lofiPlaylist) == 0 {
		m.lofiCursor = 0
		m.lofiOffset = 0
		return
	}
	m.lofiCursor += delta
	if m.lofiCursor < 0 {
		m.lofiCursor = 0
	}
	if m.lofiCursor >= len(m.lofiPlaylist) {
		m.lofiCursor = len(m.lofiPlaylist) - 1
	}
	m.ensureLofiVisible()
}

func (m *Model) applyLofiPlaylist(msg lofiPlaylistMsg) {
	if msg.err != nil {
		m.lofi.err = msg.err.Error()
		return
	}
	if len(msg.tracks) == 0 {
		m.lofi.err = "No tracks found."
		return
	}
	m.lofiPlaylist = msg.tracks
	m.lofi.err = ""
	if m.lofiCursor >= len(m.lofiPlaylist) || m.lofiCursor < 0 {
		m.lofiCursor = 0
	}
	if m.lofiNow >= len(m.lofiPlaylist) || m.lofiNow < 0 {
		m.lofiNow = 0
	}
	m.ensureLofiVisible()
}

func (m *Model) applyLofiIndex(msg lofiIndexMsg) tea.Cmd {
	if msg.index < 0 {
		return nil
	}
	if err := m.sendLofiCommand("playlist-play-index", msg.index); err != nil {
		m.lofi.err = err.Error()
		return nil
	}
	m.lofi.status = lofiStatusLoading
	m.lofi.err = ""
	return m.pollLofiPlayback(0)
}

func (m *Model) ensureLofiVisible() {
	if len(m.lofiPlaylist) == 0 {
		m.lofiOffset = 0
		return
	}
	visible := m.lofiListHeight
	if visible < 1 {
		m.lofiOffset = 0
		return
	}
	maxLines := models.LofiVisibleCap*2 - 1
	if visible > maxLines {
		visible = maxLines
	}
	totalLines := len(m.lofiPlaylist)*2 - 1
	if totalLines <= visible {
		m.lofiOffset = 0
		return
	}
	lineIdx := m.lofiCursor * 2
	if lineIdx < m.lofiOffset {
		m.lofiOffset = lineIdx
	}
	if lineIdx >= m.lofiOffset+visible {
		m.lofiOffset = lineIdx - visible + 1
	}
	maxOffset := totalLines - visible
	if m.lofiOffset > maxOffset {
		m.lofiOffset = maxOffset
	}
	if m.lofiOffset < 0 {
		m.lofiOffset = 0
	}
}

func (m *Model) handleLofiExit(msg lofiExitMsg) {
	m.lofi.cmd = nil
	m.lofi.socketPath = ""
	m.lofi.status = lofiStatusStopped
	if msg.err != nil {
		m.lofi.err = "Lofi player stopped."
	} else {
		m.lofi.err = ""
	}
}

func (m *Model) applyLofiPlayback(msg lofiPlaybackMsg) tea.Cmd {
	if msg.err != nil {
		m.lofi.err = msg.err.Error()
		return nil
	}
	if msg.playing {
		m.lofi.status = lofiStatusPlaying
		m.lofi.err = ""
		return nil
	}
	if m.lofi.status != lofiStatusLoading {
		return nil
	}
	if msg.attempts >= 30 {
		return nil
	}
	return m.pollLofiPlayback(msg.attempts + 1)
}

func (m *Model) playLofiAt(index int) tea.Cmd {
	if err := m.validateLofi(); err != nil {
		m.lofi.err = err.Error()
		return nil
	}
	if index < 0 || index >= len(m.lofiPlaylist) {
		m.lofi.err = "Select a track first."
		return nil
	}
	m.lofiNow = index
	m.lofiCursor = index
	m.ensureLofiVisible()

	var cmd tea.Cmd
	if m.lofi.cmd == nil {
		cmd = m.startLofi()
	} else if m.lofi.status == lofiStatusStopped {
		_ = m.sendLofiCommand("loadfile", m.lofi.url, "replace")
	}
	m.lofi.status = lofiStatusLoading
	m.lofi.err = ""
	return tea.Batch(cmd, deferLofiIndex(index), m.pollLofiPlayback(0))
}

func deferLofiIndex(index int) tea.Cmd {
	if index < 0 {
		return nil
	}
	return tea.Tick(200*time.Millisecond, func(time.Time) tea.Msg {
		return lofiIndexMsg{index: index}
	})
}

func (m *Model) pollLofiPlayback(attempts int) tea.Cmd {
	return tea.Tick(400*time.Millisecond, func(time.Time) tea.Msg {
		playing, err := m.isLofiPlaying()
		return lofiPlaybackMsg{playing: playing, attempts: attempts, err: err}
	})
}

func (m *Model) toggleLofiPlayPause() tea.Cmd {
	if err := m.validateLofi(); err != nil {
		m.lofi.err = err.Error()
		return nil
	}
	if m.lofiNow < 0 && len(m.lofiPlaylist) > 0 {
		if m.lofiCursor >= 0 && m.lofiCursor < len(m.lofiPlaylist) {
			m.lofiNow = m.lofiCursor
		} else {
			m.lofiNow = 0
		}
	}
	if m.lofi.cmd == nil {
		m.lofi.status = lofiStatusLoading
		cmd := m.startLofi()
		if m.lofiNow >= 0 {
			return tea.Batch(cmd, deferLofiIndex(m.lofiNow), m.pollLofiPlayback(0))
		}
		return tea.Batch(cmd, m.pollLofiPlayback(0))
	}
	if m.lofi.status == lofiStatusStopped {
		if err := m.sendLofiCommand("loadfile", m.lofi.url, "replace"); err != nil {
			m.lofi.err = err.Error()
			return nil
		}
		m.lofi.status = lofiStatusLoading
		m.lofi.err = ""
		return m.pollLofiPlayback(0)
	}
	if err := m.sendLofiCommand("cycle", "pause"); err != nil {
		m.lofi.err = err.Error()
		return nil
	}
	if m.lofi.status == lofiStatusPaused {
		m.lofi.status = lofiStatusPlaying
	} else {
		m.lofi.status = lofiStatusPaused
	}
	m.lofi.err = ""
	return nil
}

func (m *Model) lofiNext() tea.Cmd {
	if err := m.validateLofi(); err != nil {
		m.lofi.err = err.Error()
		return nil
	}
	if m.lofi.cmd == nil {
		m.lofi.status = lofiStatusLoading
		return tea.Batch(m.startLofi(), m.pollLofiPlayback(0))
	}
	if err := m.sendLofiCommand("playlist-next", "force"); err != nil {
		m.lofi.err = err.Error()
		return nil
	}
	if len(m.lofiPlaylist) > 0 && m.lofiNow < len(m.lofiPlaylist)-1 {
		m.lofiNow++
		m.lofiCursor = m.lofiNow
		m.ensureLofiVisible()
	}
	m.lofi.status = lofiStatusLoading
	m.lofi.err = ""
	return m.pollLofiPlayback(0)
}

func (m *Model) lofiPrev() tea.Cmd {
	if err := m.validateLofi(); err != nil {
		m.lofi.err = err.Error()
		return nil
	}
	if m.lofi.cmd == nil {
		m.lofi.status = lofiStatusLoading
		return tea.Batch(m.startLofi(), m.pollLofiPlayback(0))
	}
	if err := m.sendLofiCommand("playlist-prev", "force"); err != nil {
		m.lofi.err = err.Error()
		return nil
	}
	if len(m.lofiPlaylist) > 0 && m.lofiNow > 0 {
		m.lofiNow--
		m.lofiCursor = m.lofiNow
		m.ensureLofiVisible()
	}
	m.lofi.status = lofiStatusLoading
	m.lofi.err = ""
	return m.pollLofiPlayback(0)
}

func (m *Model) lofiStop() {
	if m.lofi.cmd == nil {
		m.lofi.status = lofiStatusStopped
		return
	}
	if err := m.sendLofiCommand("stop"); err != nil {
		m.lofi.err = err.Error()
		return
	}
	m.lofi.status = lofiStatusStopped
	m.lofi.err = ""
}

func (m *Model) startLofi() tea.Cmd {
	if err := m.validateLofi(); err != nil {
		m.lofi.err = err.Error()
		return nil
	}
	if m.lofi.cmd != nil {
		if err := m.sendLofiCommand("loadfile", m.lofi.url, "replace"); err != nil {
			m.lofi.err = err.Error()
			return nil
		}
		m.lofi.status = lofiStatusLoading
		m.lofi.err = ""
		return nil
	}

	socketPath := filepath.Join(os.TempDir(), "student-manager-mpv.sock")
	_ = os.Remove(socketPath)
	cmd := exec.Command("mpv", "--no-video", "--really-quiet", "--input-ipc-server="+socketPath, "--idle=yes")
	if err := cmd.Start(); err != nil {
		m.lofi.err = fmt.Sprintf("mpv error: %v", err)
		return nil
	}
	m.lofi.cmd = cmd
	m.lofi.socketPath = socketPath
	m.lofi.status = lofiStatusLoading
	m.lofi.err = ""

	if err := m.sendLofiCommand("loadfile", m.lofi.url, "replace"); err != nil {
		m.lofi.err = err.Error()
	}

	return func() tea.Msg {
		err := cmd.Wait()
		return lofiExitMsg{err: err}
	}
}

func (m *Model) shutdownLofi() {
	if m.lofi.cmd == nil {
		m.lofi.status = lofiStatusStopped
		return
	}
	_ = m.sendLofiCommand("quit")
	if m.lofi.cmd.Process != nil {
		_ = m.lofi.cmd.Process.Kill()
	}
	m.lofi.cmd = nil
	m.lofi.socketPath = ""
	m.lofi.status = lofiStatusStopped
}

func (m *Model) validateLofi() error {
	if !m.lofi.enabled {
		return fmt.Errorf("Enable Lofi in Settings first.")
	}
	if strings.TrimSpace(m.lofi.url) == "" {
		return fmt.Errorf("Set a playlist URL in Settings.")
	}
	return nil
}

func loadLofiPlaylist(url string) tea.Cmd {
	url = strings.TrimSpace(url)
	if url == "" {
		return nil
	}
	return func() tea.Msg {
		tracks, err := fetchLofiPlaylist(url)
		return lofiPlaylistMsg{tracks: tracks, err: err}
	}
}

type ytPlaylist struct {
	Entries []struct {
		Title    string `json:"title"`
		Uploader string `json:"uploader"`
		Channel  string `json:"channel"`
	} `json:"entries"`
}

func fetchLofiPlaylist(url string) ([]models.LofiTrack, error) {
	cmd := exec.Command("yt-dlp", "--flat-playlist", "-J", url)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("yt-dlp error: %s", strings.TrimSpace(string(output)))
	}
	var payload ytPlaylist
	if err := json.Unmarshal(output, &payload); err != nil {
		return nil, fmt.Errorf("yt-dlp parse error: %v", err)
	}
	if len(payload.Entries) == 0 {
		return nil, fmt.Errorf("Playlist is empty.")
	}
	tracks := make([]models.LofiTrack, 0, len(payload.Entries))
	for _, entry := range payload.Entries {
		title := strings.TrimSpace(entry.Title)
		if title == "" {
			continue
		}
		note := strings.TrimSpace(entry.Uploader)
		if note == "" {
			note = strings.TrimSpace(entry.Channel)
		}
		tracks = append(tracks, models.LofiTrack{Title: title, Note: note})
	}
	if len(tracks) == 0 {
		return nil, fmt.Errorf("Playlist is empty.")
	}
	return tracks, nil
}

func (m *Model) sendLofiCommand(args ...any) error {
	conn, err := m.dialLofiSocket()
	if err != nil {
		return err
	}
	defer conn.Close()
	payload, err := json.Marshal(map[string]any{"command": args})
	if err != nil {
		return err
	}
	payload = append(payload, '\n')
	_, err = conn.Write(payload)
	return err
}

type mpvResponse struct {
	Error string      `json:"error"`
	Data  interface{} `json:"data"`
}

func (m *Model) isLofiPlaying() (bool, error) {
	resp, err := m.sendLofiRequest("get_property", "playback-time")
	if err != nil {
		return false, err
	}
	if resp.Error != "success" {
		return false, nil
	}
	val, ok := resp.Data.(float64)
	if !ok {
		return false, nil
	}
	return val > 0.05, nil
}

func (m *Model) sendLofiRequest(args ...any) (mpvResponse, error) {
	if m.lofi.socketPath == "" {
		return mpvResponse{}, fmt.Errorf("Player not running.")
	}
	conn, err := net.Dial("unix", m.lofi.socketPath)
	if err != nil {
		return mpvResponse{}, err
	}
	defer conn.Close()
	payload, err := json.Marshal(map[string]any{"command": args, "request_id": 1})
	if err != nil {
		return mpvResponse{}, err
	}
	payload = append(payload, '\n')
	if _, err := conn.Write(payload); err != nil {
		return mpvResponse{}, err
	}
	reader := bufio.NewReader(conn)
	line, err := reader.ReadBytes('\n')
	if err != nil {
		return mpvResponse{}, err
	}
	var resp mpvResponse
	if err := json.Unmarshal(line, &resp); err != nil {
		return mpvResponse{}, err
	}
	return resp, nil
}

func (m *Model) dialLofiSocket() (net.Conn, error) {
	if m.lofi.socketPath == "" {
		return nil, fmt.Errorf("Player not running.")
	}
	var lastErr error
	for i := 0; i < 6; i++ {
		conn, err := net.Dial("unix", m.lofi.socketPath)
		if err == nil {
			return conn, nil
		}
		lastErr = err
		time.Sleep(50 * time.Millisecond)
	}
	return nil, lastErr
}
