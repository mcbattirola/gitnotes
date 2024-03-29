package config

import (
	"bufio"
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/mcbattirola/gitnotes/pkg/gn"
)

//go:embed gn.conf
var defaultConfig string

// ReadConfigFile reads a config file and sets the values found
// into gn
func ReadConfigFile(gn *gn.GN, path string, fileName string) error {
	// create config file if it doesnt exist
	filePath := fmt.Sprintf("%s/%s", path, fileName)
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			err = createConfigFile(path, fileName)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	// scan config file line by line
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), "=")
		if len(s) < 2 {
			continue
		}

		switch s[0] {
		case "editor":
			gn.Editor = parseInput(s[1])
		case "notes":
			gn.NotesPath = parseInput(s[1])
		case "always-commit":
			if parseInput(s[1]) == "true" {
				gn.AlwaysCommit = true
			}

		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan file error: %v", err)
	}

	return nil
}

// parseInput reads the input, removes comments and
// trim whitespaces
func parseInput(i string) string {
	s := strings.Split(i, "#")[0]
	s = strings.TrimSpace(s)
	return os.ExpandEnv(s)
}

func createConfigFile(configPath string, fileName string) error {
	if err := os.MkdirAll(configPath, os.ModeDir|0700); err != nil {
		return err
	}

	f, err := os.OpenFile(fmt.Sprintf("%s/%s", configPath, fileName), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.Write([]byte(defaultConfig))
	if err != nil {
		return err
	}

	return nil
}
