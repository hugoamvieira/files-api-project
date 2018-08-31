package main

import (
	"log"
	"os"
)

/*

We would also need to get some statistics per folder basis and retrieve them through another entry point. These statistics are:
Total number of files in that folder.
Average number of alphanumeric characters per text file (and standard deviation) in that folder.
Average word length (and standard deviation) in that folder.
Total number of bytes stored in that folder.
Note: All these computations must be calculated recursively from the provided path to the entry point.

*/

const (
	defaultAPIAddr = ":8000"
)

func main() {
	log.Println("Starting File API...")

	if _, err := os.Stat(RootPath); os.IsNotExist(err) {
		err := os.Mkdir(RootPath, DefaultFolderPermission)
		if err != nil {
			log.Fatalf("Failed to create root folder. Error: %v\n", err)
		}
		log.Println("Root folder did not exist - It has been created.")
	}

	StartAPI(defaultAPIAddr) // Start API
}
