package github

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/google/go-github/v71/github"
)

func trimString(s string, n int) string {
	if len(s) > n {
		return s[:n-3] + "..."
	}
	return s + strings.Repeat(" ", n-len(s))
}

func (a *Analyzer) AnalyzeOrg() error {
	allRepos := []*github.Repository{}
	page := 1
	for {
		opts := &github.RepositoryListByOrgOptions{}
		opts.ListOptions.PerPage = 100
		opts.ListOptions.Page = page
		repos, res, err := a.client.Repositories.ListByOrg(context.Background(), a.config.GitHubOrganization, opts)
		if err != nil {
			return fmt.Errorf("failed to list org repositories: %w", err)
		}
		allRepos = append(allRepos, repos...)
		if res.NextPage == 0 {
			break
		}
		page++
	}

	fmt.Printf("\n\x1b[1;36m## Processing ...\x1b[0m\n\n")
	type RepoEventCount struct {
		RepoName          string
		EventTypeCount    map[string]int
		PRIssueEventCount map[string]int
		PRIssueTitle      map[string]string
		TotalEvents       int
	}

	var wg sync.WaitGroup
	repoEventCountsChan := make(chan RepoEventCount, len(allRepos))
	for _, repo := range allRepos {
		wg.Add(1)
		go func(repo *github.Repository) {
			defer wg.Done()
			eventCounts, prIssueCounts, prIssueTitles, totalCount := a.getRepoEventsLastXDays(repo.Owner.GetLogin(), repo.GetName())
			if totalCount > 0 {
				repoEventCountsChan <- RepoEventCount{
					RepoName:          repo.Owner.GetLogin() + "/" + repo.GetName(),
					EventTypeCount:    eventCounts,
					PRIssueEventCount: prIssueCounts,
					PRIssueTitle:      prIssueTitles,
					TotalEvents:       totalCount,
				}
				fmt.Printf("\x1b[1;33m%s/%s\x1b[0m \x1b[1;37mTotalEvents=\x1b[1;32m%d\x1b[0m\n",
					repo.Owner.GetLogin(), repo.GetName(), totalCount)
			}
		}(repo)
	}
	wg.Wait()
	close(repoEventCountsChan)

	repoEventCounts := []RepoEventCount{}
	for repoEventCount := range repoEventCountsChan {
		repoEventCounts = append(repoEventCounts, repoEventCount)
	}

	fmt.Printf("\n\x1b[1;36m## Ordered Results\x1b[0m\n\n")
	sort.Slice(repoEventCounts, func(i, j int) bool {
		return repoEventCounts[i].TotalEvents > repoEventCounts[j].TotalEvents
	})

	for _, repoEventCount := range repoEventCounts {
		fmt.Printf("\x1b[1;33m### %s\x1b[0m \x1b[1;37mTotalEvents=\x1b[1;32m%d\x1b[0m\n\n",
			repoEventCount.RepoName, repoEventCount.TotalEvents)

		EventTypeCountSortedSlice := sortMap(repoEventCount.EventTypeCount)
		fmt.Printf("\x1b[1;37mEvent Types:\x1b[0m\n")
		for _, pair := range EventTypeCountSortedSlice {
			fmt.Printf("- \x1b[1;34m%s\x1b[0m: \x1b[1;32m%d\x1b[0m\n", pair.Key, pair.Value)
		}

		fmt.Printf("\n\x1b[1;37mTop PRs/Issues:\x1b[0m\n")
		count := 0
		PRIssuesSortedSlice := sortMap(repoEventCount.PRIssueEventCount)
		for _, pair := range PRIssuesSortedSlice {
			title := trimString(repoEventCount.PRIssueTitle[pair.Key], 48)
			fmt.Printf("- [\x1b[1;34m%s\x1b[0m](%s): \x1b[1;32m%d\x1b[0m\n",
				title, pair.Key, pair.Value)
			count++
			if count >= a.config.TopXActivities {
				break
			}
		}
		fmt.Printf("\n")
	}

	return nil
}

func (a *Analyzer) getRepoEventsLastXDays(owner string, repo string) (map[string]int, map[string]int, map[string]string, int) {
	eventCounts := make(map[string]int)
	prIssueCounts := make(map[string]int)
	prIssueTitle := make(map[string]string)
	totalCount := 0

	stop := false
	page := 1
	for {
		opts := github.ListOptions{PerPage: 100, Page: page}
		events, res, _ := a.client.Activity.ListRepositoryEvents(context.Background(), owner, repo, &opts)
		for _, event := range events {
			if !isWithinTimeRange(event, a.config.Days) {
				stop = true
				break
			}

			eventCounts[event.GetType()]++
			totalCount++

			payload, _ := event.ParsePayload()
			switch event.GetType() {
			case "PullRequestEvent":
				prEvent := payload.(*github.PullRequestEvent)
				prIssueCounts[prEvent.PullRequest.GetHTMLURL()]++
				prIssueTitle[prEvent.PullRequest.GetHTMLURL()] = prEvent.PullRequest.GetTitle()
			case "PullRequestReviewEvent":
				prReviewEvent := payload.(*github.PullRequestReviewEvent)
				prIssueCounts[prReviewEvent.PullRequest.GetHTMLURL()]++
				prIssueTitle[prReviewEvent.PullRequest.GetHTMLURL()] = prReviewEvent.PullRequest.GetTitle()
			case "PullRequestReviewCommentEvent":
				prReviewCommentEvent := payload.(*github.PullRequestReviewCommentEvent)
				prIssueCounts[prReviewCommentEvent.PullRequest.GetHTMLURL()]++
				prIssueTitle[prReviewCommentEvent.PullRequest.GetHTMLURL()] = prReviewCommentEvent.PullRequest.GetTitle()
			case "PullRequestReviewThreadEvent":
				prReviewThreadEvent := payload.(*github.PullRequestReviewThreadEvent)
				prIssueCounts[prReviewThreadEvent.PullRequest.GetHTMLURL()]++
				prIssueTitle[prReviewThreadEvent.PullRequest.GetHTMLURL()] = prReviewThreadEvent.PullRequest.GetTitle()
			case "PullRequestTargetEvent":
				prTargetEvent := payload.(*github.PullRequestTargetEvent)
				prIssueCounts[prTargetEvent.PullRequest.GetHTMLURL()]++
				prIssueTitle[prTargetEvent.PullRequest.GetHTMLURL()] = prTargetEvent.PullRequest.GetTitle()
			case "IssuesEvent":
				issuesEvent := payload.(*github.IssuesEvent)
				prIssueCounts[issuesEvent.Issue.GetHTMLURL()]++
				prIssueTitle[issuesEvent.Issue.GetHTMLURL()] = issuesEvent.Issue.GetTitle()
			case "IssueCommentEvent":
				issueCommentEvent := payload.(*github.IssueCommentEvent)
				prIssueCounts[issueCommentEvent.Issue.GetHTMLURL()]++
				prIssueTitle[issueCommentEvent.Issue.GetHTMLURL()] = issueCommentEvent.Issue.GetTitle()
			}
		}
		if stop {
			break
		}
		if res.NextPage == 0 {
			break
		}
		page++
	}

	return eventCounts, prIssueCounts, prIssueTitle, totalCount
}
