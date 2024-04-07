# Github Summary

CLI tool to summarize activity of a repository/organization/user over the last X days.

## Usage

```bash
export DAYS=
export GITHUB_ORGANIZATION=
export GITHUB_PERSONAL_ACCESS_TOKEN=

go run main.go
```

## CLI Setup

```bash
go build -o ghs
cp ghs /usr/local/bin

# Add environment variables
echo 'export DAYS=' >> ~/.zshrc
echo 'export GITHUB_ORGANIZATION=' >> ~/.zshrc
echo 'export GITHUB_PERSONAL_ACCESS_TOKEN=' >> ~/.zshrc

# Update PATH to include the ghs binary directory
echo 'export PATH=$PATH:/usr/local/bin/ghs' >> ~/.zshrc

source ~/.zshrc
```

<!-- 
TODO: 
- For each of the results, show all PRs and Issues in the last X days.
- Just show top PRs and Issues?
- Rank repo or issues/prs by activity?
- Activity includes #commit, #prs, #issues ?
- Be able to show old PRs which are being commented on?

- Use a repo's pushed_at or updated_at to quickly filter out repos?
  - https://stackoverflow.com/questions/15918588/github-api-v3-what-is-the-difference-between-pushed-at-and-updated-at
  - https://docs.github.com/en/rest/repos/repos?apiVersion=2022-11-28#list-organization-repositories

- Event types include Push, PR Comment, PR Review, 

	// client.Repositories.ListCommitActivity()
	// client.PullRequests.List()
	// client.Issues.List()
	// client.Activity.ListRepositoryNotifications()

- Frontend?
-->

<!-- 
DONE:
- First start by listing repositories which had the most activity in the past DAYS
- Don't list repos which have zero activity.
-->