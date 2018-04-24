# Trello Plugin for Gitdo
### General notice
This is a plugin for [Gitdo](https://github.com/nebloc/Gitdo) built in [Python 3](https://www.python.org/downloads/release/latest). It is bundled with the Gitdo application, and should be considered an example implementation of a plugin; as Gitdo's plugin interaction changes, this will be the most up to date example on how it works.

The [interp file](https://github.com/nebloc/Trello-GitdoPlugin/blob/master/interp) tells Gitdo the command to run the files, if your Python 3 uses something different, e.g. `python`, change this file before running gitdo.

#### NEEDS GITDO VERSION HIGHER THAN V0.0.8 as changes were made to create, that switches which order arguments are given to the create function. For versions lower, please download the previous Gitdo release and use the plugin bundled there.

For more information on how Gitdo plugins can be developed or used, see the [Gitdo Wiki](https://github.com/nebloc/Gitdo/wiki/Plugins).

## Install
Windows:
```
git clone https://github.com/nebloc/Trello-GitdoPlugin.git %AppData%\Gitdo\plugins\Trello
```
Mac and Linux:
```
git clone https://github.com/nebloc/Trello-GitdoPlugin.git ~/.gitdo/plugins/Trello
```

## How it works
The plugin is comprised of 4 'functions':
* [`setup`](#setup)
* [`getid`](#get-id)
* [`create`](#create)
* [`done`](#done)

### [Setup](setup)
Asks the user for information that the plugin needs, in order to interact with Trello.

The idea of Gitdo is to configure early, and then get out of the way. So this function will save the settings of the project which are:
1. User's Trello Key - [Found here][trello]
1. User's Trello Token - [Found here][trello]
1. The list ID for new TODOs to go to - [Help](#list-id)
1. The list ID for done TODOs to go to - [Help](#list-id)
1. Remote link [Help](#remote-link)
---
##### List ID
The easiest way I have found to get the ID of a list on Trello, is:
1. Open a card on the list
1. Click "Share and more..." in the bottom right
1. Click "Export JSON"
1. Find the attribute idList

---
##### Remote link
This is an example of something that can be added to a card in Trello to provide the most information possible.

The link is to point to the exact location in the remote git repository of when the commit that contained the TODO happened, i.e. for Github, given a remote link such as,
```
# GITHUB
https://github.com/<username>/<your_project>/blob/{hash}/{file_name}#L{file_line}
# BITBUCKET
https://bitbucket.org/<username>/<your_project>/src/{hash}/{file_name}#{file_name}-{file_line}
```
the card added to Trello will have a permalink to the comment.

An example of this being used is the link https://github.com/nebloc/gitdo/blob/4892f877b299c00220c16f43ce377d1ca45b6a51/commit.go#L55 on [this](https://trello.com/c/G8F6PYby) Trello card.

### [Get ID](getid)
1. Converts the argument given (string of task in JSON format) to JSON.
1. Creates a task_name card with the description of the author, file_name, and file_line.
1. Card is added to the corresponding list (no. 3 in [Setup](#setup)).
1. Card is marked as `closed` so that it is added to the Archive. This is to hide the task as it relates to code that cannot be viewed by the team yet.
1. Gets the short ID of the new card and prints it to Stdout for Gitdo to use.

### [Create](create)
1. Process the ID and task JSON given.
1. Creates a new description for the card that adds the Git hash of the added line, and the branch it was on. If provided in setup, the remote link is also added with the hash, file_name, and file_line, as a link.
1. The card has `closed` set to false so it is taken out the Archive and moved to the new list.

The Card is now on the list, and has all the information it should need, for a team to make a decision on it.

No actions on the card will effect the source code. Ideally it would be moved to in progress for a sprint, and can be discussed, have more information like implementation choices added to it, and then eventually handled by the team.



### [Done](done)
Is given an ID as an argument, that is used to find the card, which is then simply moved to the list specified in [Setup #4](#setup).

[trello]: https://trello.com/app-key
