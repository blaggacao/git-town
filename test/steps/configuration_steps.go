package steps

import (
	"fmt"
	"strconv"

	"github.com/cucumber/godog"
)

// ConfigurationSteps defines Cucumber step implementations around configuration.
// nolint:funlen,gocognit
func ConfigurationSteps(suite *godog.Suite, state *ScenarioState) {
	suite.Step(`^Git Town is no longer configured for this repository$`, func() error {
		res, err := state.gitEnv.DevRepo.HasGitTownConfigNow()
		if err != nil {
			return err
		}
		if res {
			return fmt.Errorf("unexpected Git Town configuration")
		}
		return nil
	})

	suite.Step(`^I haven't configured Git Town yet$`, func() error {
		state.gitEnv.DevRepo.Configuration(false).DeleteMainBranchConfiguration()
		state.gitEnv.DevRepo.Configuration(false).DeletePerennialBranchConfiguration()
		return nil
	})

	suite.Step(`^my repo has "color\.ui" set to "([^"]*)"$`, func(value string) error {
		_ = state.gitEnv.DevRepo.Configuration(false).SetColorUI(value)
		return nil
	})

	suite.Step(`^my repo has "git-town.sync-upstream" set to (true|false)$`, func(text string) error {
		value, err := strconv.ParseBool(text)
		if err != nil {
			return err
		}
		_ = state.gitEnv.DevRepo.Configuration(false).SetShouldSyncUpstream(value)
		return nil
	})

	suite.Step(`^my repo has "git-town.code-hosting-driver" set to "([^"]*)"$`, func(value string) error {
		_ = state.gitEnv.DevRepo.Configuration(false).SetCodeHostingDriver(value)
		return nil
	})

	suite.Step(`^my repo has "git-town.code-hosting-origin-hostname" set to "([^"]*)"$`, func(value string) error {
		_ = state.gitEnv.DevRepo.Configuration(false).SetCodeHostingOriginHostname(value)
		return nil
	})

	suite.Step(`^my repo has "git-town.ship-delete-remote-branch" set to "(true|false)"$`, func(value string) error {
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		_ = state.gitEnv.DevRepo.Configuration(false).SetShouldShipDeleteRemoteBranch(parsed)
		return nil
	})

	suite.Step(`^my repo is now configured with no perennial branches$`, func() error {
		branches := state.gitEnv.DevRepo.Configuration(true).GetPerennialBranches()
		if len(branches) > 0 {
			return fmt.Errorf("expected no perennial branches, got %q", branches)
		}
		return nil
	})

	suite.Step(`^the new-branch-push-flag configuration is now (true|false)$`, func(text string) error {
		want, err := strconv.ParseBool(text)
		if err != nil {
			return err
		}
		have := state.gitEnv.DevRepo.Configuration(true).ShouldNewBranchPush()
		if have != want {
			return fmt.Errorf("expected global new-branch-push-flag to be %t, but was %t", want, have)
		}
		return nil
	})

	suite.Step(`^the global new-branch-push-flag configuration is (true|false)$`, func(text string) error {
		b, err := strconv.ParseBool(text)
		if err != nil {
			return err
		}
		_ = state.gitEnv.DevRepo.Configuration(false).SetNewBranchPush(b, true)
		return nil
	})

	suite.Step(`^the new-branch-push-flag configuration is (true|false)$`, func(value string) error {
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		state.gitEnv.DevRepo.Configuration(false).SetNewBranchPush(b, false)
		return nil
	})

	suite.Step(`^the main branch is configured as "([^"]+)"$`, func(name string) error {
		state.gitEnv.DevRepo.Configuration(false).SetMainBranch(name)
		return nil
	})

	suite.Step(`^the main branch is now configured as "([^"]+)"$`, func(name string) error {
		actual := state.gitEnv.DevRepo.Configuration(true).GetMainBranch()
		if actual != name {
			return fmt.Errorf("expected %q, got %q", name, actual)
		}
		return nil
	})

	suite.Step(`^the main branch name is not configured$`, func() error {
		state.gitEnv.DevRepo.Configuration(false).DeleteMainBranchConfiguration()
		return nil
	})

	suite.Step(`^the perennial branches are not configured$`, func() error {
		state.gitEnv.DevRepo.Configuration(false).DeletePerennialBranchConfiguration()
		return nil
	})

	suite.Step(`^the perennial branches are configured as "([^"]+)"$`, func(name string) error {
		state.gitEnv.DevRepo.Configuration(false).AddToPerennialBranches(name)
		return nil
	})

	suite.Step(`^the perennial branches are configured as "([^"]+)" and "([^"]+)"$`, func(branch1, branch2 string) error {
		state.gitEnv.DevRepo.Configuration(false).AddToPerennialBranches(branch1, branch2)
		return nil
	})

	suite.Step(`^the perennial branches are now configured as "([^"]+)"$`, func(name string) error {
		actual := state.gitEnv.DevRepo.Configuration(true).GetPerennialBranches()
		if len(actual) != 1 {
			return fmt.Errorf("expected 1 perennial branch, got %q", actual)
		}
		if actual[0] != name {
			return fmt.Errorf("expected %q, got %q", name, actual[0])
		}
		return nil
	})

	suite.Step(`^the perennial branches are now configured as "([^"]+)" and "([^"]+)"$`, func(branch1, branch2 string) error {
		actual := state.gitEnv.DevRepo.Configuration(true).GetPerennialBranches()
		if len(actual) != 2 {
			return fmt.Errorf("expected 2 perennial branches, got %q", actual)
		}
		if actual[0] != branch1 || actual[1] != branch2 {
			return fmt.Errorf("expected %q, got %q", []string{branch1, branch2}, actual)
		}
		return nil
	})

	suite.Step(`^the pull-branch-strategy configuration is "(merge|rebase)"$`, func(value string) error {
		state.gitEnv.DevRepo.Configuration(false).SetPullBranchStrategy(value)
		return nil
	})

	suite.Step(`^the pull-branch-strategy configuration is now "(merge|rebase)"$`, func(want string) error {
		have := state.gitEnv.DevRepo.Configuration(false).GetPullBranchStrategy()
		if have != want {
			return fmt.Errorf("expected pull-branch-strategy to be %q but was %q", want, have)
		}
		return nil
	})

	suite.Step(`^my repo is now configured with no perennial branches$`, func() error {
		actual := state.gitEnv.DevRepo.Configuration(true).GetPerennialBranches()
		if len(actual) > 0 {
			return fmt.Errorf("expected no perennial branches, got %q", actual)
		}
		return nil
	})
}
