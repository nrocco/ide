package main

import "golang.org/x/sys/unix"

// setNonDumpable blocks ptrace from other same-user processes via
// PR_SET_DUMPABLE (arg 0 = not dumpable).
func setNonDumpable() error {
	return unix.Prctl(unix.PR_SET_DUMPABLE, 0, 0, 0, 0)
}
