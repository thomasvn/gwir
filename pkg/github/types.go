package github

import (
	"sort"
	"time"

	"github.com/google/go-github/v71/github"
)

type Flags struct {
	Days                      int
	TopXActivities            int
	GitHubOrganization        string
	GitHubUser                string
	GitHubPersonalAccessToken string
}

type pair struct {
	Key   string
	Value int
}

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

func isWithinTimeRange(event *github.Event, days int) bool {
	return event.GetCreatedAt().After(time.Now().Add(time.Duration(-days) * time.Hour * 24))
}
