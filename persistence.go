package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Handles persistence to local .json file.
func LoadDataFromDisk() {
	// Restore data from disk storage before running (persistence.json).
	jsonStore, err := ioutil.ReadFile("persistence.json")
	if err != nil {
		log.Println("Couldn't read persistence.json. Perhaps, this is the first time the application is running?")
	} else {
		jsonerr := json.Unmarshal(jsonStore, &queue)
		if jsonerr != nil {
			log.Println("Couldn't unmarshal persistence.json. This is likely the result of data corruption.", jsonerr)
			log.Println("210queue is starting with a fresh new datastore.")
		} else {
			log.Println("Restarting with data from persistence.json.")
			queue.Mutex.Lock()
			numEntries := len(queue.Entries)
			queue.Mutex.Unlock()
			log.Println("Restarting with", numEntries, "entries in the queue.")
		}
	}
}

func UpdateDiskCopy() {
	queueJson, _ := json.Marshal(queue)
	err := ioutil.WriteFile("persistence.json", queueJson, 0644)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Updated persistence.json with new data.")
}
