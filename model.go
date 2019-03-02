package main

import (
	"log"
	"os"
	"sync"
	"time"
)

type QueueEntry struct {
	CSid      string
	Name      string
	TaskInfo  string
	JoinedAt  time.Time
	WasServed bool
	ServedAt  time.Time
}

type Queue struct {
	Mutex   sync.Mutex
	Entries []QueueEntry
}

var queue = Queue{Entries: []QueueEntry{}}

// Returns the estimated wait time in seconds.
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

// Deletes the persistence file and starts the queue from scratch. To be used by TAs only.
func NukeAllTheThings(ip string) {
	log.Println("User @", ip, "asked for database deletion. Will do.")
	queue = Queue{Entries: []QueueEntry{}}
	err := os.Remove("persistence.json")
	if err != nil {
		log.Println("Unable to delete persistence.json", err)
	}
}