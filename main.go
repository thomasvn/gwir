package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/google/go-github/v61/github"
)

type ArgOpts struct {
	DAYS                         int
	TOPXACTIVITIES               int
	GITHUB_ORGANIZATION          string
	GITHUB_USER                  string
	GITHUB_PERSONAL_ACCESS_TOKEN string
}

func main() {
	// Parse arguments
	opts := ArgOpts{}
	flag.IntVar(&opts.DAYS, "days", 7, "How many days back to analyze")
	flag.IntVar(&opts.TOPXACTIVITIES, "top", 5, "How many top PRs/Issues to show")
	flag.StringVar(&opts.GITHUB_ORGANIZATION, "org", "", "GitHub organization to analyze")
	flag.StringVar(&opts.GITHUB_USER, "usr", "", "GitHub user to analyze")
	flag.StringVar(&opts.GITHUB_PERSONAL_ACCESS_TOKEN, "token", "", "Optional. Passing a GitHub Personal Access Token allows you to view private repositories and make more API requests per hour. You can also set this token as an environment variable GITHUB_PERSONAL_ACCESS_TOKEN.")
	flag.Parse()

	// Prioritize any Github PAT which was passed as a flag, then check env vars
	var client *github.Client
	if opts.GITHUB_PERSONAL_ACCESS_TOKEN != "" {
		client = github.NewClient(nil).WithAuthToken(opts.GITHUB_PERSONAL_ACCESS_TOKEN)
	} else if os.Getenv("GITHUB_PERSONAL_ACCESS_TOKEN") != "" {
		client = github.NewClient(nil).WithAuthToken(os.Getenv("GITHUB_PERSONAL_ACCESS_TOKEN"))
	} else {
		client = github.NewClient(nil)
	}

	if opts.GITHUB_ORGANIZATION == "" && opts.GITHUB_USER == "" {
		fmt.Println("Either -org or -usr must be set. Use -h to view all options.")
		return
	} else if opts.GITHUB_ORGANIZATION != "" && opts.GITHUB_USER != "" {
		fmt.Println("Cannot analyze both -org and -user. Only set one.")
		return
	} else if opts.GITHUB_ORGANIZATION != "" {
		AnalyzeOrgActivity(client, opts)
	} else if opts.GITHUB_USER != "" {
		AnalyzeUserActivity(client, opts)
	}
}
