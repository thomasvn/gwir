# Github Week in Review

Spotify Wrapped meets Github Activity. A CLI tool that generates an overview of a Github organization's activity over the last X days. Along with links to its most active Issues and PRs.

```txt
$ gwir -org opencost

## Processing ...

opencost/opencost-grafana-dashboard. TotalEvents=16
opencost/opencost-parquet-exporter. TotalEvents=1
opencost/opencost-website. TotalEvents=1
opencost/opencost-helm-chart. TotalEvents=35
opencost/opencost-plugins. TotalEvents=31
opencost/opencost. TotalEvents=101

## Ordered results ...

### opencost/opencost. TotalEvents=101
  - IssueCommentEvent : 44
  - PushEvent : 18
  - PullRequestEvent : 13
  - WatchEvent : 11
  - CreateEvent : 7
  - DeleteEvent : 4
  - IssuesEvent : 2
  - PullRequestReviewEvent : 2
Top PRs/Issues:
  - [TypeUtil Enhancements                           ](https://github.com/opencost/opencost/pull/2707) : 6
  - [Intermittent "error":"vector cannot contain m...](https://github.com/opencost/opencost/issues/2704) : 6
  - [AWS IRSA authorizer for cloud integrations      ](https://github.com/opencost/opencost/pull/2710) : 5

...

```

## Usage

```bash
$ gwir -h
  -days int
    	How many days back to analyze (default 7)
  -org string
    	GitHub organization to analyze
  -token string
    	Optional. Passing a GitHub Personal Access Token allows you to view private repositories and make more API requests per hour. You can also set this token as an environment variable GITHUB_PERSONAL_ACCESS_TOKEN.
  -top int
    	How many top PRs/Issues to show (default 5)
  -usr string
    	GitHub user to analyze
```

## Install

```bash
RELEASE=v0.2
ARCH=macos-arm64  # macos-amd64, linux-amd64, windows-amd64

curl -L -O https://github.com/thomasvn/gwir/releases/download/$RELEASE/gwir.$ARCH.tar.gz
tar -xvf gwir.$ARCH.tar.gz
sudo mv gwir /usr/local/bin
```

<!--
TODO: 
- Enrich data when -usr flag is passsed? Not all events have associated HTML URLs.
- TUI
  - https://github.com/avelino/awesome-go?tab=readme-ov-file#command-line
  - https://github.com/charmbracelet/bubbletea
- Frontend?
- Provide a --version flag
- Automate releases via Github workflows?
- CLI tool downloadable via `go get` or `brew install`
- Use a repo's pushed_at or updated_at to quickly filter out repos?
  - https://stackoverflow.com/questions/15918588/github-api-v3-what-is-the-difference-between-pushed-at-and-updated-at
  - https://docs.github.com/en/rest/repos/repos?apiVersion=2022-11-28#list-organization-repositories
- Use a pretty image for the README. https://github.com/charmbracelet/vhs. Keep it up to date with vhs-actions
  - Asciicinema? https://github.com/kubecost/kubectl-cost/blob/main/assets/presentation-script.md
- Pipe to Glow?
  - echo "[Glow](https://github.com/charmbracelet/glow)" | glow -
- Other APIs to investigate.
    // client.Activity.ListEventsPerformedByUser()
    // client.Activity.ListEventsReceivedByUser()
    // client.Activity.ListUserEventsForOrganization()
    // client.Activity.ListEventsForOrganization()
    // client.Activity.ListFeeds()
-->

<!-- 
DONE (newest to oldest):
- First implementation of -usr flag.
- Prettify output. Specifically PR/Issue title length?
- Github MultiArch releases
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
