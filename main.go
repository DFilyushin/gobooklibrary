package main

import (
	"bufio"
	"fmt"
	"github.com/DFilyushin/gobooklibrary/database"
	"github.com/DFilyushin/gobooklibrary/loader"
	"log"
	"os"
	"time"
)
const indexFilesPath = "c://var//library//index//"
const indexFileName = "c://var//library//index//fb2-060424-074391.inp"
const processedFileName = "processed.log"
const processIgnoreFile = true
const mongoConnectionString = "mongodb://root:qwerty@localhost:27017/"
const mongoDatabase = "books_work"

func saveProcessedFiles(fileName string, files []string) error {
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for _, line := range files {
		fmt.Fprintln(w, line)
	}
	w.Flush()
	return nil
}

func getProcessedFiles(fileName string) ([]string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	result := make([]string, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}
	return result, scanner.Err()
}

func loadIndexFiles(ignoreEarlyProcessedFiles bool) {
	var ignoredFiles []string
	var err error

	if ignoreEarlyProcessedFiles {
		ignoredFiles, err = getProcessedFiles(processedFileName)
		if err != nil {
			fmt.Printf("Error reading ignore file! Scan all files. Error: %v\n", err)
		}
	}

	processedFiles := loader.ProcessIndexPath(indexFilesPath, ignoredFiles)

	for _, fileName := range processedFiles {
		fmt.Printf("%s Ok\n", fileName)
	}

	if ignoreEarlyProcessedFiles {
		err = saveProcessedFiles(processedFileName, processedFiles)
		if err != nil {
			fmt.Printf("Error save processed files %s. Error: %v\n", processedFileName, err)
		}
	}
}

func loadIndexFile(fileName string)  {
	loader.ProcessIndexFile(fileName)
}

func main()  {
	log.Println("Loader started...")

	database.SetupDatabase(mongoConnectionString, mongoDatabase)

	start := time.Now()
	loadIndexFiles(processIgnoreFile)
	//loadIndexFile(indexFileName)
	duration := time.Since(start)

	log.Println("Duration loading:", duration)
}