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
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/mcbattirola/gitnotes/pkg/errflags"
	"github.com/mcbattirola/gitnotes/pkg/log"
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
	// commit message is the message to be used on commit
	CommitMessage string
	// AlwaysCommit indicates if all note edit should be commited
	AlwaysCommit bool
	// RemoteURL is the URL to the remote repository
	RemoteURL string
	author    Author
	log       log.Logger
}

// New creates a new GN
// with required internal fields set
func New(debug bool) *GN {
	log := log.New(debug)

	a, err := readGlobalGitAuthor()
	if err != nil {
		// it is ok to ignore this error
		log.Debug("failed to read global git author: %s", err.Error)
	}
	return &GN{
		author: a,
		log:    log,
	}
}

// Edit opens the user's current project and branch on
// the selected editor. The behaviour of this method depends on the
// working directory, since it uses the current dir to find the project's name
func (gn *GN) Edit() error {
	var err error

	if gn.AlwaysCommit {
		defer func() {
			err := gn.Commit()
			if err != nil {
				gn.log.Info("failed to commit: %s", err.Error())
			}
		}()
	}

	// run `git init` into notes path
	// we can still procceed if it errors
	if err := gn.init(); err != nil {
		gn.log.Debug("failed to init: %s", err.Error())
	}

	project, err := gn.findProject()
	if err != nil {
		return err
	}

	branch, err := gn.findBranch()
	if err != nil {
		return err
	}

	return gn.edit(project, branch)
}

// findProject returns the name of the project
// if no project is set, it finds and returns the current
// working directory project
func (gn *GN) findProject() (string, error) {
	project := gn.Project
	// if didn't received project name, find it
	if project == "" {
		// read current project name and branch
		project, err := getProjectRoot()
		if err != nil {
			gn.log.Debug("could not find project root: %s", err.Error())
			return "", err
		}
		gn.log.Debug("found project root: %s", project)
		return project, nil
	}

	return project, nil
}

func (gn *GN) findBranch() (string, error) {
	branch := gn.Branch
	// if didn't received branch name, use current working branch
	if branch == "" {
		// get user working repo
		dir, err := os.Getwd()
		if err != nil {
			gn.log.Debug("could not find branch: %s", err.Error())
			return "", err
		}
		r, err := git.PlainOpenWithOptions(dir, &git.PlainOpenOptions{DetectDotGit: true})
		if err != nil {
			gn.log.Debug("could not open repository to look for branch: %s", err.Error())
			return "", err
		}

		branch, err = getCurrentBranch(r)
		if err != nil {
			return "", err
		}
		gn.log.Debug("found branch: %s", branch)
		return branch, nil
	}

	return branch, nil
}

// edit opens a specific project/branch on the selected editor
// If project is empty, uses current project
func (gn *GN) edit(project string, branch string) error {
	gn.log.Debug("editing project %s and branch %s", project, branch)

	err := gn.createNotesPath()
	if err != nil {
		return err
	}

	notePath := getNotePath(gn.NotesPath, project, branch)
	gn.log.Debug("note path: %s", notePath)

	// make the directory with the notePath instead of project path because
	// a branch name may contain slashes. In that case, we want to make the full path
	// and the slashes in branch name will become directories (which is ok for now)
	if err := os.MkdirAll(filepath.Dir(notePath), os.ModeDir|0700); err != nil {
		return err
	}

	gn.log.Debug("opening note file %s", notePath)
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

// createNotesPath creates the note file if it doesn't exist
// it does nothing if the file already exists
func (gn *GN) createNotesPath() error {
	_, err := os.Stat(gn.NotesPath)
	if os.IsNotExist(err) {
		gn.log.Debug("notes path %s does not exist", gn.NotesPath)
		err := os.MkdirAll(gn.NotesPath, os.ModeDir|0700)
		if err != nil {
			if errors.Is(err, fs.ErrPermission) {
				return errflags.Flag(err, errflags.NotAuthorized)
			}
			return err
		}
		gn.log.Debug("notes path created without errors")
	} else if err != nil {
		return err
	}

	return nil
}

// init initializes a git repo into the notes path.
// Running git init in an existing repository is safe. It will not overwrite things that are already there
func (gn *GN) init() error {
	_, err := git.PlainInit(gn.NotesPath, false)
	if err != nil {
		return err
	}

	return nil
}

// Push pushes git notes to the remote repository
func (gn *GN) Push() error {
	// run `git init` into notes path
	// we can still procceed if it errors
	if err := gn.init(); err != nil {
		gn.log.Debug("failed to init: %s", err.Error())
	}

	r, err := git.PlainOpen(gn.NotesPath)
	if err != nil {
		return err
	}

	// run `git add .`
	w, err := r.Worktree()
	if err != nil {
		return err
	}

	_, err = w.Add(".")
	if err != nil {
		return err
	}

	// create commit
	if gn.CommitMessage == "" {
		gn.CommitMessage = fmt.Sprintf("Update notes - %s", time.Now().Local().String())
	}
	err = gn.commit(gn.CommitMessage, w)
	if err != nil {
		return err
	}

	if err = gn.checkAndAddOrigin(r); err != nil {
		return err
	}

	err = r.Push(&git.PushOptions{})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		if err == git.ErrRemoteNotFound {
			return errflags.Flag(err, errflags.NoRemote)
		}
		return err
	}

	return nil
}

