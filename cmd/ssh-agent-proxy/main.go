// ssh-agent-split: creates one filtered socket per key in an upstream SSH agent.
//
// Usage:
//
//	ssh-agent-split --dir /tmp/agent-sockets
//
// Creates sockets named by key fingerprint:
//
//	/tmp/agent-sockets/SHA256_abc123...sock  (for key "my-github-key")
//
// A symlink named by comment is also created for convenience:
//
//	/tmp/agent-sockets/my-github-key.sock
package main

import (
	"crypto/sha256"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// ----- filter agent (same as before, keyed by blob) -------------------------

type filterAgent struct {
	upstream agent.ExtendedAgent
	keyBlob  string // the specific key blob we're exposing
}

func keyFingerprint(key agent.Key) string {
	digest := sha256.Sum256(key.Marshal())
	return "SHA256:" + base64.RawStdEncoding.EncodeToString(digest[:])
}

func (f *filterAgent) List() ([]*agent.Key, error) {
	keys, err := f.upstream.List()
	if err != nil {
		return nil, err
	}
	for _, k := range keys {
		if string(k.Marshal()) == f.keyBlob {
			return []*agent.Key{k}, nil
		}
	}
	return []*agent.Key{}, nil
}

func (f *filterAgent) Sign(key ssh.PublicKey, data []byte) (*ssh.Signature, error) {
	if string(key.Marshal()) != f.keyBlob {
		return nil, fmt.Errorf("agent failure: key not allowed")
	}
	return f.upstream.Sign(key, data)
}

func (f *filterAgent) SignWithFlags(key ssh.PublicKey, data []byte, flags agent.SignatureFlags) (*ssh.Signature, error) {
	if string(key.Marshal()) != f.keyBlob {
		return nil, fmt.Errorf("agent failure: key not allowed")
	}
	return f.upstream.SignWithFlags(key, data, flags)
}

func (f *filterAgent) Extension(extensionType string, contents []byte) ([]byte, error) {
	return f.upstream.Extension(extensionType, contents)
}

func (f *filterAgent) Add(_ agent.AddedKey) error {
	return fmt.Errorf("agent failure: not permitted")
}

func (f *filterAgent) Remove(_ ssh.PublicKey) error {
	return fmt.Errorf("agent failure: not permitted")
}

func (f *filterAgent) RemoveAll() error {
	return fmt.Errorf("agent failure: not permitted")
}

func (f *filterAgent) Lock(_ []byte) error {
	return fmt.Errorf("agent failure: not permitted")
}

func (f *filterAgent) Unlock(_ []byte) error {
	return fmt.Errorf("agent failure: not permitted")
}

func (f *filterAgent) Signers() ([]ssh.Signer, error) {
	return nil, fmt.Errorf("agent failure: not permitted")
}

// ----- per-key listener ------------------------------------------------------

type keyListener struct {
	key      agent.Key
	blob     string
	sockPath string
	listener net.Listener
}

func sanitizeComment(comment string) string {
	r := strings.NewReplacer("/", "_", " ", "_", ":", "_", "\\", "_")
	return r.Replace(comment)
}

func socketName(key agent.Key) string {
	fp := keyFingerprint(key)
	// Replace characters that are awkward in filenames
	safe := strings.ReplaceAll(fp, "/", "_")
	return safe + ".sock"
}

func startListener(key agent.Key, dir string, upstreamPath string) (*keyListener, error) {
	sockPath := filepath.Join(dir, socketName(key))

	// Remove stale socket
	if err := os.Remove(sockPath); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("remove stale socket: %w", err)
	}

	l, err := net.Listen("unix", sockPath)
	if err != nil {
		return nil, fmt.Errorf("listen: %w", err)
	}
	if err := os.Chmod(sockPath, 0600); err != nil {
		l.Close()
		return nil, fmt.Errorf("chmod: %w", err)
	}

	// Symlink by comment for convenience, e.g. my-github-key.sock
	if key.Comment != "" {
		linkPath := filepath.Join(dir, sanitizeComment(key.Comment)+".sock")
		_ = os.Remove(linkPath)
		if err := os.Symlink(sockPath, linkPath); err != nil {
			log.Printf("warning: could not create symlink %s: %v", linkPath, err)
		}
	}

	kl := &keyListener{
		key:      key,
		blob:     string(key.Marshal()),
		sockPath: sockPath,
		listener: l,
	}

	go kl.serve(upstreamPath)

	log.Printf("listening: %s  (%s  %s)", sockPath, keyFingerprint(key), key.Comment)

	return kl, nil
}

