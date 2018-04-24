package versioncontrol

import "errors"

// VersionControl is an interface for defining the functions that are needed to interact with version control systems on ahost machine
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
	PathOfTopLevel() string

	// Add changed tasks back to staging
	RestageTasks(fileName string) error

	// Set the hooks that are needed for the VC during init
	SetHooks(homeDir string) error
}

var (
	// ErrNotVCDir is thrown when the version control system that is trying to be used has not been initialised.
	ErrNotVCDir = errors.New("directory is not a git or mercurial repo")
	// ErrNoDiff is thrown when the current repository is clean.
	ErrNoDiff = errors.New("diff is empty")
)
