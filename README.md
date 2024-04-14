# Github Week in Review

Spotify Wrapped meets Github Activity. A CLI tool that generates an overview of a Github organization's activity over the last X days. Along with links to its most active Issues and PRs.

```sh
$ gwir

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
  - [Add config option to not query for local disk cost](https://github.com/opencost/opencost/pull/2441) : 4

opencost/opencost-helm-chart. TotalEvents=61
  - PullRequestReviewCommentEvent : 23
  - PullRequestReviewEvent : 23
  - IssueCommentEvent : 5
  - PullRequestEvent : 3
  - ForkEvent : 3
  - PushEvent : 2
  - CreateEvent : 1
  - ReleaseEvent : 1
Top PRs/Issues:
  - [Add Parquet-Exporter helm chart](https://github.com/opencost/opencost-helm-chart/pull/195) : 24
  - [Option to use existing kubernetes secret ](https://github.com/opencost/opencost-helm-chart/pull/196) : 20
  - [chore(kubeStateMetrics): make it clear that the keys of kubeStat](https://github.com/opencost/opencost-helm-chart/pull/193) : 9
  - [CloudCost QueryAthena Issues](https://github.com/opencost/opencost-helm-chart/issues/194) : 1

opencost/opencost-parquet-exporter. TotalEvents=14
  - PullRequestReviewEvent : 5
  - IssueCommentEvent : 4
  - PullRequestReviewCommentEvent : 4
  - PullRequestEvent : 1
Top PRs/Issues:
  - [Add azure storage account export target](https://github.com/opencost/opencost-parquet-exporter/pull/11) : 12
  - [Add helm chart](https://github.com/opencost/opencost-parquet-exporter/pull/10) : 2

opencost/opencost-grafana-dashboard. TotalEvents=1
  - WatchEvent : 1
Top PRs/Issues:
```

## Usage

```bash
export DAYS=
export GITHUB_ORGANIZATION=
export GITHUB_PERSONAL_ACCESS_TOKEN=  # optional

go run main.go
```

## CLI Setup

```bash
go build -o gwir
sudo cp gwir /usr/local/bin

# Add environment variables
echo 'export DAYS=' >> ~/.zshrc
echo 'export GITHUB_ORGANIZATION=' >> ~/.zshrc
echo 'export GITHUB_PERSONAL_ACCESS_TOKEN=' >> ~/.zshrc

source ~/.zshrc
```

<!--
TODO: 
- Take params via args instead of env vars.
- Define defaults for the env vars, so the user doens't have to set them.
- Make all PR/Issue titles the same length to "prettify" the output?
- Use a repo's pushed_at or updated_at to quickly filter out repos?
  - https://stackoverflow.com/questions/15918588/github-api-v3-what-is-the-difference-between-pushed-at-and-updated-at
  - https://docs.github.com/en/rest/repos/repos?apiVersion=2022-11-28#list-organization-repositories
- TUI
  - https://github.com/avelino/awesome-go?tab=readme-ov-file#command-line
  - https://github.com/charmbracelet/bubbletea
- Frontend?
-->

<!-- 
DONE (newest to oldest):
- Concurrency
- Include a snippet of the name of the PR/Issue.
- Only show top X PRs and Issues?
- Order the PRs and Issues
- For each of the results, show all PRs and Issues in the last X days.
- First start by listing repositories which had the most activity in the past DAYS
- Don't list repos which have zero activity.
-->