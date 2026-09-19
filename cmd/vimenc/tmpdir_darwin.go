package main

import (
	"fmt"
	"os"
)

// secureTmpDir falls back to the regular OS temp directory, since Darwin
// has no tmpfs mounted by default the way Linux does with /dev/shm. This is
// a weaker guarantee than Linux: the plaintext scratch file can be paged to
// disk here. Callers should not treat this as equivalent hardening.
func secureTmpDir() string {
	fmt.Fprintln(os.Stderr, "warning: no tmpfs on Darwin — plaintext scratch file will be written to disk-backed temp storage")
	return os.TempDir()
}
