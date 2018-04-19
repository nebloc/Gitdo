package main

type VersionControl interface {
	// Initialise the VC system on init
	Init() error

	// Get information from the version control
	GetDiff() (string, error)
	GetEmail() (string, error)
	GetBranch() (string, error)
	GetHash() (string, error)

	// Get details of the version control being used
	NameOfDir() string
	NameOfVC() string

	// Get the top level directory for the project
	SetTopLevel(topLevel string)
	GetTopLevel() string

	// Add changed tasks back to staging
	RestageTasks(task Task) error

	// Set the hooks that are needed for the VC during init
	SetHooks(homeDir string) error
}
