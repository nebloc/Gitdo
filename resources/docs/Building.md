1. [Install the Go programming language](https://golang.org/doc/install)
1. Run `go get -u golang.org/x/vgo` in a terminal
1. Run `go get -u github.com/nebloc/app/...`
1. Run `vgo install` to install gitdo to your $GOPATH/bin directory.

### Installing Plugins and Hooks
As we have installed Gitdo to `$GOPATH/bin`, which should be a part of your `$PATH` as outlined in the [How to Write Go Code docs,](https://golang.org/doc/code.html#GOPATH) we just need to copy the Plugins and Hooks.

Gitdo has it's home folder in either `%AppData%\Gitdo` or `~/.gitdo` (For this doc, replace `$GITDO` with the correct one).
The folder should look like:
```
$GITDO/
  plugins/
    List of plugins
  hooks/
    pre-commit
    post-Commit
    push
```
1. Create `$GITDO` directory, i.e. `mkdir ~/.gitdo`
1. Copy the hooks folder into the directory, i.e. `cp -r ./hooks ~/.gitdo`
1. Copy the plugins folder into the directory, i.e. `cp -r ./plugins ~/.gitdo`


##### Version numbers
*Optionally, running this command will change the version number to `0.0.0-A5-Hash` where hash is the current commit hash*
```
vgo install -ldflags "-X main.version=0.0.0-A5-$(git rev-parse --short HEAD)"
```


### Testing
Run unit tests using the command `vgo test github.com/nebloc/gitdo/app/...`
The unit tests need a lot of work at the moment, and do use the temporary directory to create folders with git and hg repositories.
