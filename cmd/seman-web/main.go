//go:build js && wasm

// Command seman-web runs seman in the browser. It is the same Bubble Tea
// program as cmd/seman, with the terminal replaced by xterm.js and the JSON
// file replaced by the browser's localStorage.
package main

import (
	"encoding/json"
	"io"
	"syscall/js"

	"github.com/OmarIbrahim22/seman/internal/app"
	"github.com/OmarIbrahim22/seman/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

const storageKey = "seman.semester"

// browserStore implements storage.Store on top of localStorage.
type browserStore struct{}

func (browserStore) Load() (storage.SemesterData, bool, error) {
	raw := js.Global().Get("localStorage").Call("getItem", storageKey)
	if raw.IsNull() || raw.IsUndefined() {
		return storage.SemesterData{}, false, nil
	}
	var data storage.SemesterData
	if err := json.Unmarshal([]byte(raw.String()), &data); err != nil {
		return storage.SemesterData{}, true, err
	}
	return data, true, nil
}

func (browserStore) Save(data storage.SemesterData) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}
	js.Global().Get("localStorage").Call("setItem", storageKey, string(payload))
	return nil
}

// keyReader hands keystrokes from xterm.js to Bubble Tea.
type keyReader struct {
	keys chan []byte
	rest []byte
}

func (r *keyReader) Read(p []byte) (int, error) {
	if len(r.rest) == 0 {
		chunk, ok := <-r.keys
		if !ok {
			return 0, io.EOF
		}
		r.rest = chunk
	}
	n := copy(p, r.rest)
	r.rest = r.rest[n:]
	return n, nil
}

// termWriter sends rendered frames to xterm.js.
type termWriter struct{}

func (termWriter) Write(p []byte) (int, error) {
	js.Global().Call("semanWrite", string(p))
	return len(p), nil
}

func main() {
	// There is no terminal to query, so state what xterm.js supports.
	lipgloss.SetColorProfile(termenv.TrueColor)
	lipgloss.SetHasDarkBackground(true)

	var program *tea.Program
	size := tea.WindowSizeMsg{Width: 120, Height: 34}
	input := &keyReader{keys: make(chan []byte, 256)}

	js.Global().Set("semanInput", js.FuncOf(func(_ js.Value, args []js.Value) any {
		select {
		case input.keys <- []byte(args[0].String()):
		default: // drop keys if the app is not reading
		}
		return nil
	}))
	js.Global().Set("semanResize", js.FuncOf(func(_ js.Value, args []js.Value) any {
		size = tea.WindowSizeMsg{Width: args[0].Int(), Height: args[1].Int()}
		if p := program; p != nil {
			go p.Send(size)
		}
		return nil
	}))
	js.Global().Call("semanReady")

	// Quitting restarts the app, since a browser tab has no shell to return to.
	for {
		store := browserStore{}
		data, found, err := store.Load()
		if err != nil || !found {
			data = demoData()
		}
		program = tea.NewProgram(
			app.NewModel(store, data, true),
			tea.WithAltScreen(),
			tea.WithInput(input),
			tea.WithOutput(termWriter{}),
		)
		start := size
		go program.Send(start)
		if _, err := program.Run(); err != nil {
			js.Global().Call("semanExited", err.Error())
			return
		}
		js.Global().Call("semanExited", "")
	}
}
