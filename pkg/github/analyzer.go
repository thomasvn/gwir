package github

import (
	"fmt"

	"github.com/google/go-github/v71/github"
)

type Analyzer struct {
	client *github.Client
	config Flags
}

func NewAnalyzer(config Flags) *Analyzer {
	client := github.NewClient(nil)
	if config.GitHubPersonalAccessToken != "" {
		client = client.WithAuthToken(config.GitHubPersonalAccessToken)
	}
	return &Analyzer{client: client, config: config}
}

func (a *Analyzer) Analyze() error {
	if a.config.GitHubOrganization == "" && a.config.GitHubUser == "" {
		return fmt.Errorf("must specify either organization or user to analyze")
	}
	if a.config.GitHubOrganization != "" && a.config.GitHubUser != "" {
		return fmt.Errorf("cannot analyze both organization and user simultaneously")
	}
	if a.config.GitHubOrganization != "" {
		return a.AnalyzeOrg()
	}
	return a.AnalyzeUser()
}
