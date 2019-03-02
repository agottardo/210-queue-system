package main

type HomePageValues struct {
	Error string
}

type JoinedPageValues struct {
	AheadOfMe         string
	HasEstimate       bool
	EstimatedWaitTime string
	JoinedAt          string
	Name              string
}

type StatusPageValues struct {
	Entries []QueueEntry
}
