package main

// HomePageValues represents the values used in the homepage.
type HomePageValues struct {
	Error string
}

// JoinedPageValues represents the values used in the "queue joined" page.
type JoinedPageValues struct {
	AheadOfMe         string
	HasEstimate       bool
	EstimatedWaitTime string
	JoinedAt          string
	Name              string
}

// RejectedPageValues represents the values used in the queue rejected page
type RejectedPageValues struct {
	NumTimesJoined int
	Name           string
}

// StatusPageValues represents the values used in the "current queue status" page.
type StatusPageValues struct {
	Entries []QueueEntry
}
