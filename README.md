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
- First start by listing repositories which had the most activity in the past DAYS

- Rank by activity?
- Activity includes #commit, #prs, #issues ?

- Event types include Push, PR Comment, PR Review, 

	// client.Repositories.ListCommitActivity()
	// client.PullRequests.List()
	// client.Issues.List()
	// client.Activity.ListRepositoryNotifications()
-->