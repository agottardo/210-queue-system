package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// LoadDataFromDisk fills the in-memory data structure by reading
// the persistence.json file. To be called when booting/restarting
// the application.
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

// UpdateDiskCopy converts the data structure to a JSON file
// and saves it to disk. The caller of this function should
// have locked the mutex before calling it.
func UpdateDiskCopy() {
	queueJSON, _ := json.Marshal(queue)
	err := ioutil.WriteFile("persistence.json", queueJSON, 0644)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Updated persistence.json with new data.")
}

func LoadPasswordsFromDisk() map[string]string {
	passwordsStore, err := ioutil.ReadFile("authdb.json")
	if err != nil {
		log.Fatalln("Couldn't read authdb.json. Create it before running this application.")
	}
	theMap := map[string]string{}
	jsonerr := json.Unmarshal(passwordsStore, &theMap)
	if jsonerr != nil {
		log.Fatalln("I couldn't unmarshall authdb.json. The JSON syntax is probably bad.")
	}
	return theMap
}
