package main

import (
	"fmt"
	"github.com/DFilyushin/gobooklibrary/database"
	"github.com/DFilyushin/gobooklibrary/helpers"
	"github.com/DFilyushin/gobooklibrary/loader"
	"log"
	"time"
)

const indexFilesPath = "c://var//library//index//"
const processedFileName = "processed.log"
const processIgnoreFile = true
const mongoConnectionString = "mongodb://root:qwerty@localhost:27017/"
const mongoDatabase = "books_work"


func loadIndexFiles(ignoreEarlyProcessedFiles bool) {
	var ignoredFiles []string
	var err error

	if ignoreEarlyProcessedFiles {
		ignoredFiles, err = helpers.ReadTextFromFile(processedFileName)
		if err != nil {
			fmt.Printf("Error reading ignore file! Scan all files. Error: %v\n", err)
		}
	}

	processedFiles := loader.ProcessIndexPath(indexFilesPath, ignoredFiles)

	for _, fileName := range processedFiles {
		fmt.Printf("%s Ok\n", fileName)
	}

	if ignoreEarlyProcessedFiles {
		err = helpers.WriteTextToFile(processedFileName, processedFiles)
		if err != nil {
			fmt.Printf("Error save processed files %s. Error: %v\n", processedFileName, err)
		}
	}
}

func loadIndexFile(fileName string)  {
	loader.ProcessFbFile(fileName)
}

func main()  {
	log.Println("Loader started...")

	database.SetupDatabase(mongoConnectionString, mongoDatabase)

	start := time.Now()
	//loadIndexFiles(processIgnoreFile)
	loadIndexFile("C:\\var\\library\\data\\91846.fb2")
	duration := time.Since(start)

	log.Println("Duration loading:", duration)
}