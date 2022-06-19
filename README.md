# gitnotes

Commands:

```bash
gn edit # edit current branch notes; alternativaly accepts a branch name
gn commit # (not implemented) commits notes to remote
gn push # (not implemented) push notes to remote
gn pull # (not implemented) pull notes from remote
gn sync # (not implemented) pull and push
gn path # prints notes path into stdout
```

- `.config` file to configurations later
- works locally, can be synced

---

## TODO

- add log and log level
- add a header to each new note (notes on branch xxx)
- consider replacing git lib for git syscalls
- implement missing commands
- make a real README
- setup ci pipeline

### Config fields

Config fields pending implementation:

- always-commit=true/false
- always-push=true/false
- header=true/false
- header-template=path to header template file
