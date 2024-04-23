package main

import "github.com/google/go-github/v61/github"

func AnalyzeUserActivity(client *github.Client, argopts ArgOpts) {
	// TODO: Use one of these APIs
	// client.Activity.ListEventsPerformedByUser()
	// client.Activity.ListEventsReceivedByUser()
	// client.Activity.ListUserEventsForOrganization()
	// client.Activity.ListEventsForOrganization()
	// client.Activity.ListFeeds()
}
