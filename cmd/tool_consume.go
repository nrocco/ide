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

type sessionConfig struct {
	Consume consumeConfig              `json:"consume"`
	Session map[string]json.RawMessage `json:"-"`
}

func (c *sessionConfig) UnmarshalJSON(data []byte) error {
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
				name := strings.TrimSuffix(strings.TrimSuffix(entry.Name(), ".gpg"), ".json")
				sessions = append(sessions, name)
			}
			return sessions, cobra.ShellCompDirectiveNoFileComp
		default:
			return nil, cobra.ShellCompDirectiveDefault
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if args[1] != "as" {
			return fmt.Errorf("expected 'as' keyword, got '%s'", args[1])
		}

		host := args[0]
		sessionName := args[2]
		httpieArgs := args[3:]

		config, err := loadSessionFromFile(filepath.Join(consumeBaseDir(), host, sessionName+".json"))
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
	},
}

func loadSessionFromFile(path string) (*sessionConfig, error) {
	if _, err := os.Stat(path + ".gpg"); err == nil {
		return loadSessionFromFileGPG(path + ".gpg")
	} else if _, err := os.Stat(path); err == nil {
		return loadSessionFromFileJSON(path)
	}
	return nil, fmt.Errorf("%s does not exist", path)
}

func loadSessionFromFileGPG(path string) (*sessionConfig, error) {
	command := exec.Command("gpg", "-q", "--no-tty", "--decrypt", path)
	command.Stdin = os.Stdin
	output, err := command.Output()
	if err != nil {
		return nil, err
	}
	var data sessionConfig
	if err := json.Unmarshal(output, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func loadSessionFromFileJSON(path string) (*sessionConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var data sessionConfig
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return nil, err
	}
	return &data, nil
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
	toolCmd.AddCommand(consumeCmd)
}
