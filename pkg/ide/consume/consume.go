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

// LoadSessionFromFile loads a session config from a .json or .gpg file
func LoadSessionFromFile(path string) (*SessionConfig, error) {
	if _, err := os.Stat(path + ".gpg"); err == nil {
		return loadSessionFromFileGPG(path + ".gpg")
	} else if _, err := os.Stat(path); err == nil {
		return loadSessionFromFileJSON(path)
	}
	return nil, fmt.Errorf("%s does not exist", path)
}

func loadSessionFromFileGPG(path string) (*SessionConfig, error) {
	command := exec.Command("gpg", "-q", "--no-tty", "--decrypt", path)
	command.Stdin = os.Stdin
	output, err := command.Output()
	if err != nil {
		return nil, err
	}
	var data SessionConfig
	if err := json.Unmarshal(output, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func loadSessionFromFileJSON(path string) (*SessionConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var data SessionConfig
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return nil, err
	}
	return &data, nil
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
		name := strings.TrimSuffix(strings.TrimSuffix(entry.Name(), ".gpg"), ".json")
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
