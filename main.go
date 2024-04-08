package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/google/go-github/v61/github"
)

func main() {
	DAYS := os.Getenv("DAYS")
	DAYS_INT, _ := strconv.Atoi(DAYS)
	GITHUB_ORGANIZATION := os.Getenv("GITHUB_ORGANIZATION")
	GITHUB_PERSONAL_ACCESS_TOKEN := os.Getenv("GITHUB_PERSONAL_ACCESS_TOKEN")

	var client *github.Client
	if GITHUB_PERSONAL_ACCESS_TOKEN == "" {
		client = github.NewClient(nil)
	} else {
		client = github.NewClient(nil).WithAuthToken(GITHUB_PERSONAL_ACCESS_TOKEN)
	}

	// Get all repos in this GITHUB_ORGANIZATION
	allRepos := []*github.Repository{}
	page := 1
	for {
		opts := &github.RepositoryListByOrgOptions{}
		opts.ListOptions.PerPage = 100
		opts.ListOptions.Page = page
		repos, res, _ := client.Repositories.ListByOrg(context.Background(), GITHUB_ORGANIZATION, opts)
		allRepos = append(allRepos, repos...)
		if res.NextPage == 0 {
			break
		}
		page++
	}

	// For each repo, print the number/type of events in the last X days
	fmt.Println("\nProcessing ...")
	type RepoEventCount struct {
		RepoName    string
		EventsMap   map[string]int
		PRIssuesMap map[string]int
		TotalEvents int
	}
	repoEventCounts := []RepoEventCount{}
	for _, repo := range allRepos {
		eventCounts, prIssueCounts, totalCount := getRepoEventsLastXDays(client, repo.Owner.GetLogin(), repo.GetName(), DAYS_INT)
		if totalCount > 0 {
			repoEventCounts = append(repoEventCounts, RepoEventCount{RepoName: repo.Owner.GetLogin() + "/" + repo.GetName(), EventsMap: eventCounts, PRIssuesMap: prIssueCounts, TotalEvents: totalCount})
			fmt.Printf("%s/%s. TotalEvents=%d\n", repo.Owner.GetLogin(), repo.GetName(), totalCount)
		}
	}

	// Order the results from above by TotalEvents
	fmt.Println("\nOrdered results:")
	sort.Slice(repoEventCounts, func(i, j int) bool {
		return repoEventCounts[i].TotalEvents > repoEventCounts[j].TotalEvents
	})
	for _, repoEventCount := range repoEventCounts {
		fmt.Printf("\n%s. TotalEvents=%d\n", repoEventCount.RepoName, repoEventCount.TotalEvents)
		for k, v := range repoEventCount.EventsMap {
			fmt.Printf("  - %s=%d\n", k, v)
		}
		fmt.Printf("PRs/Issues:\n")
		for k, v := range repoEventCount.PRIssuesMap {
			fmt.Printf("  - %s : %d\n", k, v)
		}
	}
}

// getRepoEventsLastXDays analyzes all activity in a repo within the last x
// days. It returns 1) a map of event types and their counts, 2) a map of
// PRs/Issues and their event counts, and 3) the total number of events.
func getRepoEventsLastXDays(client *github.Client, owner string, repo string, x int) (map[string]int, map[string]int, int) {
	eventCounts := make(map[string]int)
	prIssueCounts := make(map[string]int)
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
			case "PullRequestReviewEvent":
				prReviewEvent := payload.(*github.PullRequestReviewEvent)
				if val, ok := prIssueCounts[prReviewEvent.PullRequest.GetHTMLURL()]; ok {
					prIssueCounts[prReviewEvent.PullRequest.GetHTMLURL()] = val + 1
				} else {
					prIssueCounts[prReviewEvent.PullRequest.GetHTMLURL()] = 1
				}
			case "PullRequestReviewCommentEvent":
				prReviewCommentEvent := payload.(*github.PullRequestReviewCommentEvent)
				if val, ok := prIssueCounts[prReviewCommentEvent.PullRequest.GetHTMLURL()]; ok {
					prIssueCounts[prReviewCommentEvent.PullRequest.GetHTMLURL()] = val + 1
				} else {
					prIssueCounts[prReviewCommentEvent.PullRequest.GetHTMLURL()] = 1
				}
			case "PullRequestReviewThreadEvent":
				prReviewThreadEvent := payload.(*github.PullRequestReviewThreadEvent)
				if val, ok := prIssueCounts[prReviewThreadEvent.PullRequest.GetHTMLURL()]; ok {
					prIssueCounts[prReviewThreadEvent.PullRequest.GetHTMLURL()] = val + 1
				} else {
					prIssueCounts[prReviewThreadEvent.PullRequest.GetHTMLURL()] = 1
				}
			case "PullRequestTargetEvent":
				prTargetEvent := payload.(*github.PullRequestTargetEvent)
				if val, ok := prIssueCounts[prTargetEvent.PullRequest.GetHTMLURL()]; ok {
					prIssueCounts[prTargetEvent.PullRequest.GetHTMLURL()] = val + 1
				} else {
					prIssueCounts[prTargetEvent.PullRequest.GetHTMLURL()] = 1
				}
			case "IssuesEvent":
				issuesEvent := payload.(*github.IssuesEvent)
				if val, ok := prIssueCounts[issuesEvent.Issue.GetHTMLURL()]; ok {
					prIssueCounts[issuesEvent.Issue.GetHTMLURL()] = val + 1
				} else {
					prIssueCounts[issuesEvent.Issue.GetHTMLURL()] = 1
				}
			case "IssueCommentEvent":
				issueCommentEvent := payload.(*github.IssueCommentEvent)
				if val, ok := prIssueCounts[issueCommentEvent.Issue.GetHTMLURL()]; ok {
					prIssueCounts[issueCommentEvent.Issue.GetHTMLURL()] = val + 1
				} else {
					prIssueCounts[issueCommentEvent.Issue.GetHTMLURL()] = 1
				}
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

	return eventCounts, prIssueCounts, totalCount
}
