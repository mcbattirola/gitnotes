package config

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/mcbattirola/gitnotes/pkg/gn"
)

// TODO use some package for testing

func TestReadConfigFile(t *testing.T) {
	// it creates a config file if it doesn't exist
	gn := gn.GN{}
	testDir := t.TempDir()
	fileName := "test.conf"
	err := ReadConfigFile(&gn, testDir, fileName)
	if err != nil {
		log.Fatal(err)
	}

	filePath := fmt.Sprintf("%s/%s", testDir, fileName)
	if _, err := os.Stat(filePath); err != nil {
		log.Fatalf("error reading config file: %s", err.Error())
	}

	// expect default values
	if gn.Editor != "vi" {
		log.Fatalf("default editor expected to be vi but got '%s'", gn.Editor)
	}
	if gn.NotesPath != os.ExpandEnv("$HOME/gitnotes") {
		log.Fatalf("default notes path should be %s but got %s", "$HOME/gitnotes", gn.NotesPath)
	}

	// it reads the fields from the config file if it exist

}
