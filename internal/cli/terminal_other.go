//go:build !darwin && !dragonfly && !freebsd && !linux && !netbsd && !openbsd

package cli

import "os"

func fileTerminalWidth(_ *os.File) (int, bool) {
	return 0, false
}
