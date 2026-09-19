package main

// secureTmpDir returns a tmpfs-backed directory so the plaintext scratch
// file never touches disk or swap.
func secureTmpDir() string {
	return "/dev/shm"
}
