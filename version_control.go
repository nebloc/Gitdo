package main

type VersionControl interface {
	Init() error
	GetDiff() (string, error)
	GetEmail() (string, error)

	GetBranch() (string, error)
	GetHash() (string, error)

	NameOfDir() string
	NameOfVC() string

	SetTopLevel(topLevel string)
	GetTopLevel() string

	RestageTasks(task Task) error

	SetHooks(homeDir string) error
}
