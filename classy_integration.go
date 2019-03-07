package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

// This file contains functions used by the application to integrate with
// Classy, the course management system used at UBC Computer Science.

// For instance, it allows us to give priority in the queue to students
// who are registered in a specific lab section, based on their SSC
// registration information.

var students map[string]string // CSid -> LabSection

const CLASSY_ENDPOINT = "https://cs210.ugrad.cs.ubc.ca/portal/admin/students"
const CLASSY_USER = "queueapp"
const CLASSY_TOKEN = "surely-not-posting-this-on-github-dude"

// LoadClassyData connects to Classy over its REST endpoint, downloads
// and parses student information, then loads it into memory.
// Returns true if loading was successful, false otherwise.
func LoadClassyData() bool {
	students = map[string]string{}

	request, err := http.NewRequest("GET", CLASSY_ENDPOINT, nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("token", CLASSY_TOKEN)
	request.Header.Set("user", CLASSY_USER)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println("Failed to connect to Classy to retrieve registration info:", err)
		return false
	}
	if response.StatusCode != http.StatusOK {
		log.Println("Classy returned a non-200 status code while fetching registration info:", err)
		return false
	}

	// TODO: fill the `students` map here by parsing the json response
	_, _ = io.Copy(os.Stdout, response.Body)
	_ = response.Body.Close()
	return len(students) > 0
}

// LabSectionForStudent returns true if the given string contains
// the CS ID of a student that is registered for the course.
// If the student is registered, it also returns their lab section,
// otherwise it produces the empty string.
func LabSectionForStudent(CSid string) (isRegistered bool, labSection string) {
	labSection, isRegistered = students[CSid]
	return
}
