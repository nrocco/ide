package consume

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Config holds the consume-specific configuration from a session file
type Config struct {
	DefaultOptions []string `json:"default_options"`
	Scheme         string   `json:"scheme"`
	Host           string   `json:"host"`
	Path           string   `json:"path"`
	AccountID      string   `json:"account_id"`
}

// SessionConfig holds a parsed session file with its consume config and raw session data
type SessionConfig struct {
	Consume Config                     `json:"consume"`
	Session map[string]json.RawMessage `json:"-"`
}

// UnmarshalJSON implements custom unmarshalling to separate consume config from session data
func (c *SessionConfig) UnmarshalJSON(data []byte) error {
	raw := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if consumeRaw, ok := raw["consume"]; ok {
		if err := json.Unmarshal(consumeRaw, &c.Consume); err != nil {
			return err
		}
		delete(raw, "consume")
	}
	c.Session = raw
	return nil
}

// LoadSessionFromFile loads a session config from a .json file, resolving any op:// URIs via the op CLI
func LoadSessionFromFile(path string) (*SessionConfig, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("%s does not exist", path)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var first SessionConfig
	if err := json.Unmarshal(raw, &first); err != nil {
		return nil, err
	}
	if first.Consume.AccountID == "" {
		return &first, nil
	}
	resolved, err := resolveOPSecrets(raw, first.Consume.AccountID)
	if err != nil {
		return nil, err
	}
	var data SessionConfig
	if err := json.Unmarshal(resolved, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// resolveOPSecrets pipes JSON through `op inject` to replace op:// URIs with their 1Password values.
// If accountID is non-empty, passes --account to target a specific 1Password account.
// If the JSON contains no op:// references, op inject returns it unchanged.
func resolveOPSecrets(data []byte, accountID string) ([]byte, error) {
	args := []string{"inject"}
	if accountID != "" {
		args = append(args, "--account", accountID)
	}
	cmd := exec.Command("op", args...)
	cmd.Stdin = strings.NewReader(string(data))
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("op inject: %w", err)
	}
	return out, nil
}

// ListHosts returns all host directories in the base directory
func ListHosts() ([]string, error) {
	entries, err := os.ReadDir(BaseDir())
	if err != nil {
		return nil, err
	}
	var hosts []string
	for _, entry := range entries {
		if entry.IsDir() {
			hosts = append(hosts, entry.Name())
		}
	}
	return hosts, nil
}

// ListSessionsForHost returns all session names for the given host
func ListSessionsForHost(host string) ([]string, error) {
	entries, err := os.ReadDir(filepath.Join(BaseDir(), host))
	if err != nil {
		return nil, err
	}
	var sessions []string
	for _, entry := range entries {
		name := strings.TrimSuffix(entry.Name(), ".json")
		sessions = append(sessions, name)
	}
	return sessions, nil
}

// Run executes the consume command with the given arguments
func Run(host, sessionName string, httpieArgs []string) error {
	config, err := LoadSessionFromFile(filepath.Join(BaseDir(), host, sessionName+".json"))
	if err != nil {
		return err
	}
	if config.Consume.Scheme == "" {
		config.Consume.Scheme = "http"
	}
	if config.Consume.Host == "" {
		config.Consume.Host = host
	}
	config.Consume.Host = strings.TrimRight(config.Consume.Host, "/")
	config.Consume.Path = strings.Trim(config.Consume.Path, "/")

	tmpFile, err := os.CreateTemp("", "consume-session-*.json")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())

	if err := json.NewEncoder(tmpFile).Encode(config.Session); err != nil {
		tmpFile.Close()
		return err
	}
	tmpFile.Close()

	httpie := os.Getenv("CONSUME_HTTPIE")
	if httpie == "" {
		httpie = "http"
	}

	execArgs := []string{"--session-read-only", tmpFile.Name()}
	execArgs = append(execArgs, config.Consume.DefaultOptions...)

	for _, arg := range httpieArgs {
		if strings.HasPrefix(arg, "/") {
			url := fmt.Sprintf("%s://%s/", config.Consume.Scheme, config.Consume.Host)
			if config.Consume.Path != "" {
				url += config.Consume.Path + "/"
			}
			url += strings.TrimLeft(arg, "/")
			execArgs = append(execArgs, url)
		} else {
			execArgs = append(execArgs, arg)
		}
	}

	command := exec.Command(httpie, execArgs...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	return command.Run()
}

// BaseDir returns the base directory for consume session files
func BaseDir() string {
	if baseDir := os.Getenv("CONSUME_BASEDIR"); baseDir != "" {
		return baseDir
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(homeDir, ".config", "httpie", "sessions")
}
