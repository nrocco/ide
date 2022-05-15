package ide

import (
	"os"
)

// AutoDetectLanguage guesses the language of a git repository based on some
// standard files and defaults to plain
func (project *Project) AutoDetectLanguage() string {
	if _, err := os.Stat("setup.py"); err == nil {
		return "python"
	} else if _, err := os.Stat("composer.json"); err == nil {
		return "php"
	} else if _, err := os.Stat("package.json"); err == nil {
		return "javascript"
	} else if _, err := os.Stat("go.mod"); err == nil {
		return "go"
	} else if _, err := os.Stat("Gemfile"); err == nil {
		return "ruby"
	} else if _, err := os.Stat("metadata.rb"); err == nil {
		return "chef"
	}

	return "plain"
}

// Language returns the language of the ide project as stored in .git/config file
func (project *Project) Language() string {
	return project.config.Raw.Section("ide").Option("language")
}

// SetLanguage stores the given language in the .git/config file of the ide project
func (project *Project) SetLanguage(language string) error {
	project.config.Raw.Section("ide").SetOption("language", language)

	return project.repository.Storer.SetConfig(project.config)
}
