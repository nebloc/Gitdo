package versioncontrol

import "errors"

// VersionControl is the interface for different version control systems
type VersionControl interface {
	// Initialise the VC system on init
	Init() error

	// Get information from the version control
	GetDiff() (string, error)
	GetEmail() (string, error)
	GetBranch() (string, error)
	GetHash() (string, error)
	GetTrackedFiles(branch string) ([]string, error)

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

	NewCommit(message string) error
	CheckClean() bool
}

var (
	// ErrNotVCDir is thrown when the current directory is not inside a repository
	ErrNotVCDir = errors.New("directory is not a git or mercurial repo")
	// ErrNoDiff is thrown when the diff output is empty
	ErrNoDiff = errors.New("diff is empty")
)

// NewBranchName is the name of the branch that force-all does it's tagging on
const NewBranchName = "gitdo/taggingall"
