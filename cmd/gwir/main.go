package main

import (
	"flag"
	"log"
	"os"

	"github.com/thomasvn/gwir/pkg/github"
)

func main() {
	flags := github.Flags{}
	flag.IntVar(&flags.Days, "days", 7, "How many days back to analyze")
	flag.IntVar(&flags.TopXActivities, "top", 5, "How many top PRs/Issues to show")
	flag.StringVar(&flags.GitHubOrganization, "org", "", "GitHub organization to analyze")
	flag.StringVar(&flags.GitHubUser, "usr", "", "GitHub user to analyze")
	flag.StringVar(
		&flags.GitHubPersonalAccessToken,
		"token",
		os.Getenv("GITHUB_PERSONAL_ACCESS_TOKEN"),
		"GitHub Personal Access Token (can also be set via GITHUB_PERSONAL_ACCESS_TOKEN env var)",
	)
	flag.Parse()

	if err := github.NewAnalyzer(flags).Analyze(); err != nil {
		log.Fatal(err)
	}
}
