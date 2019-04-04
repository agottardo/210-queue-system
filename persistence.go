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
		jsonErr := json.Unmarshal(jsonStore, &queue)
		if jsonErr != nil {
			log.Println("Couldn't unmarshal persistence.json. This is likely the result of data corruption.", jsonErr)
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
	jsonErr := json.Unmarshal(passwordsStore, &theMap)
	if jsonErr != nil {
		log.Fatalln("I couldn't unmarshal authdb.json. The JSON syntax is probably bad.")
	}
	return theMap
}

// A type that stores the application configuration.
type Config struct {
	ListenAt          string
	AuthSecret        string
	MaxNumTimesHelped uint
}

// Reads the system configuration from the config.json file.
func ReadConfig() Config {

	config := Config{}
	configStore, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalln("Couldn't read config.json. Create and fill it before running this application.")
	}
	theMap := map[string]interface{}{}
	jsonErr := json.Unmarshal(configStore, &theMap)
	if jsonErr != nil {
		log.Fatalln("I couldn't unmarshal config.json. The JSON syntax is probably bad.")
	}

	// ListenAt is the HTTP port the web-server should listen at for incoming
	// connections.
	config.ListenAt = theMap["ListenAt"].(string)

	// AuthSecret is a random string that is used when hashing CSids and
	// storing them in a cookie.
	// This kind of crypto is used to ensure that only whoever joined
	// the queue can actually leave it manually (/leaveearly).
	config.AuthSecret = theMap["AuthSecret"].(string)

	// MaxNumTimesHelped is a constant which represents the maximum number of
	// times a student can seek help within a 24 hour timeframe.
	config.MaxNumTimesHelped = uint(theMap["MaxNumTimesHelped"].(float64))

	return config
}
