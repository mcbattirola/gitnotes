package gn

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type GN struct {
	Editor    string
	NotesPath string
	Project   string
	Branch    string
}

//Edit opens the user's current project and branch on
//the selected editor. The behaviour of this method depends on the
//working directory, since it uses the current dir to find the project's name
func (gn *GN) Edit() error {
	// if received project or branch name, use it
	if gn.Project != "" {
		return gn.edit(gn.Project, gn.Branch)
	}

	// read current project name and branch
	project, err := getProjectRoot()
	if err != nil {
		return err
	}

	// get user working repo
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	// if received branch, use it instead of checking the current one
	if gn.Branch != "" {
		return gn.edit(project, gn.Branch)
	}

	r, err := git.PlainOpen(dir)
	if err != nil {
		return err
	}

	branch, err := getCurrentBranch(r)
	if err != nil {
		return err
	}

	return gn.edit(project, branch)
}

// edit opens a specific project/branch
// on the selected editor
// If project is empty, uses current project
func (gn *GN) edit(project string, branch string) error {
	_, err := os.Stat(gn.NotesPath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(gn.NotesPath, os.ModeDir|0700)
		if err != nil {
			if errors.Is(err, fs.ErrPermission) {
				fmt.Fprintf(os.Stderr, "No permision to create directory %s", gn.NotesPath)
			}
			return err
		}
	} else if err != nil {
		// TODO improve error handling
		return err
	}

	projectPath := fmt.Sprintf("%s/%s", gn.NotesPath, project)
	notePath := fmt.Sprintf("%s/%s", projectPath, branch)

	// we make the directory with the notePath instead of project path because
	// a branch name may contain slashes. In that case, we want to make the full path
	// and the slashes in branch name will become directories (which is ok for now)
	if err := os.MkdirAll(filepath.Dir(notePath), os.ModeDir|0700); err != nil {
		return err
	}

	// TODO if file doesn't exist, create it with a header
	_, err = os.OpenFile(notePath, os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	editor := gn.Editor
	if editor == "" {
		editor = "vi"
	}

	cmd := exec.Command(editor, notePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// getProjectRoot runs git through a syscall to get the top level directory
// we do it this way because go-git does not implement rev-parse
func getProjectRoot() (string, error) {
	path, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", err
	}

	s := strings.Split(strings.TrimSpace(string(path)), "/")
	return s[len(s)-1], nil
}

func getCurrentBranch(r *git.Repository) (string, error) {
	h, err := r.Reference(plumbing.HEAD, false)
	if err != nil {
		return "", err
	}

	target := h.Target()

	s := strings.Split(string(target), "refs/heads/")
	if len(s) < 2 {
		return "", errors.New("couldn't find project's root")
	}
	return s[1], nil
}
