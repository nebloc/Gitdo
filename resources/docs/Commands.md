There are a few commands that are meant to be ran by developers, and not as part of a git hook:
* list
* list --config
* init
* init --with-git (currently broken)
* destroy

### List
This pretty prints the current tasks that are new, and have been found in a commit but not added to the task manager in a Push yet, as well as the ID's of tasks that have been completed, removed from the source code in a commit.
i.e.
```
$ gitdo list

===New Tasks===
main.go#9:   Finish the documentation of Gitdo's List function      id#56qn0ORJ  
main.go#10:  Create a mock output of List to show in documentation  id#oQcrxtdl  
===Completed Tasks===
Done: nDbVH2fg
Done: z33eOILg
Done: Wf5SEuPO
```

### List Config
This is a flag for the list command that instead just prints the current contents of the configuration file. This allows users to quickly check which plugin they are using on a project, etc.
```
$ gitdo list --config

Author: email@example.com
Plugin: Trello
Interpreter: python3
```

### Init
For adding Gitdo to a new project, it creates the configuration file, copies the git hooks and starts the plugins setup function.
```
$ gitdo init

Copying from: /Users/bencoleman/.gitdo/hooks to .git/hooks
Using email@example.com
Available plugins:
1: Omnifocus
2: Test
3: Trello
What plugin would you like to use (1-3): 2
Using Test
Currently all plugins made as an example need python 3 set up in path. Redesign of plugin language choice and use coming soon.
What interpreter for this plugin (i.e. python3/node/python): python3
Using python3
No setup required.
Done
```

### Init --with-git (currently broken)
This is meant to be a flag that makes Gitdo run `git init` on the command line first, however, currently gitdo is checking that the directory is a git directory, and failing when it is not, before the flag is checked.

### Destroy
simply deletes the tasks file, meaning none of the new tasks or done tasks will be processed in the next `git push`.

Recommend against doing this, and the command is mostly there for testing purposes.
