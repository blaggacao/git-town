package steps

import (
	"github.com/git-town/git-town/src/script"
)

// FetchUpstreamStep brings the Git history of the local repository
// up to speed with activities that happened in the upstream remote.
type FetchUpstreamStep struct {
	NoOpStep
	BranchName string
}

// Run executes this step.
func (step *FetchUpstreamStep) Run() error {
	return script.RunCommand("git", "fetch", "upstream", step.BranchName)
}
