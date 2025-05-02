package github

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/google/go-github/v71/github"
)

type RepoStats struct {
	Name          string
	EventCounts   map[string]int
	PRIssueCounts map[string]int
	PRIssueTitles map[string]string
	TotalEvents   int
}

func (a *Analyzer) AnalyzeOrg() error {
	allRepos, err := a.fetchAllRepos()
	if err != nil {
		return err
	}

	fmt.Printf("\n\x1b[1;36m## Processing ...\x1b[0m\n\n")

	results := make(chan RepoStats, len(allRepos))
	var wg sync.WaitGroup

	for _, repo := range allRepos {
		wg.Add(1)
		go func(repo *github.Repository) {
			defer wg.Done()
			a.processRepo(repo, results)
		}(repo)
	}

	wg.Wait()
	close(results)

	var stats []RepoStats
	for repoStats := range results {
		stats = append(stats, repoStats)
	}

	a.displayResults(stats)
	return nil
}

func (a *Analyzer) fetchAllRepos() ([]*github.Repository, error) {
	var allRepos []*github.Repository
	page := 1

	for {
		opts := &github.RepositoryListByOrgOptions{
			ListOptions: github.ListOptions{
				PerPage: 100,
				Page:    page,
			},
		}

		repos, res, err := a.client.Repositories.ListByOrg(context.Background(), a.config.GitHubOrganization, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to list org repositories: %w", err)
		}

		allRepos = append(allRepos, repos...)
		if res.NextPage == 0 {
			break
		}
		page++
	}

	return allRepos, nil
}

func (a *Analyzer) processRepo(repo *github.Repository, results chan<- RepoStats) {
	eventCounts := make(map[string]int)
	prIssueCounts := make(map[string]int)
	prIssueTitles := make(map[string]string)
	totalCount := 0

	owner := repo.Owner.GetLogin()
	repoName := repo.GetName()

	page := 1
	for {
		opts := github.ListOptions{PerPage: 100, Page: page}
		events, res, err := a.client.Activity.ListRepositoryEvents(context.Background(), owner, repoName, &opts)
		if err != nil {
			fmt.Printf("Error fetching events for %s/%s: %v\n", owner, repoName, err)
			return
		}

		for _, event := range events {
			if !isWithinTimeRange(event, a.config.Days) {
				break
			}

			eventCounts[event.GetType()]++
			totalCount++

			payload, err := event.ParsePayload()
			if err != nil {
				continue
			}

			switch e := payload.(type) {
			case *github.PullRequestEvent:
				updatePRIssueStats(prIssueCounts, prIssueTitles, e.PullRequest)
			case *github.PullRequestReviewEvent:
				updatePRIssueStats(prIssueCounts, prIssueTitles, e.PullRequest)
			case *github.PullRequestReviewCommentEvent:
				updatePRIssueStats(prIssueCounts, prIssueTitles, e.PullRequest)
			case *github.PullRequestReviewThreadEvent:
				updatePRIssueStats(prIssueCounts, prIssueTitles, e.PullRequest)
			case *github.PullRequestTargetEvent:
				updatePRIssueStats(prIssueCounts, prIssueTitles, e.PullRequest)
			case *github.IssuesEvent:
				updatePRIssueStats(prIssueCounts, prIssueTitles, e.Issue)
			case *github.IssueCommentEvent:
				updatePRIssueStats(prIssueCounts, prIssueTitles, e.Issue)
			}
		}

		if res.NextPage == 0 {
			break
		}
		page = res.NextPage
	}

	if totalCount > 0 {
		results <- RepoStats{
			Name:          owner + "/" + repoName,
			EventCounts:   eventCounts,
			PRIssueCounts: prIssueCounts,
			PRIssueTitles: prIssueTitles,
			TotalEvents:   totalCount,
		}
		fmt.Printf("\x1b[1;33m%s/%s\x1b[0m \x1b[1;37mTotalEvents=\x1b[1;32m%d\x1b[0m\n", owner, repoName, totalCount)
	}
}

func updatePRIssueStats(counts map[string]int, titles map[string]string, issue interface{}) {
	var url, title string
	switch i := issue.(type) {
	case *github.PullRequest:
		url = i.GetHTMLURL()
		title = i.GetTitle()
	case *github.Issue:
		url = i.GetHTMLURL()
		title = i.GetTitle()
	}
	counts[url]++
	titles[url] = title
}

func (a *Analyzer) displayResults(stats []RepoStats) {
	fmt.Printf("\n\x1b[1;36m## Ordered Results\x1b[0m\n\n")
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].TotalEvents > stats[j].TotalEvents
	})

	for _, repo := range stats {
		fmt.Printf("\x1b[1;33m### %s\x1b[0m \x1b[1;37mTotalEvents=\x1b[1;32m%d\x1b[0m\n\n",
			repo.Name, repo.TotalEvents)

		fmt.Printf("\x1b[1;37mEvent Types:\x1b[0m\n")
		for _, pair := range sortMap(repo.EventCounts) {
			fmt.Printf("- \x1b[1;34m%s\x1b[0m: \x1b[1;32m%d\x1b[0m\n", pair.Key, pair.Value)
		}

		fmt.Printf("\n\x1b[1;37mTop PRs/Issues:\x1b[0m\n")
		count := 0
		for _, pair := range sortMap(repo.PRIssueCounts) {
			title := trimString(repo.PRIssueTitles[pair.Key], 48)
			fmt.Printf("- [\x1b[1;34m%s\x1b[0m](%s): \x1b[1;32m%d\x1b[0m\n",
				title, pair.Key, pair.Value)
			count++
			if count >= a.config.TopXActivities {
				break
			}
		}
		fmt.Printf("\n")
	}
}

func trimString(s string, n int) string {
	if len(s) > n {
		return s[:n-3] + "..."
	}
	return s + strings.Repeat(" ", n-len(s))
}
