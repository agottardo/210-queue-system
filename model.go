package main

import (
	"log"
	"os"
	"sync"
	"time"
)

// This file contains the logic for our application.
// See persistence.go for disk (JSON) storage.

// A ticket in the Queue. A user can have multiple tickets, but only one
// has WasServed set to true. We store all tickets so that we can compute
// statistics if necessary by parsing the persistence.json file.
type QueueEntry struct {
	CSid      string
	Name      string
	TaskInfo  string
	JoinedAt  time.Time
	WasServed bool
	ServedAt  time.Time
}

type Queue struct {
	Mutex   sync.Mutex   // To handle concurrency, prevents multiple users from touching the DS.
	Entries []QueueEntry // Contains the actual tickets.
}

// Main in-memory data structure.
var queue = Queue{Entries: []QueueEntry{}}

// Adds the student with name and CSid to the queue.
// Returns how many students are ahead of the new student in the queue,
// and the estimated wait time in seconds.
func JoinQueue(name string, CSid string, taskInfo string) (int, int) {
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

// Returns true if the user with given CSid has joined the queue
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

// Marks the student with given CSid as served.
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

// Returns all tickets that have not been served yet.
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

// Returns the number of times the given CSid was helped in the last 12 hours.
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

// Returns the estimated wait time in seconds for a students that joins the
// queue right now.
func EstimatedWaitTime() float64 {
	var acc time.Duration
	count := 0
	queue.Mutex.Lock()
	for _, entry := range queue.Entries {
		if entry.WasServed {
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

// Deletes the persistence file and starts the queue from scratch. To be used by TAs only.
func NukeAllTheThings(ip string) {
	log.Println("User @", ip, "asked for database deletion. Will do.")
	queue = Queue{Entries: []QueueEntry{}}
	err := os.Remove("persistence.json")
	if err != nil {
		log.Println("Unable to delete persistence.json", err)
	}
}
