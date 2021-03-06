package steps

import (
	"fmt"

	"github.com/cucumber/godog"
)

// Merge defines Gherkin step implementations around merges.
func MergeSteps(suite *godog.Suite, state *ScenarioState) {
	suite.Step(`^my repo (?:now|still) has a merge in progress$`, func() error {
		hasMerge, err := state.gitEnv.DevRepo.HasMergeInProgress()
		if err != nil {
			return err
		}
		if !hasMerge {
			return fmt.Errorf("expected merge in progress")
		}
		return nil
	})

	suite.Step(`^there is no merge in progress$`, func() error {
		hasMerge, err := state.gitEnv.DevRepo.HasMergeInProgress()
		if err != nil {
			return err
		}
		if hasMerge {
			return fmt.Errorf("expected no merge in progress")
		}
		return nil
	})
}
