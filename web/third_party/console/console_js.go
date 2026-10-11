//go:build js

package console

// There is no real console in the browser; bubbletea falls back to plain
// reader/writer I/O when these report that.

func checkConsole(File) error { return ErrNotAConsole }

func newMaster(File) (Console, error) { return nil, ErrNotAConsole }
