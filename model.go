package main

import (
	"log"
	"os"
	"sync"
	"time"
)

// This file contains the logic for our application.
// See persistence.go for disk (JSON) storage.

// A QueueEntry represents a ticket in the Queue. A user can have multiple
// tickets, but only one has WasServed set to true. We store all tickets
// so that we can compute statistics later by parsing the persistence.json file.
type QueueEntry struct {
	CSid      string
	Name      string
	TaskInfo  string
	JoinedAt  time.Time
	WasServed bool
	ServedAt  time.Time
}

// MaxNumTimesHelped is a constant which represents the maximum number of
// times a student can seek help within a 24 hour timeframe.
const MaxNumTimesHelped = 5

// Queue is the underlying thread-safe data structure (mutex + queue).
type Queue struct {
	Mutex   sync.Mutex   // To handle concurrency, prevents multiple users from touching the DS.
	Entries []QueueEntry // Contains the actual tickets.
}

// Main in-memory data structure.
var queue = Queue{Entries: []QueueEntry{}}

// JoinQueue adds the student with name and CSid to the queue.
// Returns how many students are ahead of the new student in the queue,
// and the estimated wait time in seconds.
// If the student has requested help more than MaxNumTimesHelped,
// returns how many times the students has asked for help already, and -1
func JoinQueue(name string, CSid string, taskInfo string) (int, int) {
	timesHelped := NumTimesHelped(CSid)
	if timesHelped < MaxNumTimesHelped {
		entry := QueueEntry{CSid, name, taskInfo, time.Now(), false, time.Now()}
		queue.Mutex.Lock()
		// How many un-served students joined before me?
		rsf := 0
		for _, entry := range queue.Entries {
			if !entry.WasServed {
				rsf++
			}
		}
		queue.Entries = append(queue.Entries, entry)
		UpdateDiskCopy()
		queue.Mutex.Unlock()
		return rsf, int(EstimatedWaitTime())
	}
	return timesHelped, -1
}

// HasJoinedQueue returns true if the user with given CSid has joined the queue
// and has not been served yet.
func HasJoinedQueue(CSid string) bool {
	queue.Mutex.Lock()
	for _, entry := range queue.Entries {
		if entry.CSid == CSid && !entry.WasServed {
			queue.Mutex.Unlock()
			return true
		}
	}
	queue.Mutex.Unlock()
	return false
}

// ServeStudent marks the student with given CSid as served.
func ServeStudent(CSid string) {
	queue.Mutex.Lock()
	for i, entry := range queue.Entries {
		if entry.CSid == CSid {
			if !entry.WasServed {
				queue.Entries[i].ServedAt = time.Now()
			}
			queue.Entries[i].WasServed = true
		}
	}
	UpdateDiskCopy()
	queue.Mutex.Unlock()
}

// UnservedEntries returns all tickets that have not been served yet.
func UnservedEntries() []QueueEntry {
	acc := []QueueEntry{}
	queue.Mutex.Lock()
	for _, entry := range queue.Entries {
		if !entry.WasServed {
			acc = append(acc, entry)
		}
	}
	queue.Mutex.Unlock()
	return acc
}

// NumTimesHelped returns the number of times the given CSid was helped in the last 24 hours.
func NumTimesHelped(CSid string) int {
	acc := 0
	queue.Mutex.Lock()
	for _, entry := range queue.Entries {
		if entry.CSid == CSid && entry.WasServed && entry.ServedAt.After(time.Now().AddDate(0, 0, -1)) {
			acc++
		}
	}
	queue.Mutex.Unlock()
	return acc
}

// EstimatedWaitTime returns the estimated wait time in seconds for a students that joins the
// queue right now, based on served entries from the past 30 minutes.
func EstimatedWaitTime() float64 {
	thirtyMinsAgo := time.Now().Add(-30 * time.Minute)
	var acc time.Duration
	count := 0
	queue.Mutex.Lock()
	for _, entry := range queue.Entries {
		if entry.WasServed && entry.ServedAt.After(thirtyMinsAgo) {
			waitTime := entry.ServedAt.Sub(entry.JoinedAt)
			acc += waitTime
			count++
		}
	}
	queue.Mutex.Unlock()
	if count == 0 {
		// Nobody was served yet, no estimate available.
		return 0
	}
	return acc.Seconds() / float64(count)
}

// NukeAllTheThings deletes the persistence file and starts the queue from scratch.
// To be used by TAs only.
func NukeAllTheThings(ip string) {
	log.Println("User @", ip, "asked for database deletion. Will do.")
	queue = Queue{Entries: []QueueEntry{}}
	err := os.Remove("persistence.json")
	if err != nil {
		log.Println("Unable to delete persistence.json", err)
	}
}

// QueuePositionForCSID returns whether the given CSid is waiting in the queue,
// and their position.
func QueuePositionForCSID(CSid string) (bool, uint) {
	entries := UnservedEntries()
	var acc uint = 0
	for _, entry := range entries {
		if entry.CSid == CSid {
			return true, acc
		}
		acc++
	}
	return false, 0
}
