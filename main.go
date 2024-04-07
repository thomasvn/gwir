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
		TotalEvents int
	}
	repoEventCounts := []RepoEventCount{}
	for _, repo := range allRepos {
		eventsMap, totalCount := getRepoEventsLastXDays(client, repo.Owner.GetLogin(), repo.GetName(), DAYS_INT)
		if totalCount > 0 {
			repoEventCounts = append(repoEventCounts, RepoEventCount{RepoName: repo.GetName(), EventsMap: eventsMap, TotalEvents: totalCount})
			fmt.Printf("%s/%s. TotalEvents=%d\n", repo.Owner.GetLogin(), repo.GetName(), totalCount)
			for k, v := range eventsMap {
				fmt.Printf("  - %s=%d\n", k, v)
			}
		}
	}

	// Order the results from above by TotalEvents
	fmt.Println("\nOrdered results:")
	sort.Slice(repoEventCounts, func(i, j int) bool {
		return repoEventCounts[i].TotalEvents > repoEventCounts[j].TotalEvents
	})
	for _, repoEventCount := range repoEventCounts {
		fmt.Printf("%s. TotalEvents=%d\n", repoEventCount.RepoName, repoEventCount.TotalEvents)
		for k, v := range repoEventCount.EventsMap {
			fmt.Printf("  - %s=%d\n", k, v)
		}
	}

	// // TODO: It's not pulling all PRs correctly.
	// fmt.Println("\nRecent PRs:")
	// for _, repoCount := range repoCounts {
	// 	if repoCount.Count == 0 {
	// 		break
	// 	}
	// 	fmt.Println(repoCount.RepoName, "(", repoCount.Count, ")", ":")
	// 	prs := getRepoPRsLastXDays(client, GITHUB_ORGANIZATION, repoCount.RepoName, DAYS_INT)
	// 	for _, pr := range prs {
	// 		fmt.Println(" - ", pr.GetTitle(), ":", pr.GetHTMLURL())
	// 	}
	// }
}

func getRepoEventsLastXDays(client *github.Client, owner string, repo string, x int) (map[string]int, int) {
	totalCount := 0
	results := make(map[string]int)
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
			if val, ok := results[event.GetType()]; ok {
				results[event.GetType()] = val + 1
			} else {
				results[event.GetType()] = 1
			}
			totalCount++

			// TODO: Enrich information about events
			// payload, _ := event.ParsePayload()
			// fmt.Printf("\tDEBUG. Type: %v, Event payload: %v\n", reflect.TypeOf(payload), payload)
			// switch event.GetType() {
			// case "PullRequestEvent":
			// 	prEvent := payload.(*github.PullRequestEvent)

		}
		if stop {
			break
		}
		if res.NextPage == 0 {
			break
		}
		page++
	}
	return results, totalCount
}

func getRepoPRsLastXDays(client *github.Client, owner string, repo string, x int) []*github.PullRequest {
	prs := []*github.PullRequest{}
	stop := false
	page := 1
	for {
		opts := &github.PullRequestListOptions{}
		opts.State = "all"
		opts.Sort = "updated"
		opts.ListOptions.PerPage = 100
		opts.ListOptions.Page = page
		pulls, res, _ := client.PullRequests.List(context.Background(), owner, repo, opts)
		for _, pull := range pulls {
			if pull.GetUpdatedAt().Time.Before(time.Now().AddDate(0, 0, -1*x)) {
				stop = true
				break
			}
			prs = append(prs, pull)
		}
		if stop {
			break
		}
		if res.NextPage == 0 {
			break
		}
		page++
	}
	return prs
}
