Gitdo - Track TODO comments in your task manager.
=================

![AppLogo](https://github.com/nebloc/gitdo/wiki/images/GitdoLogo.png)

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

### [Install](https://github.com/nebloc/Gitdo/wiki/Install)

### [Usage](https://github.com/nebloc/Gitdo/wiki/Usage)

### [Building](https://github.com/nebloc/Gitdo/wiki/Building)

# Architecture
![Architecture design](https://github.com/nebloc/gitdo/wiki/images/Architecture.png)

## Sequence Diagrams
[Git Commit Sequence Diagram](https://github.com/nebloc/gitdo/wiki/images/sequence_diagram_commit.png)

[Git Push Sequence Diagram](https://github.com/nebloc/gitdo/wiki/images/sequence_diagram_push.png)

# Known Issues
1. If a commit message is empty, gitdo will run, even though git will fail.
1. If a task has a new line character it will not look at the second line, and treat only the first line.
1. If a plugin fails, it is not clear whether gitdo will fail to, and stop the commit.
1. Cannot get it to run when committing from Eclipse. Running Git from CLI is best
1. Intellij runs hook if selected, but will not give information unless it fails. Running Git from CLI is best.
