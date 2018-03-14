# Gitdo

## NOT yet ready for pull requests.
## Using experimental vgo tool for dependencies.
install: `go get -u golang/x/vgo`
[See research by Russ Cox here](http://https://research.swtch.com/vgo)

Suboptimal code is not correctly tracked and corrected in modern software methodologies, often leading to quality or performance issues later in the application lifecycle.

When pushed for time or pressured from other priorities, developers may annotate unfinished code, with certain terms, in order to show that the following code is not complete. However, these annotations are often not visible enough and go long periods without being fixed.

This project will aim to aid development in this problem domain by providing an open source, extendable tool, that uses common practices and technologies to track unclean code. This will include analyzing code for comments with keywords, at different times in the version control cycle, to garner as much information as possible, and provide a bridge between the offending code and a task management system, mainly focusing on Kanban board services, such as Trello. This process will be automated using Git hooks to start analysis and extraction of annotations. The extendable and open source nature will allow further development of other services, such as personal Todo apps, and corporate ticket systems, i.e. Jira.

The hope is that a dedicated tool will push teams to formalise their use of TODO annotations and allow them to be tracked more efficiently.

