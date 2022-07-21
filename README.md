# gitnotes :notebook:

gitnotes is a simple git-aware notes manager. It makes taking notes easy while working on multiple projects/branchs by opening the correct note when you run `gn edit`.

gitnotes will use your own text editor, and save notes as regular files in the default directory `$HOME/gitnotes`.

```bash
# example usage
$ cd my-project # on branch main of my-project-

$ gn edit # opens my-project/main note on editor

$ git checkout -b another-branch

$ gn edit # opens my-project/another-branch note

$ git checkout main

$ gn edit # back to notes of my-project/main

$ vim $(gn path) # opens directory of all notes
```

## Usage

The most common use is to just `gn edit`, take notes and save. Later, you may checkout another branch or work on another project, then you just `gn edit` again and take notes. Once you go back to the original project/branch, your notes are stored and you can pick up from where you left.

You can use the flags `-b` and `-p` to edit notes from a different branch and project, respectivelly. If you just want to read the notes in the terminal without opening an editor, you can use `gn print` instead of `gn edit`.

All your notes will be stored in `$HOME/gitnotes` (by default), making them easy to version. gitnotes comes with commands to help you version your own notes on git, like `gn pull`, `gn commit` and `gn push`.

If you try to run `gn edit` on a directory that is not a git repository without providing a project and branch, it will error.

Run `gn help` for more details.

```bash
usage: gn [-d] <command> <args>
Available commands:
- edit: edit the git note
- push: push notes to remote
- pull: pull notes from remote
- commit: commit notes
- path: prints the notes path to stdio
- print: prints the note to stdio
- delete: delete notes
run 'gn [command] -h' for more details on each command
```

## Instalation

Dependencies:

- Golang (building)
- git

1. Download source code

```bash
git clone https://github.com/mcbattirola/gitnotes.git
```

2. Build and install

```bash
make install
```

This will build the binary and move it to `/usr/local/bin/gn`. You can run `make build` and move `./dist/gn` to another directory if you prefer.

## Config file

gitnotes will create a config file if it doesn't find `$HOME/.config/gitnotes/gn.conf`.

### Default config file

```bash
editor=vim # binary name of the code editor
notes=$HOME/gitnotes # path in which notes will be stored
always-commit=false # commit after each `gn edit` (true/false)
```

## Troubleshooting

If you have problems running `gn push` to github, try running the following:

```bash
ssh-keyscan -t rsa github.com > ~/.ssh/known_hosts
ssh-keyscan -t ecdsa github.com >> ~/.ssh/known_hosts
```

See [this](https://github.com/go-git/go-git/issues/411) issue for more details.
