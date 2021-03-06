package steps

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v10"
	"github.com/git-town/git-town/src/command"
	"github.com/git-town/git-town/test"
)

// RunSteps defines Gherkin step implementations around running things in subshells.
// nolint:funlen
func RunSteps(suite *godog.Suite, state *ScenarioState) {
	suite.Step(`^I run "(.+)"$`, func(command string) error {
		state.runRes, state.runErr = state.gitEnv.DevShell.RunString(command)
		return nil
	})

	suite.Step(`^I run "([^"]+)" and answer the prompts:$`, func(cmd string, input *messages.PickleStepArgument_PickleTable) error {
		state.runRes, state.runErr = state.gitEnv.DevShell.RunStringWith(cmd, command.Options{Input: tableToInput(input)})
		return nil
	})

	suite.Step(`^I run "([^"]*)", answer the prompts, and close the next editor:$`, func(cmd string, input *messages.PickleStepArgument_PickleTable) error {
		env := append(os.Environ(), "GIT_EDITOR=true")
		state.runRes, state.runErr = state.gitEnv.DevShell.RunStringWith(cmd, command.Options{Env: env, Input: tableToInput(input)})
		return nil
	})

	suite.Step(`^I run "([^"]*)" and close the editor$`, func(cmd string) error {
		env := append(os.Environ(), "GIT_EDITOR=true")
		state.runRes, state.runErr = state.gitEnv.DevShell.RunStringWith(cmd, command.Options{Env: env})
		return nil
	})

	suite.Step(`^I run "([^"]*)" and enter an empty commit message$`, func(cmd string) error {
		state.runRes, state.runErr = state.gitEnv.DevShell.RunStringWith(cmd, command.Options{Input: []string{"dGZZ"}})
		return nil
	})

	suite.Step(`^I run "([^"]+)" in the "([^"]+)" folder$`, func(cmd, folderName string) error {
		state.runRes, state.runErr = state.gitEnv.DevShell.RunStringWith(cmd, command.Options{Dir: folderName})
		return nil
	})

	suite.Step(`^"([^"]*)" launches a new pull request with this url in my browser:$`, func(tool string, url *messages.PickleStepArgument_PickleDocString) error {
		want := fmt.Sprintf("%s called with: %s", tool, url.Content)
		want = strings.ReplaceAll(want, "?", `\?`)
		regex := regexp.MustCompile(want)
		have := state.runRes.OutputSanitized()
		if !regex.MatchString(have) {
			return fmt.Errorf("EXPECTED: a regex matching %q\nGOT: %q", want, have)
		}
		return nil
	})

	suite.Step(`^it runs no commands$`, func() error {
		commands := test.GitCommandsInGitTownOutput(state.runRes.Output())
		if len(commands) > 0 {
			for _, command := range commands {
				fmt.Println(command)
			}
			return fmt.Errorf("expected no commands but found %d commands", len(commands))
		}
		return nil
	})

	suite.Step(`^it runs the commands$`, func(input *messages.PickleStepArgument_PickleTable) error {
		commands := test.GitCommandsInGitTownOutput(state.runRes.Output())
		table := test.RenderExecutedGitCommands(commands, input)
		dataTable := test.FromGherkin(input)
		expanded := dataTable.Expand(
			state.gitEnv.DevRepo.Dir,
			&state.gitEnv.DevRepo,
			state.gitEnv.OriginRepo,
		)
		diff, errorCount := table.EqualDataTable(expanded)
		if errorCount != 0 {
			fmt.Printf("\nERROR! Found %d differences in the commands run\n\n", errorCount)
			fmt.Println(diff)
			return fmt.Errorf("mismatching commands run, see diff above")
		}
		return nil
	})
}

func tableToInput(table *messages.PickleStepArgument_PickleTable) []string {
	var result []string
	for i := 1; i < len(table.Rows); i++ {
		row := table.Rows[i]
		answer := row.Cells[1].Value
		answer = strings.ReplaceAll(answer, "[ENTER]", "\n")
		answer = strings.ReplaceAll(answer, "[DOWN]", "\x1b[B")
		answer = strings.ReplaceAll(answer, "[UP]", "\x1b[A")
		answer = strings.ReplaceAll(answer, "[SPACE]", " ")
		result = append(result, answer)
	}
	return result
}
