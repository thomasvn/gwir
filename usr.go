package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v61/github"
)

func AnalyzeUserActivity(client *github.Client, argopts ArgOpts) {
	// Paginated query to get all events for the user
	allUserEvents := []*github.Event{}
	page := 1
	for {
		opts := github.ListOptions{PerPage: 100, Page: page}
		events, res, _ := client.Activity.ListEventsPerformedByUser(context.Background(), argopts.GITHUB_USER, false, &opts)
		allUserEvents = append(allUserEvents, events...)
		// TODO: filter for timestamps here to reduce the # of API requests made
		if res.NextPage == 0 {
			break
		}
		page++
	}

	// Filter for events within the timestamp
	recentUserEvents := []*github.Event{}
	for _, event := range allUserEvents {
		if event.GetCreatedAt().After(time.Now().Add(time.Duration(-argopts.DAYS) * time.Hour * 24)) {
			recentUserEvents = append(recentUserEvents, event)
		}
	}

	// Process the event types
	type EventCount struct {
		RepoCount          map[string]int
		RepoEventTypeCount map[string]map[string]int
	}
	eventCounts := EventCount{RepoCount: make(map[string]int), RepoEventTypeCount: make(map[string]map[string]int)}
	for _, event := range recentUserEvents {
		eventCounts.RepoCount[event.GetRepo().GetName()]++
		if _, ok := eventCounts.RepoEventTypeCount[event.GetRepo().GetName()]; !ok {
			eventCounts.RepoEventTypeCount[event.GetRepo().GetName()] = make(map[string]int)
		}
		eventCounts.RepoEventTypeCount[event.GetRepo().GetName()][event.GetType()]++
	}

	// Print sorted results
	repoCountSortedSlice := sortMap(eventCounts.RepoCount)
	for _, pair := range repoCountSortedSlice {
		fmt.Printf("\n%s: %d\n", pair.Key, pair.Value)
		repoEventCountSortedSlice := sortMap(eventCounts.RepoEventTypeCount[pair.Key])
		for _, pair := range repoEventCountSortedSlice {
			fmt.Printf("  - %s: %d\n", pair.Key, pair.Value)
		}
	}

	// TODO: Depending on the event type, print details and provide a link?
	// TODO: Provide links to PRs/Issues/Repos?
}
