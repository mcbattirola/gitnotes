package gn

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/mcbattirola/gitnotes/pkg/errflags"
)

type GN struct {
	// Editor is the name of the binary of the text editor
	Editor string
	// NotesPath is the path in which the notes are stored
	NotesPath string
	// Project is the name of the project for the notes
	Project string
	// Branch is the name of the branch for the notes
	Branch string
	author Author
}

// New creates a new GN
// with required internal fields set
func New() *GN {
	// read git global configs
	// this is necessary because go-git wont read it
	// this can be removed if se use only syscalls instead of
	// using go-git
	return &GN{
		author: readGlobalGitAuthor(),
	}
}

// Edit opens the user's current project and branch on
// the selected editor. The behaviour of this method depends on the
// working directory, since it uses the current dir to find the project's name
func (gn *GN) Edit() error {
	var err error

	// run `git init` into notes path
	// we can still procceed if it errors
	// TODO log this error if in debug/verbose mode
	gn.init()

	project := gn.Project
	// if didn't received project name, find it
	if project == "" {
		// read current project name and branch
		project, err = getProjectRoot()
		if err != nil {
			return err
		}
	}

	branch := gn.Branch
	// if didn't received branch name, use current working branch
	if branch == "" {
		// get user working repo
		dir, err := os.Getwd()
		if err != nil {
			return err
		}
		r, err := git.PlainOpen(dir)
		if err != nil {
			return err
		}

		branch, err = getCurrentBranch(r)
		if err != nil {
			return err
		}
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
				errflags.Flag(err, errflags.NotAuthorized)
			}
			return err
		}
	} else if err != nil {
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

// init initializes a git repo into the notes path.
// Running git init in an existing repository is safe. It will not overwrite things that are already there
func (gn *GN) init() error {
	_, err := exec.Command("git", "init", "--quiet", gn.NotesPath).Output()
	if err != nil {
		return err
	}

	return nil
}

// Push pushes git notes to the remote repository
func (gn *GN) Push() error {
	// run `git init` into notes path
	// we can still procceed if it errors
	// TODO log this error if in debug/verbose mode
	gn.init()

	r, err := git.PlainOpen(gn.NotesPath)
	if err != nil {
		return err
	}

	// create an upstream branch if it doesn't exist

	// run `git add .`
	w, err := r.Worktree()
	if err != nil {
		return err
	}

	_, err = w.Add(".")
	if err != nil {
		return err
	}

	_, err = w.Commit(fmt.Sprintf("Update notes - %s", time.Now().Local().String()), &git.CommitOptions{
		Author: &object.Signature{When: time.Now()},
	})
	if err != nil {
		return err
	}

	if _, err := r.Remote("origin"); err != nil {
		if err == git.ErrRemoteNotFound {
			return errflags.Flag(err, errflags.NoRemote)
		}
		return err
	}

	// push
	// TODO make branch name variable / read it from current branch
	cmd := exec.Command("git", "-C", gn.NotesPath, "push", "--set-upstream", "origin", "master")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Printf("push failed: %s\n", err.Error())
		return err
	}
	// err = r.Push(&git.PushOptions{
	// 	Auth: &http.BasicAuth{
	// 		Username: "username",
	// 		Password: "password",
	// 	},
	// })
	// if err != nil && err != git.NoErrAlreadyUpToDate {
	// 	if err == git.ErrRemoteNotFound {
	// 		return errflags.Flag(err, errflags.NoRemote)
	// 	}
	// 	return err
	// }

	return nil
}

func (gn *GN) AddOrigin(url string) error {
	r, err := git.PlainOpen(gn.NotesPath)
	if err != nil {
		return err
	}

	_, err = r.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{url},
	})

	return err
}

// Path returns the path in which the notes are stored
func (gn *GN) Path() string {
	return gn.NotesPath
}
