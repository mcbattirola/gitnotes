package gn

import (
	"os/exec"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/mcbattirola/gitnotes/pkg/errflags"
)

type Author struct {
	Name  string
	Email string
}

func getCurrentBranch(r *git.Repository) (string, error) {
	h, err := r.Reference(plumbing.HEAD, false)
	if err != nil {
		return "", err
	}

	target := h.Target()

	s := strings.Split(string(target), "refs/heads/")
	if len(s) < 2 {
		return "", errflags.New("couldn't find project branch", errflags.BadParameter)
	}
	return s[1], nil
}

func readGlobalGitAuthor() Author {
	// run git config --list
	// then find user name and email from it
	// TODO
	return Author{}
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
