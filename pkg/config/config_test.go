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
	assert.Equal(t, "vim", gn.Editor)
	assert.Equal(t, os.ExpandEnv("$HOME/gitnotes"), gn.NotesPath)
	assert.Equal(t, false, gn.AlwaysCommit)
}

func TestParseInput(t *testing.T) {
	tt := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normal input",
			input:    "value",
			expected: "value",
		},
		{
			name:     "input with trailing whitespaces",
			input:    "  value   ",
			expected: "value",
		},
		{
			name:     "input whith comment",
			input:    "value#comment",
			expected: "value",
		},
		{
			name:     "input whith comment and whitespaces",
			input:    " value # comment",
			expected: "value",
		},
		{
			name:     "input whith multiple comments and whitespaces",
			input:    " value    #    # multiple comments ",
			expected: "value",
		},
		{
			name:     "input with env vars",
			input:    "$HOME",
			expected: os.ExpandEnv("$HOME"),
		},
		{
			name:     "input with env vars, comments and spaces",
			input:    " value$HOME  # comment ",
			expected: "value" + os.ExpandEnv("$HOME"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, parseInput(tc.input))
		})
	}
}
