package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v71/github"
)

func (a *Analyzer) AnalyzeUser() error {
	allUserEvents := []*github.Event{}
	page := 1
	for {
		opts := github.ListOptions{PerPage: 100, Page: page}
		events, res, err := a.client.Activity.ListEventsPerformedByUser(context.Background(), a.config.GitHubUser, false, &opts)
		if err != nil {
			return fmt.Errorf("failed to list user events: %w", err)
		}
		allUserEvents = append(allUserEvents, events...)
		if res.NextPage == 0 {
			break
		}
		page++
	}

	recentUserEvents := []*github.Event{}
	for _, event := range allUserEvents {
		if isWithinTimeRange(event, a.config.Days) {
			recentUserEvents = append(recentUserEvents, event)
		}
	}

	type EventCount struct {
		RepoCount          map[string]int
		RepoEventTypeCount map[string]map[string]int
	}
	eventCounts := EventCount{
		RepoCount:          make(map[string]int),
		RepoEventTypeCount: make(map[string]map[string]int),
	}

	for _, event := range recentUserEvents {
		eventCounts.RepoCount[event.GetRepo().GetName()]++
		if _, ok := eventCounts.RepoEventTypeCount[event.GetRepo().GetName()]; !ok {
			eventCounts.RepoEventTypeCount[event.GetRepo().GetName()] = make(map[string]int)
		}
		eventCounts.RepoEventTypeCount[event.GetRepo().GetName()][event.GetType()]++
	}

	fmt.Printf("\n\x1b[1;36m## User Activity Summary\x1b[0m\n\n")
	repoCountSortedSlice := sortMap(eventCounts.RepoCount)
	for _, pair := range repoCountSortedSlice {
		fmt.Printf("\x1b[1;33m### %s\x1b[0m \x1b[1;37mTotalEvents=\x1b[1;32m%d\x1b[0m\n\n",
			pair.Key, pair.Value)

		fmt.Printf("\x1b[1;37mEvent Types:\x1b[0m\n")
		repoEventCountSortedSlice := sortMap(eventCounts.RepoEventTypeCount[pair.Key])
		for _, pair := range repoEventCountSortedSlice {
			fmt.Printf("- \x1b[1;34m%s\x1b[0m: \x1b[1;32m%d\x1b[0m\n", pair.Key, pair.Value)
		}
		fmt.Printf("\n")
	}

	return nil
}
