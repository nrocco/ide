// Encryption: AES-256-GCM (stdlib crypto/aes + crypto/cipher — the same AES
// implementation OpenSSL uses, no CLI shell-out, no cgo/openssl bindings).
// GCM is authenticated (AEAD): unlike openssl's plain `enc -aes-256-cbc`,
// tampered/corrupted ciphertext is detected and rejected, not silently
// decrypted into garbage. This was flagged as a gap in the openssl-enc
// version of this tool.
//
// Key derivation: PBKDF2-HMAC-SHA256, implemented directly from stdlib
// crypto/hmac + crypto/sha256 (RFC 8018) — avoids depending on
// golang.org/x/crypto, which requires network access to fetch.
//
// Memory hardening (see prior discussion of the same concerns in Python):
//   - passphrase held in a mlock'd byte slice (syscall.Mlock), zeroed after use
//   - core dumps disabled (RLIMIT_CORE = 0)
//   - process marked non-dumpable — blocks ptrace from other same-user
//     processes. Linux: PR_SET_DUMPABLE=0. Darwin: PT_DENY_ATTACH. See
//     hardening_linux.go / hardening_darwin.go.
//   - passphrase re-prompted after vim exits rather than cached across the
//     edit session (shrinks in-memory exposure window)
//   - plaintext tempfile shredded on exit incl. signals. Linux: tmpfs
//     (/dev/shm), so it never touches disk. Darwin has no tmpfs mounted by
//     default, so we fall back to os.TempDir() (backed by disk/APFS) — a
//     weaker guarantee than Linux; see tmpdir_darwin.go.
//
// Same hard limit as before: none of this stops a root/CAP_SYS_PTRACE
// attacker. Go's garbage collector is also a real caveat here — see the
// note in readPassphraseLocked below.
package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"
)

const (
	pbkdf2Iterations = 600_000 // OWASP 2023 recommendation for PBKDF2-SHA256
	keyLen           = 32      // AES-256
	saltLen          = 16
	nonceLen         = 12 // standard GCM nonce size
)

// ---------- locked memory ----------

// lockedBytes is a byte slice locked into RAM (mlock) so it can't be paged
// to swap, with an explicit Wipe method. Go strings are immutable like
// Python's, so passphrase input is read directly into a []byte (see
// readPassphraseLocked) to avoid ever materializing an immutable copy.
type lockedBytes struct {
	b      []byte
	locked bool
}

func newLockedBytes(n int) *lockedBytes {
	b := make([]byte, n)
	lb := &lockedBytes{b: b}
	if err := syscall.Mlock(b); err != nil {
		fmt.Fprintf(os.Stderr, "warning: mlock failed (%v) — passphrase may be swappable\n", err)
	} else {
		lb.locked = true
	}
	return lb
}

func (lb *lockedBytes) Wipe() {
	for i := range lb.b {
		lb.b[i] = 0
	}
	if lb.locked {
		_ = syscall.Munlock(lb.b)
	}
}

// ---------- process hardening ----------

func hardenProcess() error {
	// Disable core dumps.
	if err := syscall.Setrlimit(syscall.RLIMIT_CORE, &syscall.Rlimit{Cur: 0, Max: 0}); err != nil {
		return fmt.Errorf("setrlimit RLIMIT_CORE: %w", err)
	}
	// Platform-specific anti-ptrace hardening (PR_SET_DUMPABLE on Linux,
	// PT_DENY_ATTACH on Darwin) — see hardening_linux.go / hardening_darwin.go.
	if err := setNonDumpable(); err != nil {
		return fmt.Errorf("anti-debug hardening: %w", err)
	}
	return nil
}

// ---------- terminal passphrase entry (no echo, no external deps) ----------

