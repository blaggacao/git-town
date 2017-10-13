# frozen_string_literal: true

Given(/^Git Town is in offline mode$/) do
  set_global_configuration 'offline', true
end


Given(/^Git Town is not in offline mode$/) do
  set_global_configuration 'offline', false
end


Given(/^I don't have a main branch name configured$/) do
  delete_main_branch_configuration
end


Given(/^my perennial branches are not configured$/) do
  delete_perennial_branches_configuration
end


Given(/^my repository has the "([^"]*)" configuration set to "([^"]*)"$/) do |configuration, value|
  set_configuration configuration, value
end


Given(/^I have configured the main branch name as "(.*)"$/) do |main_branch_name|
  set_configuration 'main-branch-name', main_branch_name
end


Given(/^my perennial branches are configured as (.*)$/) do |data|
  branch_names = Kappamaki.from_sentence data
  set_configuration 'perennial-branch-names', branch_names.join(' ')
end


Given(/I haven't configured Git Town yet/) do
  delete_main_branch_configuration
  delete_perennial_branches_configuration
end


Given(/^I configure "([^"]*)" as "([^"]*)"$/) do |key, value|
  run "git config '#{key}' '#{value}'"
end



Then(/^my repo is configured with perennial branches as "(.*)"$/) do |data|
  branch_names = Kappamaki.from_sentence(data)
  expect(perennial_branch_configuration.split(' ').map(&:strip)).to match_array branch_names
end


Then(/^my repo is configured with no perennial branches$/) do
  expect(perennial_branch_configuration).to be_empty
end


Then(/^my repo is configured with the main branch as "([^"]*)"$/) do |branch_name|
  expect(main_branch_configuration).to eql branch_name
end


Then(/^my repo is now configured with "([^"]*)" set to "(.+?)"$/) do |configuration, value|
  expect(get_configuration(configuration)).to eql value
end


Then(/^Git Town is (?:no longer|still not) configured for this repository$/) do
  expect(git_town_configuration).to be_empty
end


Then(/^I see the initial configuration prompt$/) do
  step %(I see "Git Town needs to be configured")
end


Then(/^offline mode is enabled$/) do
  expect(get_configuration('offline')).to eql 'true'
end


Then(/^offline mode is disabled$/) do
  expect(get_configuration('offline')).to eql 'false'
end
