//go:build darwin || dragonfly || freebsd || linux || netbsd || openbsd

package cli

import (
	"os"
	"syscall"
	"unsafe"
)

type winsize struct {
	rows uint16
	cols uint16
	x    uint16
	y    uint16
}

func fileTerminalWidth(file *os.File) (int, bool) {
	var size winsize
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		file.Fd(),
		syscall.TIOCGWINSZ,
		uintptr(unsafe.Pointer(&size)),
	)
	if errno != 0 || size.cols == 0 {
		return 0, false
	}
	return int(size.cols), true
}