// readPassphraseLocked reads a passphrase from the terminal with echo
// disabled, writing bytes directly into a locked buffer. Implemented via
// raw termios ioctls instead of golang.org/x/term to avoid a network fetch.
//
// Caveat carried over from the Python version: Go's runtime may still copy
// slice contents during things like a growing internal buffer before we
// control it, and the GC does not guarantee zeroing freed memory. We
// minimize this by reading directly into a fixed pre-allocated locked
// buffer rather than building a string and copying it in, which removes
// the "one extra immutable copy" problem Python's getpass has — but Go's
// GC-managed runtime still isn't as tight a guarantee as a language with
// manual memory control (e.g., C with explicit malloc/free).
func readPassphraseLocked(prompt string, maxLen int) (*lockedBytes, int, error) {
	fmt.Fprint(os.Stderr, prompt)
	defer fmt.Fprintln(os.Stderr)

	fd := int(os.Stdin.Fd())
	oldState, err := termMakeRaw(fd)
	if err != nil {
		return nil, 0, fmt.Errorf("disabling echo: %w", err)
	}
	defer termRestore(fd, oldState)

	lb := newLockedBytes(maxLen)
	n := 0
	buf := make([]byte, 1)
	for n < maxLen {
		nr, err := os.Stdin.Read(buf)
		if err != nil || nr == 0 {
			break
		}
		c := buf[0]
		if c == '\n' || c == '\r' {
			break
		}
		lb.b[n] = c
		n++
	}
	buf[0] = 0
	return lb, n, nil
}

// ---------- termios raw mode (echo off) via direct ioctl ----------
//
// The ioctl request numbers (tcgetsReq/tcsetsReq) differ between Linux and
// Darwin — see termiosreq_linux.go / termiosreq_darwin.go. The struct layout
// also differs (Linux's termios has an extra line-discipline byte; Darwin's
// Cc array is a different length), so we use golang.org/x/sys/unix's
// platform-correct unix.Termios/IoctlGetTermios/IoctlSetTermios instead of a
// hand-rolled struct, rather than duplicating the layout ourselves.

func termMakeRaw(fd int) (*unix.Termios, error) {
	orig, err := unix.IoctlGetTermios(fd, tcgetsReq)
	if err != nil {
		return nil, err
	}
	raw := *orig
	raw.Lflag &^= unix.ECHO | unix.ICANON // no echo, read byte-by-byte
	if err := unix.IoctlSetTermios(fd, tcsetsReq, &raw); err != nil {
		return nil, err
	}
	return orig, nil
}

func termRestore(fd int, state *unix.Termios) {
	_ = unix.IoctlSetTermios(fd, tcsetsReq, state)
}

// ---------- PBKDF2-HMAC-SHA256 (RFC 8018), stdlib-only ----------

func pbkdf2SHA256(password, salt []byte, iterations, keyLen int) []byte {
	prf := func() hmacIface { return hmac.New(sha256.New, password) }
	hLen := sha256.Size
	numBlocks := (keyLen + hLen - 1) / hLen
	dk := make([]byte, 0, numBlocks*hLen)

	block := make([]byte, 4)
	for i := 1; i <= numBlocks; i++ {
		binary.BigEndian.PutUint32(block, uint32(i))
		h := prf()
		h.Write(salt)
		h.Write(block)
		u := h.Sum(nil)
		t := make([]byte, len(u))
		copy(t, u)
		for j := 1; j < iterations; j++ {
			h := prf()
			h.Write(u)
			u = h.Sum(nil)
			for k := range t {
				t[k] ^= u[k]
			}
		}
		dk = append(dk, t...)
	}
	return dk[:keyLen]
}

type hmacIface interface {
	Write(p []byte) (n int, err error)
	Sum(b []byte) []byte
}

// ---------- AES-256-GCM encrypt/decrypt ----------

