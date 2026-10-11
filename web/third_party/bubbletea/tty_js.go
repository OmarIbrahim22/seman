//go:build js

package tea

import (
	"errors"
	"os"
)

// In the browser, input and output are wired to xterm.js through
// WithInput/WithOutput, so there is no TTY to configure.

func (p *Program) initInput() error { return nil }

func (p *Program) restoreInput() error { return nil }

func openInputTTY() (*os.File, error) {
	return nil, errors.New("bubbletea: no TTY in the browser")
}

// Resizes are delivered from JavaScript as WindowSizeMsg instead of SIGWINCH.
func (p *Program) listenForResize(done chan struct{}) {
	<-p.ctx.Done()
	close(done)
}
