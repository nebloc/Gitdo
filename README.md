Gitdo - Track TODO comments in your task manager.
=================

![App Icon](./docs/images/Logo.png)

1. [Introduction](#introduction)
1. [Install](#install)
1. [Usage](#usage)
1. [Architecture](#architecture)
1. [Issues](#known-issues)

# Introduction
#### Using experimental vgo tool for dependencies.
install: `go get -u golang.org/x/vgo`
[See research by Russ Cox here](https://research.swtch.com/vgo)

## Mission
Suboptimal code is not correctly tracked and corrected in modern software methodologies, often leading to quality or performance issues later in the application lifecycle.

When pushed for time or pressured from other priorities, developers may annotate unfinished code, with certain terms, in order to show that the following code is not complete. However, these annotations are often not visible enough and go long periods without being fixed.

This project will aim to aid development in this problem domain by providing an open source, extendable tool, that uses common practices and technologies to track unclean code. This will include analyzing code for comments with keywords, at different times in the version control cycle, to garner as much information as possible, and provide a bridge between the offending code and a task management system, mainly focusing on Kanban board services, such as Trello. This process will be automated using Git hooks to start analysis and extraction of annotations. The extendable and open source nature will allow further development of other services, such as personal Todo apps, and corporate ticket systems, i.e. Jira.

The hope is that a dedicated tool will push teams to formalise their use of TODO annotations and allow them to be tracked more efficiently.

# Install
1. [Download](https://github.com/nebloc/Gitdo/releases) and unzip the latest release
1. Navigate to the folder in your console and run install.sh (*nix) or install.bat (Windows).
1. This will place the Gitdo executable in your usr/local/bin or system32. As well as add the hooks and plugins to a home directory: `%AppData%\roaming\Gitdo` on windows and `~/.gitdo` on *nix
This is where plugins are found, so to create and use your own, you should create them inside here. More information will be in the Wiki soon or you can look at the example Trello plugin.

# Usage
1. At the moment, only todo comments that are in the form `// TODO: Something`, so if you are working on a TODO that you do not want added, use another keyword like `//HACK:` or `//!TODO:`. Work is being done on adding other comment styles such as pythons `#`
1. When you commit your work, a git hook (pre-commit) should call the Gitdo tool `gitdo commit -c`. Gitdo will analyse a `git diff` for task changes and mark the source code with a tag: 
`//TODO: Make something cool <uhvc302n>`
. The `<id>` Will be added to the end of a task comment in the source code, before being restaged for commit. 
1. When `git push` is ran, a pre-push script will activate Gitdo and it will parse all tasks that have been committed since last push to the configured plugin. Example plugin uses Trello API to add tasks to a specified list.
1. IDs of tasks that have been removed (as shown in git diff) that were tagged, will be parsed to the plugin to be marked as done in the task manager.

# Architecture
![Architecture design](./docs/images/Architecture.png)

## Sequence Diagram
[Sequence Diagram](./docs/images/SequenceDiagram.png)

# Known Issues
1. If a commit message is empty, gitdo will run, even though git will fail.
1. If a task has a new line character it will not look at the second line, and treat only the first line.
1. If a plugin fails, it is not clear whether gitdo will fail to, and stop the commit.
1. Cannot get it to run when committing from Eclipse. Running Git from CLI is best
1. Intellij runs hook if selected, but will not give information unless it fails. Running Git from CLI is best.