// Pull pushes git notes to the remote repository
func (gn *GN) Pull() error {
	if err := gn.init(); err != nil {
		gn.log.Debug("failed to init: %s", err.Error())
	}

	r, err := git.PlainOpen(gn.NotesPath)
	if err != nil {
		return err
	}

	err = gn.checkAndAddOrigin(r)
	if err != nil {
		return err
	}

	// Get the working directory for the repository
	w, err := r.Worktree()
	if err != nil {
		return err
	}
	err = w.Pull(&git.PullOptions{RemoteName: "origin", ReferenceName: plumbing.Master})
	if err != git.NoErrAlreadyUpToDate {
		return err
	}

	return nil
}

func (gn *GN) Commit() error {
	r, err := git.PlainOpen(gn.NotesPath)
	if err != nil {
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}

	if gn.CommitMessage == "" {
		gn.CommitMessage = fmt.Sprintf("Update notes - %s", time.Now().Local().String())
	}

	return gn.commit(gn.CommitMessage, w)
}

func (gn *GN) commit(msg string, w *git.Worktree) error {
	_, err := w.Commit(msg, &git.CommitOptions{
		Author: &object.Signature{
			When:  time.Now(),
			Name:  gn.author.Name,
			Email: gn.author.Email},
	})
	if err != nil {
		return err
	}

	return nil
}

// checkAndAddOrigin checks if an remote origin exists
// if it don't, it tries to add gn.RemoteURL as origin
// returns errflags.NoRemote if didn't find remote and couldn't create one
// or another error in case other operation fails
func (gn *GN) checkAndAddOrigin(r *git.Repository) error {
	if _, err := r.Remote("origin"); err != nil {
		if err != git.ErrRemoteNotFound {
			return err
		}

		// if no remote origin was found but there is a RemoteURL,
		// add it as origin and continue
		if gn.RemoteURL != "" {
			err = gn.AddOrigin(gn.RemoteURL)
			if err != nil {
				return err
			}
		} else {
			return errflags.Flag(err, errflags.NoRemote)
		}
	}
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

// ReadNote returns the content of the note
func (gn *GN) ReadNote() (string, error) {
	project, err := gn.findProject()
	if err != nil {
		return "", err
	}

	branch, err := gn.findBranch()
	if err != nil {
		return "", err
	}

	err = gn.createNotesPath()
	if err != nil {
		return "", err
	}

	notePath := getNotePath(gn.NotesPath, project, branch)

	if err := os.MkdirAll(filepath.Dir(notePath), os.ModeDir|0700); err != nil {
		return "", err
	}

	f, err := os.ReadFile(notePath)
	if err != nil {
		return "", err
	}

	return string(f), nil
}

func (gn *GN) Delete() error {
	var err error

	project, err := gn.findProject()
	if err != nil {
		return err
	}

	branch, err := gn.findBranch()
	if err != nil {
		return err
	}

	notePath := getNotePath(gn.NotesPath, project, branch)
	return os.Remove(notePath)
}

// getNotePath returns the path on the filesystem of the note
func getNotePath(notesPath string, project string, branch string) string {
	projectPath := fmt.Sprintf("%s/%s", notesPath, project)
	return fmt.Sprintf("%s/%s", projectPath, branch)
}