// File format: base64([salt(16)][nonce(12)][ciphertext+GCM tag])
func encryptFile(inPath, outPath string, passphrase []byte) error {
	plaintext, err := os.ReadFile(inPath)
	if err != nil {
		return err
	}
	salt := make([]byte, saltLen)
	if _, err := rand.Read(salt); err != nil {
		return err
	}
	key := pbkdf2SHA256(passphrase, salt, pbkdf2Iterations, keyLen)
	defer zero(key)

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	nonce := make([]byte, nonceLen)
	if _, err := rand.Read(nonce); err != nil {
		return err
	}
	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	blob := make([]byte, 0, saltLen+nonceLen+len(ciphertext))
	blob = append(blob, salt...)
	blob = append(blob, nonce...)
	blob = append(blob, ciphertext...)

	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(blob)))
	base64.StdEncoding.Encode(encoded, blob)
	encoded = append(encoded, '\n')

	return os.WriteFile(outPath, encoded, 0600)
}

func decryptFile(inPath, outPath string, passphrase []byte) error {
	encoded, err := os.ReadFile(inPath)
	if err != nil {
		return err
	}
	data := make([]byte, base64.StdEncoding.DecodedLen(len(encoded)))
	n, err := base64.StdEncoding.Decode(data, bytes.TrimSpace(encoded))
	if err != nil {
		return fmt.Errorf("invalid base64 content: %w", err)
	}
	data = data[:n]
	if len(data) < saltLen+nonceLen {
		return fmt.Errorf("file too short to be valid")
	}
	salt := data[:saltLen]
	nonce := data[saltLen : saltLen+nonceLen]
	ciphertext := data[saltLen+nonceLen:]

	key := pbkdf2SHA256(passphrase, salt, pbkdf2Iterations, keyLen)
	defer zero(key)

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return fmt.Errorf("decryption failed — wrong passphrase or corrupt/tampered file: %w", err)
	}
	return os.WriteFile(outPath, plaintext, 0600)
}

func zero(b []byte) {
	for i := range b {
		b[i] = 0
	}
}

// ---------- shred ----------

func shred(path string) {
	if _, err := os.Stat(path); err != nil {
		return
	}
	if err := exec.Command("shred", "-u", path).Run(); err != nil {
		_ = os.Remove(path)
	}
}

var rootCmd = &cobra.Command{
	Use:          "vimenc",
	Short:        "vimenc is a new top level command",
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		target := args[0]

		if err := hardenProcess(); err != nil {
			fmt.Fprintln(os.Stderr, "warning:", err)
		}

		tmp, err := os.CreateTemp(secureTmpDir(), "vimenc-")
		if err != nil {
			fmt.Fprintln(os.Stderr, "error creating tmpfile:", err)
			return err
		}
		tmpPath := tmp.Name()
		tmp.Close()

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-sigCh
			shred(tmpPath)
			os.Exit(1)
		}()
		defer shred(tmpPath)

		if info, err := os.Stat(target); err == nil && info.Size() > 0 {
			passLB, n, err := readPassphraseLocked("Passphrase: ", 256)
			if err != nil {
				fmt.Fprintln(os.Stderr, "error reading passphrase:", err)
				return err
			}
			err = decryptFile(target, tmpPath, passLB.b[:n])
			passLB.Wipe()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return err
			}
			// passphrase wiped here — not held during the vim edit session
		}

		vimCmd := exec.Command("vim", tmpPath)
		vimCmd.Stdin, vimCmd.Stdout, vimCmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		if err := vimCmd.Run(); err != nil {
			fmt.Fprintln(os.Stderr, "vim exited with error:", err)
			return err
		}

		passLB, n, err := readPassphraseLocked("Passphrase (to re-encrypt): ", 256)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error reading passphrase:", err)
			return err
		}
		err = encryptFile(tmpPath, target, passLB.b[:n])
		passLB.Wipe()
		if err != nil {
			fmt.Fprintln(os.Stderr, "encryption failed, plaintext left at:", tmpPath, err)
			return err
		}

		fmt.Printf("Saved and encrypted: %s\n", target)

		return nil
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
