package versioncontrol

import "errors"

type VersionControl interface {
	// Initialise the VC system on init
	Init() error

	// Get information from the version control
	GetDiff() (string, error)
	GetEmail() (string, error)
	GetBranch() (string, error)
	GetHash() (string, error)
	GetTrackedFiles() ([]string, error)

	// Get details of the version control being used
	NameOfDir() string
	NameOfVC() string
	PathOfTopLevel() string

	// Add changed tasks back to staging
	RestageTasks(fileName string) error
	CreateBranch() error
	SwitchBranch() error

	// Set the hooks that are needed for the VC during init
	SetHooks(homeDir string) error
}

var (
	ErrNotVCDir = errors.New("directory is not a git or mercurial repo")
	ErrNoDiff   = errors.New("diff is empty")
)

const NewBranchName = "gitdo/taggingall"
