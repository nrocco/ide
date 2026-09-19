package main

import "golang.org/x/sys/unix"

// tcgetsReq/tcsetsReq are the ioctl request numbers for reading/writing
// termios state. Darwin uses BSD-style TIOCGETA/TIOCSETA rather than
// Linux's TCGETS/TCSETS — see termiosreq_linux.go.
const (
	tcgetsReq = unix.TIOCGETA
	tcsetsReq = unix.TIOCSETA
)
