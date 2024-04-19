# Github Week in Review

Spotify Wrapped meets Github Activity. A CLI tool that generates an overview of a Github organization's activity over the last X days. Along with links to its most active Issues and PRs.

```txt
$ gwir -org opencost

## Processing ...

opencost/opencost. TotalEvents=70
opencost/opencost-helm-chart. TotalEvents=61
opencost/opencost-parquet-exporter. TotalEvents=14
opencost/opencost-grafana-dashboard. TotalEvents=1

## Ordered results:

opencost/opencost. TotalEvents=70
  - IssueCommentEvent : 38
  - WatchEvent : 16
  - PullRequestEvent : 5
  - IssuesEvent : 5
  - CreateEvent : 3
  - PullRequestReviewEvent : 2
  - DeleteEvent : 1
Top PRs/Issues:
  - [Support for specifying and attributing shared fixed costs](https://github.com/opencost/opencost/issues/2427) : 6
  - [Create a hash key when agg properties are not set for Cloud Cost](https://github.com/opencost/opencost/pull/2700) : 5
  - [Provider alibaba support RRSA authentication](https://github.com/opencost/opencost/issues/2699) : 5
  - [`QueryAthenaPaginated: start query error: not found, ResolveEndp](https://github.com/opencost/opencost/issues/2697) : 5

...

```

## Usage

```bash
$ gwir -h
Usage of gwir:
  -days int
    	How many days back to analyze (default 7)
  -org string
    	GitHub organization to analyze
  -token string
    	Optional. Passing a GitHub Personal Access Token allows you to view private repositories and make more API requests per hour. You can also set this token as an environment variable GITHUB_PERSONAL_ACCESS_TOKEN.
  -top int
    	How many top PRs/Issues to show (default 5)
```

## Install

```bash
curl -L -O https://github.com/thomasvn/gwir/releases/download/v0.1/gwir.darwin-amd64.zip
unzip gwir.darwin-amd64.zip
sudo mv gwir /usr/local/bin
```

<!--
TODO: 
- Github Releases
  - GOOS=windows GOARCH=amd64 go build -o gwir.exe
  - GOOS=darwin GOARCH=amd64 go build -o gwir
  - GOOS=linux GOARCH=amd64 go build -o gwir
  - Automate using Github Workflows
- Provide a --version flag
- CLI tool downloadable via `go get` or `brew install`
- Make all PR/Issue titles the same length to "prettify" the output?
- Use a repo's pushed_at or updated_at to quickly filter out repos?
  - https://stackoverflow.com/questions/15918588/github-api-v3-what-is-the-difference-between-pushed-at-and-updated-at
  - https://docs.github.com/en/rest/repos/repos?apiVersion=2022-11-28#list-organization-repositories
- TUI
  - https://github.com/avelino/awesome-go?tab=readme-ov-file#command-line
  - https://github.com/charmbracelet/bubbletea
- Frontend?
- Use a pretty image for the README. https://github.com/charmbracelet/vhs. Keep it up to date with vhs-actions
  - Asciicinema? https://github.com/kubecost/kubectl-cost/blob/main/assets/presentation-script.md
- Pipe to Glow?
  - echo "[Glow](https://github.com/charmbracelet/glow)" | glow -
-->

<!-- 
DONE (newest to oldest):
- Take params via args instead of env vars.
  - https://pkg.go.dev/flag
  - https://github.com/avelino/awesome-go?tab=readme-ov-file#standard-cli
- Concurrency
- Include a snippet of the name of the PR/Issue.
- Only show top X PRs and Issues?
- Order the PRs and Issues
- For each of the results, show all PRs and Issues in the last X days.
- First start by listing repositories which had the most activity in the past DAYS
- Don't list repos which have zero activity.
-->