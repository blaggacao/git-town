package steps

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sync"

	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v10"
	"github.com/git-town/git-town/test"
	"github.com/git-town/git-town/test/helpers"
)

// beforeSuiteMux ensures that we run BeforeSuite only once globally.
var beforeSuiteMux sync.Mutex

// the global GitManager instance
var gitManager *test.GitManager

var running helpers.OrderedStringSet
var runningMux sync.Mutex

// SuiteSteps defines global lifecycle step implementations for Cucumber.
func SuiteSteps(suite *godog.Suite, state *ScenarioState) {
	suite.BeforeScenario(func(scenario *messages.Pickle) {
		runningMux.Lock()
		running = running.Add(scenario.Name)
		fmt.Printf("\nStarting scenario %q, all scenarios: %s", scenario.Name, running.String())
		runningMux.Unlock()
		// create a GitEnvironment for the scenario
		gitEnvironment, err := gitManager.CreateScenarioEnvironment(scenario.GetName())
		if err != nil {
			log.Fatalf("cannot create environment for scenario %q: %s", scenario.GetName(), err)
		}
		// Godog only provides state for the entire feature.
		// We want state to be scenario-specific, hence we reset the shared state before each scenario.
		// This is a limitation of the current Godog implementation, which doesn't have a `ScenarioContext` method,
		// only a `FeatureContext` method.
		// See main_test.go for additional details.
		state.Reset(gitEnvironment)
		if hasTag(scenario, "@debug") {
			test.Debug = true
		}
	})

	suite.BeforeSuite(func() {
		// NOTE: we want to create only one global GitManager instance with one global memoized environment.
		beforeSuiteMux.Lock()
		defer beforeSuiteMux.Unlock()

		running = helpers.NewOrderedStringSet()

		if gitManager == nil {
			baseDir, err := ioutil.TempDir("", "")
			if err != nil {
				log.Fatalf("cannot create base directory for feature specs: %s", err)
			}
			// Evaluate symlinks as Mac temp dir is symlinked
			evalBaseDir, err := filepath.EvalSymlinks(baseDir)
			if err != nil {
				log.Fatalf("cannot evaluate symlinks of base directory for feature specs: %s", err)
			}
			gitManager = test.NewGitManager(evalBaseDir)
			err = gitManager.CreateMemoizedEnvironment()
			if err != nil {
				log.Fatalf("Cannot create memoized environment: %s", err)
			}
		}
	})

	suite.AfterScenario(func(scenario *messages.Pickle, e error) {
		runningMux.Lock()
		running = running.Remove(scenario.Name)
		fmt.Printf("\nFinished scenario %q, all scenarios: %s", scenario.Name, running.String())
		runningMux.Unlock()
		if e == nil {
			err := state.gitEnv.Remove()
			if err != nil {
				log.Fatalf("error removing the Git environment after scenario %q: %v", scenario.GetName(), err)
			}
		} else {
			fmt.Printf("failed scenario, investigate state in %q\n", state.gitEnv.Dir)
		}
	})
}

// hasTag indicates whether the given feature has a tag with the given name.
func hasTag(scenario *messages.Pickle, name string) bool {
	for _, tag := range scenario.GetTags() {
		if tag.Name == name {
			return true
		}
	}
	return false
}
