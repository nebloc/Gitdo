Plugins are scripts that will be ran by Gitdo with different arguments, at different stages.
The current required functions of a plugin are:
* [setup](#setup)
* [getid](#get-id)
* [create](#create)
* [done](#done)

*For this document, Gitdo's home folder, `%AppData%\Gitdo\` for windows and  `~/.gitdo` for macOS and Linux, will be noted as `$GITDO`.*

*An example plugin that is bundled with Gitdo can be [found here.](https://github.com/nebloc/Trello-GitdoPlugin)*


They can be written in any language, such as Python 2/3, bash, Javascript, etc., and are currently invoked by calling an interpreter on the local machine. Therefore, if if you build a plugin in Python 3, the user's system must have a Python 3 interpreter, i.e.  running the command `python -V` or `python3 -V` results in an output of version 3.x.


# Setup
This function is ran when `gitdo init` is called. Gitdo will add a directory for the plugin to store project specific files and configuration,  `.git/gitdo/plugins/plugin-name`, and then run `<interpreter> `.
[Example function.](https://github.com/nebloc/Trello-GitdoPlugin/blob/master/setup)

Gitdo will **always** run functions with this directory as the current path. This is so that a plugin can save and load configuration files specific to each project.

The setup command is given Gitdo's Stdin, Stderr and Stdout and so can ask for user input.

The function should ask the user for any information it will need to complete the other functions. It is also a good idea to use this function to verify that the interpreter is correct for your plugin.

# Get ID
This function is ran as a `git commit` is being processed. As Gitdo finds comments that are the correct format to be a task, it passes a string of a json task as the first argument.
```
{
  "task_name":  string,
  "file_name":  string,
  "file_line":  number,
  "author":     string,
  "hash":       nil,
  "branch":     nil,
}
```

The goal of this function is to provide a unique ID for Gitdo to add back in to the code, as well as pass to the other plugin functions to act upon. [Example function.](https://github.com/nebloc/Trello-GitdoPlugin/blob/master/getid)

The argument can then be converted to JSON in the script (i.e. `json.loads(sys.argv[1])` in Python 3), to be used to give more information to the task manager. This is because in Trello, there is no way of reserving an ID, and so instead, a task is added to Trello as closed (archived), and if a user finds the card, we can give them a little bit of help as to why it is there with the `task_name`, `author`, `file_line` and `file_name`.

This function should **only** print to Stdout the id of the card, unless it is exiting with a non 0 status. If non 0 exit status, the output from the function is given to the user, and a warning that the comment could not be added. It will **not** stop others from being passed to the function or the commit itself.

# Create
This function is ran as a `git push` is being run by the user. It will be run once for each new task in the task file, and is passed a string in JSON format of task details, and the ID of the new task, i.e.
```
{
  "task_name":  string,
  "file_name":  string,
  "file_line":  number,
  "author":     string,
  "hash":       string,
  "branch":     string,
} (required)

nDbVH2fg (optional)
```
##### BREAKING CHANGE in v0.0.8 - This is the new order, before version 0.0.8, the order was ID (which was required) and then the task string. The order was changed to lower the number of plugin calls needed when processing an existing code base (`force-all` command not implemented yet).

This can then be used to add a new / update the task with the given ID, with information on where the comment is, and a commit hash and branch, so that it can be found in the remote repository. An example of this can be found [in the Trello plugin](https://github.com/nebloc/Trello-GitdoPlugin/blob/master/create), where a user can enter a remote link such as [https://github.com/name/project/blob/{HASH}/{file_name}#L{file_line}](https://github.com/nebloc/gitdo/blob/4892f877b299c00220c16f43ce377d1ca45b6a51/commit.go#L55).
It needs to return an ID of a new task.

# Done
This is for marking the task as done, when a commit that has removed the task completely, has been pushed to the remote repository. It will be ran at `git push` time, and will be ran once for each ID in the task file that is marked as Done. It is passed a single argument of the tasks ID (i.e. `nDbVH2fg`).

In the Trello example, the function will move the card to the list ID given by the user - probably a 'Done' list, or 'QA'/'Code Review'.

Some more simple task managers, like for the GTD app [Omnifocus](https://www.omnigroup.com/omnifocus), the plugin may opt to just mark it off as [complete](https://github.com/nebloc/Omnifocus-GitdoPlugin/blob/2933ef272d8400a46d264aac03675239830dd973/done#L16).
