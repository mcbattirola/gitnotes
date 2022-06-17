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
	os.Chdir(rootPath)
	assert.Nil(t, err)

	// expect the value returned to be projName
	r, err := getProjectRoot()
	assert.Nil(t, err)
	assert.Equal(t, rootName, r)
}

func testGetCurrentBranch(t *testing.T) {
	// create temp dir and branch
	// projName := t.TempDir()
	// git init into it
	// create branch

	// get current branch
	// expect it to be right

	// change branch and see if it changes too

	// test that it returns an error if run in not a git dir
}
