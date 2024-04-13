# Github Week in Review

Spotify Wrapped meets Github Activity. A CLI tool that generates an overview of a Github organization's activity over the last X days. Along with links to its most active Issues and PRs.

```sh
$ gwir

Processing ...
opencost/opencost. TotalEvents=146
opencost/opencost-website. TotalEvents=13
opencost/opencost-helm-chart. TotalEvents=35
opencost/opencost-parquet-exporter. TotalEvents=14
opencost/opencost-grafana-dashboard. TotalEvents=2

Ordered results:

opencost/opencost. TotalEvents=146
  - IssuesEvent=12
  - WatchEvent=18
  - PullRequestEvent=24
  - PullRequestReviewEvent=18
  - DeleteEvent=1
  - IssueCommentEvent=51
  - CreateEvent=5
  - PushEvent=11
  - PullRequestReviewCommentEvent=5
  - ForkEvent=1
PRs/Issues:
  - https://github.com/opencost/opencost/pull/2687 : 7
  - https://github.com/opencost/opencost/pull/2686 : 16
  - https://github.com/opencost/opencost/pull/2634 : 3
  - https://github.com/opencost/opencost/pull/2666 : 2
  - https://github.com/opencost/opencost/pull/2691 : 3
  - https://github.com/opencost/opencost/issues/2680 : 3
  - https://github.com/opencost/opencost/issues/2683 : 3
  - https://github.com/opencost/opencost/pull/2679 : 3
  - https://github.com/opencost/opencost/pull/2350 : 1
  - https://github.com/opencost/opencost/issues/2411 : 2
  - https://github.com/opencost/opencost/pull/2685 : 7
  - https://github.com/opencost/opencost/pull/2694 : 4
  - https://github.com/opencost/opencost/issues/2692 : 1
  - https://github.com/opencost/opencost/issues/2565 : 1
  - https://github.com/opencost/opencost/issues/2695 : 2
  - https://github.com/opencost/opencost/pull/2688 : 4
  - https://github.com/opencost/opencost/pull/2678 : 11
  - https://github.com/opencost/opencost/issues/1622 : 1
  - https://github.com/opencost/opencost/issues/2681 : 6
  - https://github.com/opencost/opencost/pull/2628 : 2
  - https://github.com/opencost/opencost/pull/2690 : 4
  - https://github.com/opencost/opencost/pull/2693 : 5
  - https://github.com/opencost/opencost/issues/2655 : 2
  - https://github.com/opencost/opencost/pull/2689 : 4
  - https://github.com/opencost/opencost/issues/2682 : 5
  - https://github.com/opencost/opencost/pull/2684 : 6
  - https://github.com/opencost/opencost/issues/2673 : 2

opencost/opencost-helm-chart. TotalEvents=35
  - PushEvent=5
  - CreateEvent=2
  - IssueCommentEvent=6
  - PullRequestReviewCommentEvent=2
  - PullRequestReviewEvent=7
  - PullRequestEvent=7
  - ForkEvent=2
  - ReleaseEvent=2
  - IssuesEvent=2
PRs/Issues:
  - https://github.com/opencost/opencost-helm-chart/pull/191 : 4
  - https://github.com/opencost/opencost-helm-chart/issues/189 : 1
  - https://github.com/opencost/opencost-helm-chart/pull/186 : 2
  - https://github.com/opencost/opencost-helm-chart/pull/190 : 2
  - https://github.com/opencost/opencost-helm-chart/issues/194 : 1
  - https://github.com/opencost/opencost-helm-chart/pull/193 : 2
  - https://github.com/opencost/opencost-helm-chart/pull/192 : 10
  - https://github.com/opencost/opencost-helm-chart/issues/188 : 2

opencost/opencost-parquet-exporter. TotalEvents=14
  - IssueCommentEvent=6
  - PullRequestEvent=4
  - IssuesEvent=2
  - PushEvent=1
  - PullRequestReviewEvent=1
PRs/Issues:
  - https://github.com/opencost/opencost-parquet-exporter/issues/8 : 4
  - https://github.com/opencost/opencost-parquet-exporter/pull/9 : 3
  - https://github.com/opencost/opencost-parquet-exporter/pull/11 : 2
  - https://github.com/opencost/opencost-parquet-exporter/pull/10 : 4

opencost/opencost-website. TotalEvents=13
  - IssueCommentEvent=2
  - IssuesEvent=1
  - PullRequestReviewEvent=2
  - PushEvent=3
  - PullRequestEvent=5
PRs/Issues:
  - https://github.com/opencost/opencost-website/pull/223 : 2
  - https://github.com/opencost/opencost-website/pull/225 : 3
  - https://github.com/opencost/opencost-website/issues/218 : 1
  - https://github.com/opencost/opencost-website/pull/224 : 4

opencost/opencost-grafana-dashboard. TotalEvents=2
  - IssuesEvent=2
PRs/Issues:
  - https://github.com/opencost/opencost-grafana-dashboard/issues/2 : 1
  - https://github.com/opencost/opencost-grafana-dashboard/issues/1 : 1
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
- Only show top X PRs and Issues?
- Take params via args instead of env vars.
- Include a snippet of the name of the PR/Issue.
- Concurrency
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
- Order the PRs and Issues
- For each of the results, show all PRs and Issues in the last X days.
- First start by listing repositories which had the most activity in the past DAYS
- Don't list repos which have zero activity.
-->