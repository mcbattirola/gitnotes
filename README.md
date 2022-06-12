# gitnotes

```bash
gn init # initializes gitnote in the current dir
gn edit # edit current branch notes; alternativaly accepts a branch name
gn sync # commits local notes to remote
gn clone # clones a remote notes repo
```

- .config file to configurations later
- works locally, can be synced

---

## TODO

- tests
- error handling
- gn init (git init)
- add a header to each new note (notes on branch xxx)
- handle actual versioning of the notes repository
- consider replacing git lib for git syscalls

### Config fields

Ideas for new fields to config file:

- always-commit=true/false
- header=true/false (no customization)