func (kl *keyListener) serve(upstreamPath string) {
	for {
		conn, err := kl.listener.Accept()
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				return
			}
			log.Printf("accept error on %s: %v", kl.sockPath, err)
			return
		}
		go kl.handleConn(conn, upstreamPath)
	}
}

func (kl *keyListener) handleConn(conn net.Conn, upstreamPath string) {
	defer conn.Close()

	upstreamConn, err := net.Dial("unix", upstreamPath)
	if err != nil {
		log.Printf("failed to connect to upstream: %v", err)
		return
	}
	defer upstreamConn.Close()

	fa := &filterAgent{
		upstream: agent.NewClient(upstreamConn),
		keyBlob:  kl.blob,
	}

	if err := agent.ServeAgent(fa, conn); err != nil && err != io.EOF {
		log.Printf("serve error: %v", err)
	}
}

func (kl *keyListener) close(dir string) {
	kl.listener.Close()
	os.Remove(kl.sockPath)
	if kl.key.Comment != "" {
		os.Remove(filepath.Join(dir, sanitizeComment(kl.key.Comment)+".sock"))
	}
}

// ----- watcher ---------------------------------------------------------------

type manager struct {
	upstreamPath string
	dir          string
	mu           sync.Mutex
	listeners    map[string]*keyListener // keyed by blob
}

func newManager(upstreamPath, dir string) *manager {
	return &manager{
		upstreamPath: upstreamPath,
		dir:          dir,
		listeners:    make(map[string]*keyListener),
	}
}

func (m *manager) sync() error {
	conn, err := net.Dial("unix", m.upstreamPath)
	if err != nil {
		return fmt.Errorf("dial upstream: %w", err)
	}
	defer conn.Close()

	keys, err := agent.NewClient(conn).List()
	if err != nil {
		return fmt.Errorf("list keys: %w", err)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Track which blobs are currently present
	present := make(map[string]bool)
	for _, k := range keys {
		blob := string(k.Marshal())
		present[blob] = true

		if _, exists := m.listeners[blob]; !exists {
			kl, err := startListener(*k, m.dir, m.upstreamPath)
			if err != nil {
				log.Printf("failed to start listener for %s: %v", keyFingerprint(*k), err)
				continue
			}
			m.listeners[blob] = kl
		}
	}

	// Remove listeners for keys that have disappeared
	for blob, kl := range m.listeners {
		if !present[blob] {
			log.Printf("key removed: %s  %s", keyFingerprint(kl.key), kl.key.Comment)
			kl.close(m.dir)
			delete(m.listeners, blob)
		}
	}

	return nil
}

func (m *manager) watch(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		if err := m.sync(); err != nil {
			log.Printf("sync error: %v", err)
		}
	}
}

func (m *manager) closeAll() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for blob, kl := range m.listeners {
		kl.close(m.dir)
		delete(m.listeners, blob)
	}
}

// ----- main ------------------------------------------------------------------

func main() {
	dir := flag.String("dir", "", "Directory to create sockets in (required)")
	upstream := flag.String("upstream", os.Getenv("SSH_AUTH_SOCK"), "Upstream agent socket (defaults to $SSH_AUTH_SOCK)")
	interval := flag.Duration("interval", 5*time.Second, "How often to poll for new/removed keys")
	flag.Parse()

	if *dir == "" {
		log.Fatal("--dir is required")
	}
	if *upstream == "" {
		log.Fatal("--upstream is required and $SSH_AUTH_SOCK is not set")
	}

	if err := os.MkdirAll(*dir, 0700); err != nil {
		log.Fatalf("failed to create socket dir: %v", err)
	}

	m := newManager(*upstream, *dir)

	// Initial sync
	if err := m.sync(); err != nil {
		log.Fatalf("initial sync failed: %v", err)
	}

	// Watch for changes
	go m.watch(*interval)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	log.Println("shutting down...")
	m.closeAll()
}
