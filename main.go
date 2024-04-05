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
	// TODO: Check `repos.UpdatedAt`. Only add ones which have been updated in past few days. Prevent extra work later on.

	fmt.Println("Processing ...")
	type repoCount struct {
		RepoName string
		Count    int
	}
	repoCounts := []repoCount{}
	for _, repo := range allRepos {
		numRepoEvents := getNumRepoEventsLastXDays(client, repo.Owner.GetLogin(), repo.GetName(), DAYS_INT)
		repoCounts = append(repoCounts, repoCount{RepoName: repo.GetName(), Count: numRepoEvents})
		fmt.Println(repo.Owner.GetLogin(), "/", repo.GetName(), "=", numRepoEvents)
	}

	fmt.Println("\nOrdered results:")
	sort.Slice(repoCounts, func(i, j int) bool {
		return repoCounts[i].Count > repoCounts[j].Count
	})
	for _, repoCount := range repoCounts {
		if repoCount.Count == 0 {
			break
		}
		fmt.Println(repoCount.RepoName, "=", repoCount.Count)
	}

	// TODO: It's not pulling all PRs correctly.
	fmt.Println("\nRecent PRs:")
	for _, repoCount := range repoCounts {
		if repoCount.Count == 0 {
			break
		}
		fmt.Println(repoCount.RepoName, "(", repoCount.Count, ")", ":")
		prs := getRepoPRsLastXDays(client, GITHUB_ORGANIZATION, repoCount.RepoName, DAYS_INT)
		for _, pr := range prs {
			fmt.Println(" - ", pr.GetTitle(), ":", pr.GetHTMLURL())
		}
	}
}

func getNumRepoEventsLastXDays(client *github.Client, owner string, repo string, x int) int {
	numEvents := 0
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
			numEvents++
		}
		if stop {
			break
		}
		if res.NextPage == 0 {
			break
		}
		page++
	}
	return numEvents
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
