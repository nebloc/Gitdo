Gitdo - Track TODO comments in your task manager.
=================

![AppLogo](https://github.com/nebloc/gitdo/wiki/images/GitdoLogo.png)

# Introduction
Gitdo is a tool for formalising the tracking of task annotations in source code. The idea is that when a new git (working on mercurial support as well) commit contains a new task such as;
```
// TODO: Make sure the README makes sense in explaining this.
```
Gitdo will capture it at pre-commit time, use the chosen plugin ([see here](https://github.com/nebloc/gitdo/wiki/plugins)) to get an ID for which will be added to the source code. Then when a git push is ran it will add the new tasks to the chosen task manager.

### Example video
[![See example](https://img.youtube.com/vi/czmJEh818Qo/0.jpg)](https://youtu.be/czmJEh818Qo)

### [See Wiki for documentation](https://github.com/nebloc/gitdo/wiki)

### Task Annotations
Example|Captured
-------|--------
`//TODO: Is this captured?`| Y
`// TODO: Is this captured?`| Y
`// TODO Is this captured?` | N
`# TODO: Is this captured?` | Y
`#TODO: Is this captured?` | Y
`# TODO Is this captured?` | N
`//TODO(): Is this captured`|Y

 [Test your own string](https://play.golang.org/p/PVfowMCOkyJ)

#### Using experimental vgo tool for dependencies.
install: `go get -u golang.org/x/vgo`
[See research by Russ Cox here](https://research.swtch.com/vgo)

# Known Issues
1. If a commit message is empty, gitdo will run, even though git will fail.
1. If a task has a new line character it will not look at the second line, and treat only the first line.
1. If a plugin fails, it is not clear whether gitdo will fail to, and stop the commit.
1. ~~Cannot get it to run when committing from Eclipse. Running Git from CLI is best~~  [How to fix](https://github.com/nebloc/Gitdo/wiki/Usage#eclipse).
1. Intellij runs hook if selected, but will not give information unless it fails. Running Git from CLI is best.
