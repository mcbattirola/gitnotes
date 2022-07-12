package gn

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetProjectRoot(t *testing.T) {
	// create temp dir
	rootPath := t.TempDir()

	// rootPath will be something like /tmp/hash/number
	// so the expected root of the project is 'number'
	rootName := filepath.Base(rootPath)

	// start a git repo in it
	_, err := exec.Command("git", "init", rootPath).Output()
	assert.Nil(t, err)

	// change working directory to rootPath
	err = os.Chdir(rootPath)
	assert.Nil(t, err)

	// expect the value returned to be projName
	r, err := getProjectRoot()
	assert.Nil(t, err)
	assert.Equal(t, rootName, r)
}

func TestGetAuthorFromGitConfig(t *testing.T) {
	tt := []struct {
		name     string
		config   []byte
		expected Author
	}{
		{
			name:     "empty input should produce an empty author",
			config:   []byte(""),
			expected: Author{},
		},
		{
			name:   "config with email should populate author's email",
			config: []byte("code.editor=nano\nuser.email=test@example.com"),
			expected: Author{
				Email: "test@example.com",
			},
		},
		{
			name:   "config with email should populate author's name",
			config: []byte("code.editor=nano\nuser.name=testuser"),
			expected: Author{
				Name: "testuser",
			},
		},
		{
			name:   "config with email and name should populate both fields",
			config: []byte("code.editor=nano\nuser.name=testuser\nuser.email=test@example.com"),
			expected: Author{
				Name:  "testuser",
				Email: "test@example.com",
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, getAuthorFromGitConfig(tc.config))
		})
	}
}
