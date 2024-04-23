package main

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/go-github/v61/github"
)

// AnalyzeOrgActivity looks at all repos in an organization. For each of the
// repos, it prints the number and type of activities. It also prints the top X
// PRs/Issues based on the number of events associated with the PR/Issue.
func AnalyzeOrgActivity(client *github.Client, argopts ArgOpts) {
	// Get all repos in this GITHUB_ORGANIZATION
	allRepos := []*github.Repository{}
	page := 1
	for {
		opts := &github.RepositoryListByOrgOptions{}
		opts.ListOptions.PerPage = 100
		opts.ListOptions.Page = page
		repos, res, _ := client.Repositories.ListByOrg(context.Background(), argopts.GITHUB_ORGANIZATION, opts)
		allRepos = append(allRepos, repos...)
		if res.NextPage == 0 {
			break
		}
		page++
	}

	// For each repo, print the number/type of events in the last X days
	fmt.Printf("\n## Processing ... \n\n")
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
			eventCounts, prIssueCounts, prIssueTitles, totalCount := getRepoEventsLastXDays(client, repo.Owner.GetLogin(), repo.GetName(), argopts.DAYS)
			if totalCount > 0 {
				repoEventCountsChan <- RepoEventCount{
					RepoName:          repo.Owner.GetLogin() + "/" + repo.GetName(),
					EventTypeCount:    eventCounts,
					PRIssueEventCount: prIssueCounts,
					PRIssueTitle:      prIssueTitles,
					TotalEvents:       totalCount,
				}
				fmt.Printf("%s/%s. TotalEvents=%d\n", repo.Owner.GetLogin(), repo.GetName(), totalCount)
			}
		}(repo)
	}
	wg.Wait()
	close(repoEventCountsChan)
	repoEventCounts := []RepoEventCount{}
	for repoEventCount := range repoEventCountsChan {
		repoEventCounts = append(repoEventCounts, repoEventCount)
	}

	// Order the results from above by TotalEvents
	fmt.Printf("\n## Ordered results ... \n")
	sort.Slice(repoEventCounts, func(i, j int) bool {
		return repoEventCounts[i].TotalEvents > repoEventCounts[j].TotalEvents
	})
	for _, repoEventCount := range repoEventCounts {
		fmt.Printf("\n### %s. TotalEvents=%d  \n", repoEventCount.RepoName, repoEventCount.TotalEvents)
		EventTypeCountSortedSlice := sortMap(repoEventCount.EventTypeCount)
		for _, pair := range EventTypeCountSortedSlice {
			fmt.Printf("  - %s : %d\n", pair.Key, pair.Value)
		}
		fmt.Printf("Top PRs/Issues:  \n")
		count := 0
		PRIssuesSortedSlice := sortMap(repoEventCount.PRIssueEventCount)
		for _, pair := range PRIssuesSortedSlice {
			title := trimString(repoEventCount.PRIssueTitle[pair.Key], 48)
			fmt.Printf("  - [%s](%s) : %d\n", title, pair.Key, pair.Value)
			count++
			if count >= argopts.TOPXACTIVITIES {
				break
			}
		}
	}
}

