package main

import "golang.org/x/sys/unix"

// tcgetsReq/tcsetsReq are the ioctl request numbers for reading/writing
// termios state. These differ between Linux and Darwin — see
// termiosreq_darwin.go.
const (
	tcgetsReq = unix.TCGETS
	tcsetsReq = unix.TCSETS
)
