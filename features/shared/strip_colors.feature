Feature: Strip colors

  As a developer with Git configured so colors are always on
  I expect the tool to still work
  So that Git Town does not place addititional constraints on my computer setup


  Scenario: colors are stripped from the output of git commands run internally
    Given I haven't configured Git Town yet
    And my repo has "color.ui" set to "always"
    And I am on the "main" branch
    When I run "git-town hack new-feature" and answer the prompts:
      | PROMPT                                     | ANSWER  |
      | Please specify the main development branch | [ENTER] |
      | Please specify perennial branches          | [ENTER] |
    And Git Town is now aware of this branch hierarchy
      | BRANCH      | PARENT |
      | new-feature | main   |
