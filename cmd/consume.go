package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

type consumeConfig struct {
	DefaultOptions []string `json:"default_options"`
	Scheme         string   `json:"scheme"`
	Host           string   `json:"host"`
	Path           string   `json:"path"`
}

var consumeCmd = &cobra.Command{
	Use:   "consume [host] as [session] [httpie args...]",
	Short: "HTTP client wrapper using httpie sessions",
	Long:  "HTTP client wrapper that uses httpie session files with optional GPG encryption",
	Args:  cobra.MinimumNArgs(4),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		baseDir := consumeBaseDir()

		switch len(args) {
		case 0:
			entries, err := os.ReadDir(baseDir)
			if err != nil {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			var hosts []string
			for _, entry := range entries {
				if entry.IsDir() {
					hosts = append(hosts, entry.Name())
				}
			}
			return hosts, cobra.ShellCompDirectiveNoFileComp
		case 1:
			return []string{"as"}, cobra.ShellCompDirectiveNoFileComp
		case 2:
			hostDir := filepath.Join(baseDir, args[0])
			entries, err := os.ReadDir(hostDir)
			if err != nil {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			var sessions []string
			for _, entry := range entries {
				name := entry.Name()
				if strings.HasSuffix(name, ".json.gpg") {
					sessions = append(sessions, strings.TrimSuffix(name, ".json.gpg"))
				} else if strings.HasSuffix(name, ".json") {
					sessions = append(sessions, strings.TrimSuffix(name, ".json"))
				}
			}
			return sessions, cobra.ShellCompDirectiveNoFileComp
		default:
			return nil, cobra.ShellCompDirectiveDefault
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		host := args[0]
		if args[1] != "as" {
			return fmt.Errorf("expected 'as' keyword, got '%s'", args[1])
		}
		sessionName := args[2]
		httpieArgs := args[3:]

		baseDir := consumeBaseDir()
		sessionFile := filepath.Join(baseDir, host, sessionName+".json")
		gpgSessionFile := sessionFile + ".gpg"

		var sessionData map[string]any
		var err error

		if _, err = os.Stat(gpgSessionFile); err == nil {
			sessionData, err = loadGPGSession(gpgSessionFile)
		} else if _, err = os.Stat(sessionFile); err == nil {
			sessionData, err = loadJSONSession(sessionFile)
		} else {
			return fmt.Errorf("%s does not exist", sessionFile)
		}
		if err != nil {
			return err
		}

		var config consumeConfig
		if consumeSection, ok := sessionData["consume"].(map[string]any); ok {
			if opts, ok := consumeSection["default_options"].([]any); ok {
				for _, opt := range opts {
					if s, ok := opt.(string); ok {
						config.DefaultOptions = append(config.DefaultOptions, s)
					}
				}
			}
			if scheme, ok := consumeSection["scheme"].(string); ok {
				config.Scheme = scheme
			}
			if h, ok := consumeSection["host"].(string); ok {
				config.Host = h
			}
			if path, ok := consumeSection["path"].(string); ok {
				config.Path = path
			}
		}

		if config.Scheme == "" {
			config.Scheme = "http"
		}
		if config.Host == "" {
			config.Host = host
		}
		config.Host = strings.TrimRight(config.Host, "/")
		config.Path = strings.Trim(config.Path, "/")

		delete(sessionData, "consume")

		tmpFile, err := os.CreateTemp("", "consume-session-*.json")
		if err != nil {
			return err
		}
		defer os.Remove(tmpFile.Name())

		if err := json.NewEncoder(tmpFile).Encode(sessionData); err != nil {
			tmpFile.Close()
			return err
		}
		tmpFile.Close()

		httpie := os.Getenv("CONSUME_HTTPIE")
		if httpie == "" {
			httpie = "http"
		}

		execArgs := []string{"--session-read-only", tmpFile.Name()}
		execArgs = append(execArgs, config.DefaultOptions...)

		for _, arg := range httpieArgs {
			if strings.HasPrefix(arg, "/") {
				url := fmt.Sprintf("%s://%s/", config.Scheme, config.Host)
				if config.Path != "" {
					url += config.Path + "/"
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
	},
}

func loadGPGSession(path string) (map[string]any, error) {
	command := exec.Command("gpg", "-q", "--no-tty", "--decrypt", path)
	command.Stdin = os.Stdin
	output, err := command.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt %s: %w", path, err)
	}

	var data map[string]any
	if err := json.Unmarshal(output, &data); err != nil {
		return nil, fmt.Errorf("failed to parse decrypted session: %w", err)
	}
	return data, nil
}

func loadJSONSession(path string) (map[string]any, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data map[string]any
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to parse session file: %w", err)
	}
	return data, nil
}

func consumeBaseDir() string {
	if baseDir := os.Getenv("CONSUME_BASEDIR"); baseDir != "" {
		return baseDir
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(homeDir, ".config", "httpie", "sessions")
}

func init() {
	rootCmd.AddCommand(consumeCmd)
}
