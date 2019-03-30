package main

// HomePageValues represents the values used in the homepage.
type HomePageValues struct {
	CountHelped uint
	Error       string
}

// RejectedPageValues represents the values used in the queue rejected page
type RejectedPageValues struct {
	NumTimesJoined uint
	Name           string
}

// StatusPageValues represents the values used in the "current queue status" page.
type StatusPageValues struct {
	Entries []QueueEntry
}
