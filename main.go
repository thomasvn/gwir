package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v61/github"
)

func main() {
	DAYS := os.Getenv("DAYS")
	GITHUB_ORGANIZATION := os.Getenv("GITHUB_ORGANIZATION")

	// TODO: PAT should be optional. Make best effort GET requests.
	GITHUB_PERSONAL_ACCESS_TOKEN := os.Getenv("GITHUB_PERSONAL_ACCESS_TOKEN")
	if GITHUB_PERSONAL_ACCESS_TOKEN == "" {
		fmt.Println("GITHUB_PERSONAL_ACCESS_TOKEN is not set")
		return
	}
	client := github.NewClient(nil).WithAuthToken(GITHUB_PERSONAL_ACCESS_TOKEN)

	// Compile list of repos in this organization using Paginated API
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
	for _, repo := range allRepos {
		fmt.Println(repo)
	}
	fmt.Println(len(allRepos))

	// For each repo, find its activity over the DAYS
	// Potential APIs:
	// client.Repositories.ListCommitActivity()
	// client.PullRequests.List()
	// client.Issues.List()
	// client.Activity.ListRepositoryEvents()
	// client.Activity.ListRepositoryNotifications()
}
