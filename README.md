# gitnotes

Commands:

```bash
gn edit # edit current branch notes; alternativaly accepts a branch name
gn delete # (not implemented) deletes a note
gn commit # commits notes to remote
gn push # push notes to remote
gn pull # (not implemented) pull notes from remote
gn path # prints notes path into stdout
```

---

## TODO

- add log and log level
- consider replacing git lib for git syscalls
- implement missing commands
- setup ci pipeline
- make a real README

### Config fields

Config fields pending implementation:

- always-push=true/false

## Ideas

### Note template / headers

- add a header to each new note (notes on branch xxx)
- header=true/false
- header-template=path to header template file
