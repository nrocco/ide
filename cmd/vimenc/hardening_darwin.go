package main

import "golang.org/x/sys/unix"

// ptDenyAttach is Darwin's PT_DENY_ATTACH request (sys/ptrace.h) — it has no
// named constant in golang.org/x/sys/unix.
const ptDenyAttach = 31

// setNonDumpable blocks ptrace/debugger attachment from other processes via
// ptrace(PT_DENY_ATTACH, 0, 0, 0) — Darwin's equivalent of Linux's
// PR_SET_DUMPABLE.
func setNonDumpable() error {
	_, _, errno := unix.Syscall(unix.SYS_PTRACE, ptDenyAttach, 0, 0)
	if errno != 0 {
		return errno
	}
	return nil
}
