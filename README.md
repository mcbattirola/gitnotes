# gitnotes

Commands:

```bash
gn edit # edit current branch notes; alternativaly accepts a branch name
gn delete # deletes a note
gn commit # commits notes to remote
gn push # push notes to remote
gn pull # pull notes from remote
gn path # prints notes path into stdout
```

---

## TODO

- add log and log level
- improve tests
- refactor gn struct
- setup ci pipeline
- make a real README

## Ideas

### edit -m option

`gn edit -m "message"` appends `message` to the end of note

### Note template / headers

- add a header to each new note (notes on branch xxx)
- header=true/false
- header-template=path to header template file
