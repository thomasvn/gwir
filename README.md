# Github Week in Review

Spotify Wrapped meets Github Activity. A CLI tool that generates an overview of a Github organization's activity over the last X days. Along with links to its most active Issues and PRs.

```txt
$ gwir -org opencost

## Processing ...

opencost/opencost-helm-chart TotalEvents=1
opencost/opencost-website TotalEvents=1
opencost/opencost-infra TotalEvents=94
opencost/opencost-ui TotalEvents=5
opencost/opencost TotalEvents=232

## Ordered Results

### opencost/opencost TotalEvents=232

Event Types:
- IssueCommentEvent: 47
- PushEvent: 44
- PullRequestEvent: 34
- CreateEvent: 28
- PullRequestReviewEvent: 20
- WatchEvent: 19
- DeleteEvent: 18
- PullRequestReviewCommentEvent: 13
- ForkEvent: 7
- IssuesEvent: 2

Top PRs/Issues:
- [cleanup usage for less repeated code, and ena...](https://github.com/opencost/opencost/pull/3112): 24
- [continued integration test runner debugging     ](https://github.com/opencost/opencost/pull/3126): 10
- [update permissions to properly run steps        ](https://github.com/opencost/opencost/pull/3132): 8

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
        GitHub Personal Access Token (can also be set via GITHUB_PERSONAL_ACCESS_TOKEN env var)
  -top int
        How many top PRs/Issues to show (default 5)
  -usr string
        GitHub user to analyze
```

## Install

```bash
ARCH=macos-arm64  # macos-amd64, linux-amd64, windows-amd64

curl -L -O https://github.com/thomasvn/gwir/releases/latest/download/gwir.$ARCH.tar.gz
tar -xvf gwir.$ARCH.tar.gz
sudo mv gwir /usr/local/bin
```

Optionally, set up your Github Personal Access Token (PAT) following [this doc](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens). The PAT will need the following permission: `repo (Full control of private repositories)`.

<!--
TODO: 
- Enrich data when -usr flag is passsed? Not all events have associated HTML URLs.
- Validate the PAT has sufficient permissions?
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
