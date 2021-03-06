Feature: git town-ship: trying the ship of the current feature branch without editing the default commit message

  As a developer shipping a branch
  I want the ship to abort if I don't edit the default commit message
  So that I don't have ugly commits merged into my main branch


  Background:
    Given my repository has a feature branch named "feature"
    And the following commits exist in my repository
      | BRANCH  | LOCATION | MESSAGE        | FILE NAME    | FILE CONTENT    |
      | feature | local    | feature commit | feature_file | feature content |
    And I am on the "feature" branch
    When I run "git-town ship" and close the editor


  Scenario: result
    Then it runs the commands
      | BRANCH  | COMMAND                            |
      | feature | git fetch --prune --tags           |
      |         | git checkout main                  |
      | main    | git rebase origin/main             |
      |         | git checkout feature               |
      | feature | git merge --no-edit origin/feature |
      |         | git merge --no-edit main           |
      |         | git checkout main                  |
      | main    | git merge --squash feature         |
      |         | git commit                         |
      |         | git reset --hard                   |
      |         | git checkout feature               |
      | feature | git checkout main                  |
      | main    | git checkout feature               |
    And it prints the error:
      """
      Aborted because commit exited with error
      """
    And I am still on the "feature" branch
    And my repository is left with my original commits


  Scenario: undo
    When I run "git-town undo"
    Then it prints the error:
      """
      Nothing to undo
      """
    And I am still on the "feature" branch
    And my repository is left with my original commits