// getRepoEventsLastXDays analyzes all activity in a repo within the last x
// days. It returns 1) a map of event types and their counts, 2) a map of
// PRs/Issues and their event counts, 3) a map of PRs/Issues URLs and their
// titles, and 4) the total number of events.
func getRepoEventsLastXDays(client *github.Client, owner string, repo string, x int) (map[string]int, map[string]int, map[string]string, int) {
	eventCounts := make(map[string]int)
	prIssueCounts := make(map[string]int)
	prIssueTitle := make(map[string]string)
	totalCount := 0

	// Paginated API queries against ListRepositoryEvents()
	stop := false
	page := 1
	for {
		opts := &github.ListOptions{}
		opts.PerPage = 100
		opts.Page = page
		events, res, _ := client.Activity.ListRepositoryEvents(context.Background(), owner, repo, opts)
		for _, event := range events {
			if event.GetCreatedAt().Time.Before(time.Now().AddDate(0, 0, -1*x)) {
				stop = true
				break
			}

			// Tally the event counts
			if val, ok := eventCounts[event.GetType()]; ok {
				eventCounts[event.GetType()] = val + 1
			} else {
				eventCounts[event.GetType()] = 1
			}
			totalCount++

			// Tally PR/Issue event counts
			payload, _ := event.ParsePayload()
			switch event.GetType() {
			case "PullRequestEvent":
				prEvent := payload.(*github.PullRequestEvent)
				if val, ok := prIssueCounts[prEvent.PullRequest.GetHTMLURL()]; ok {
					prIssueCounts[prEvent.PullRequest.GetHTMLURL()] = val + 1
				} else {
					prIssueCounts[prEvent.PullRequest.GetHTMLURL()] = 1
				}
				prIssueTitle[prEvent.PullRequest.GetHTMLURL()] = prEvent.PullRequest.GetTitle()
			case "PullRequestReviewEvent":
				prReviewEvent := payload.(*github.PullRequestReviewEvent)
				if val, ok := prIssueCounts[prReviewEvent.PullRequest.GetHTMLURL()]; ok {
					prIssueCounts[prReviewEvent.PullRequest.GetHTMLURL()] = val + 1
				} else {
					prIssueCounts[prReviewEvent.PullRequest.GetHTMLURL()] = 1
				}
				prIssueTitle[prReviewEvent.PullRequest.GetHTMLURL()] = prReviewEvent.PullRequest.GetTitle()
			case "PullRequestReviewCommentEvent":
				prReviewCommentEvent := payload.(*github.PullRequestReviewCommentEvent)
				if val, ok := prIssueCounts[prReviewCommentEvent.PullRequest.GetHTMLURL()]; ok {
					prIssueCounts[prReviewCommentEvent.PullRequest.GetHTMLURL()] = val + 1
				} else {
					prIssueCounts[prReviewCommentEvent.PullRequest.GetHTMLURL()] = 1
				}
				prIssueTitle[prReviewCommentEvent.PullRequest.GetHTMLURL()] = prReviewCommentEvent.PullRequest.GetTitle()
			case "PullRequestReviewThreadEvent":
				prReviewThreadEvent := payload.(*github.PullRequestReviewThreadEvent)
				if val, ok := prIssueCounts[prReviewThreadEvent.PullRequest.GetHTMLURL()]; ok {
					prIssueCounts[prReviewThreadEvent.PullRequest.GetHTMLURL()] = val + 1
				} else {
					prIssueCounts[prReviewThreadEvent.PullRequest.GetHTMLURL()] = 1
				}
				prIssueTitle[prReviewThreadEvent.PullRequest.GetHTMLURL()] = prReviewThreadEvent.PullRequest.GetTitle()
			case "PullRequestTargetEvent":
				prTargetEvent := payload.(*github.PullRequestTargetEvent)
				if val, ok := prIssueCounts[prTargetEvent.PullRequest.GetHTMLURL()]; ok {
					prIssueCounts[prTargetEvent.PullRequest.GetHTMLURL()] = val + 1
				} else {
					prIssueCounts[prTargetEvent.PullRequest.GetHTMLURL()] = 1
				}
				prIssueTitle[prTargetEvent.PullRequest.GetHTMLURL()] = prTargetEvent.PullRequest.GetTitle()
			case "IssuesEvent":
				issuesEvent := payload.(*github.IssuesEvent)
				if val, ok := prIssueCounts[issuesEvent.Issue.GetHTMLURL()]; ok {
					prIssueCounts[issuesEvent.Issue.GetHTMLURL()] = val + 1
				} else {
					prIssueCounts[issuesEvent.Issue.GetHTMLURL()] = 1
				}
				prIssueTitle[issuesEvent.Issue.GetHTMLURL()] = issuesEvent.Issue.GetTitle()
			case "IssueCommentEvent":
				issueCommentEvent := payload.(*github.IssueCommentEvent)
				if val, ok := prIssueCounts[issueCommentEvent.Issue.GetHTMLURL()]; ok {
					prIssueCounts[issueCommentEvent.Issue.GetHTMLURL()] = val + 1
				} else {
					prIssueCounts[issueCommentEvent.Issue.GetHTMLURL()] = 1
				}
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

type pair struct {
	Key   string
	Value int
}

// sortMap sorts a map by its values in descending order.
func sortMap(m map[string]int) []pair {
	pairs := []pair{}
	for k := range m {
		pairs = append(pairs, pair{k, m[k]})
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Value > pairs[j].Value
	})
	return pairs
}

// trimString trims the input string to the specified length n. If the input
// string is shorter than n, it pads the string with spaces until it reaches
// length n.
func trimString(s string, n int) string {
	if len(s) > n {
		return s[:n-3] + "..."
	}
	return s + strings.Repeat(" ", n-len(s))
}
