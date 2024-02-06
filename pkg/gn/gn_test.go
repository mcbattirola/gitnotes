package gn

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindProject(t *testing.T) {
	if os.Getenv("GN_TEST_INTEAGRATION") != "TRUE" {
		t.Skip("skipping integration test TestFindProject")
	}

	// change dir to source code path inside test container
	err := os.Chdir(os.Getenv("GN_CODE_PATH"))
	assert.NoError(t, err)

	gn := New(false)

	// find current project when none is specified
	p, err := gn.findProject()
	assert.NoError(t, err)
	assert.Equal(t, "gitnotes", p)

	// uses the provided project
	gn.Project = "test_project"
	p, err = gn.findProject()
	assert.NoError(t, err)
	assert.Equal(t, "test_project", p)

	// assert that find project works from a inner directory
	err = os.Chdir("./test")
	assert.NoError(t, err)

	gn.Project = "test_project"
	p, err = gn.findProject()
	assert.NoError(t, err)
	assert.Equal(t, "test_project", p)
}

func TestFindBranch(t *testing.T) {
	if os.Getenv("GN_TEST_INTEAGRATION") != "TRUE" {
		t.Skip("skipping integration test TestFindBranch")
	}

	// change dir to source code path inside test container
	err := os.Chdir(os.Getenv("GN_CODE_PATH"))
	assert.NoError(t, err)

	gn := New(false)

	// find current branch when none is specified
	p, err := gn.findBranch()
	assert.NoError(t, err)
	assert.Equal(t, os.Getenv("GN_CURRENT_BRANCH"), p)

	// uses the provided branch
	gn.Branch = "test_branch"
	p, err = gn.findBranch()
	assert.NoError(t, err)
	assert.Equal(t, "test_branch", p)

	// assert that find branch works from a inner directory
	err = os.Chdir("./test")
	assert.NoError(t, err)
	gn.Branch = "test_branch"
	p, err = gn.findBranch()
	assert.NoError(t, err)
	assert.Equal(t, "test_branch", p)
}

func TestCreateNotesPath(t *testing.T) {
	tempDir := t.TempDir()
	gn := New(false)
	gn.NotesPath = fmt.Sprintf("%s/%s", tempDir, "gitnotes")

	err := gn.createNotesPath()
	assert.NoError(t, err)

	_, err = os.Stat(gn.NotesPath)
	assert.NoError(t, err)
}

func TestReadNote(t *testing.T) {
	// create a temp dir
	tempDir := t.TempDir()
	gn := New(false)
	project := "test-project"
	branch := "test-branch"

	gn.NotesPath = fmt.Sprintf("%s/%s", tempDir, "gitnotes")
	gn.Project = project
	gn.Branch = branch

	// add a note to it
	note := fmt.Sprintf("%s/%s", gn.NotesPath, project)
	err := os.MkdirAll(note, os.ModeDir|0700)
	assert.NoError(t, err)

	content := "this is the note content"
	err = os.WriteFile(note+"/"+branch, []byte(content), os.ModeDir|0700)
	assert.NoError(t, err)

	noteContent, err := gn.ReadNote()
	assert.NoError(t, err)
	assert.Equal(t, content, noteContent)
}
