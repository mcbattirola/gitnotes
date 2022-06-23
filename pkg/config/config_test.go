package config

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/mcbattirola/gitnotes/pkg/gn"
	"github.com/stretchr/testify/assert"
)

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
	assert.Equal(t, "vi", gn.Editor)
	assert.Equal(t, os.ExpandEnv("$HOME/gitnotes"), gn.NotesPath)
	assert.Equal(t, false, gn.AlwaysCommit)
}
