# Trello Plugin for Gitdo
### General notice
This is a plugin for [Gitdo](https://github.com/nebloc/Gitdo) built in [Python 3](https://www.python.org/downloads/release/latest). It is bundled with the Gitdo application, and should be considered an example implementation of a plugin; as Gitdo's plugin interaction changes, this will be the most up to date example on how it works.

For more information on how Gitdo plugins can be developed or used, see the [Gitdo Wiki](https://github.com/nebloc/Gitdo/wiki/Plugins).

# How it works
The plugin is comprised of 4 'functions':
* [`setup`](#setup)
* [`getid`](#getid)
* [`create`](#create)
* [`done`](#done)

## [Setup](setup)
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

The link is to point to the exact location in the remote git repository of when the commit that contained the TODO happened, i.e. for Github, given the remote link  
`https://github.com/nebloc/Trello-GitdoPlugin/blob/{hash}/{file_name}#L{file_line}`,
the card added to Trello will have a permalink to the comment.

An example of this being used is the link https://github.com/nebloc/gitdo/blob/4892f877b299c00220c16f43ce377d1ca45b6a51/commit.go#L55 on [this](https://trello.com/c/G8F6PYby) Trello card.
[trello]: https://trello.com/app-key)

## [Get ID](getid)
## [Create](create)
## [Done](done)
