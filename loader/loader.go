package loader

import (
	"fmt"
	"github.com/DFilyushin/gobooklibrary/book"
	"github.com/DFilyushin/gobooklibrary/database"
	"github.com/DFilyushin/gobooklibrary/extractors"
	"github.com/DFilyushin/gobooklibrary/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/fs"
	"log"
	"path/filepath"
	"sync"
	"time"
)

const IndexFileExtension = ".inp"

const maxWorkers = 30000

type AuthorCache map[book.Author]primitive.ObjectID

var authorCache = make(AuthorCache)
var mu sync.RWMutex

func checkItemInArray(item string, arr *[]string) bool {
	for _, arrItem := range *arr {
		if arrItem == item {
			return true
		}
	}
	return false
}

func getAuthors(authors []book.Author) []primitive.ObjectID {
	/*
		Get authors by list
	*/
	var result []primitive.ObjectID
	var objectId primitive.ObjectID

	for _, author := range authors {
		mu.Lock()
		cacheValue, mOK := authorCache[author]
		if mOK {
			result = append(result, cacheValue)
			mu.Unlock()
			continue

		}
		mu.Unlock()

		mongoAuthor, err := database.GetAuthorByFullNames(author.LastName, author.FirstName, author.MiddleName)
		if err != nil {
			log.Printf("Error request author %s. Error message: %v\n", author.LastName, err)
			continue
		}

		if mongoAuthor != nil {
			mu.Lock()
			authorCache[author] = mongoAuthor.ID
			mu.Unlock()
			objectId = mongoAuthor.ID
		} else {
			newAuthor := &models.AuthorModel{
				ID:         primitive.NewObjectID(),
				LastName:   author.LastName,
				FirstName:  author.FirstName,
				MiddleName: author.MiddleName,
			}
			value, err := database.AddAuthor(newAuthor)
			if err != nil {
				log.Printf("Error adding new author %s. Error message: %v\n", author.LastName, err)
				continue
			} else {
				mu.Lock()
				authorCache[author] = *value
				mu.Unlock()
				objectId = *value
			}
		}

		result = append(result, objectId)
	}
	return result
}

func checkBookExists(bookId string) (bool, error) {
	/*
		Check book exists
	*/
	mongoBook, err := database.GetBookById(bookId)
	if err != nil {
		return false, fmt.Errorf("Error on book %s. Error message: %s\n", bookId, err)
	}
	return mongoBook != nil, nil
}

func ProcessBook(book book.Book) (*models.BookModel, error) {
	/*
		Checking existing book, adding
	*/
	isExists, err := checkBookExists(book.BookId)
	if err != nil {
		return nil, err
	}
	if isExists {
		return nil, nil
	}

	authors := getAuthors(book.Authors)
	isDeleted := book.Deleted == "1"

	newBook := &models.BookModel{
		ID:        primitive.NewObjectID(),
		BookId:    book.BookId,
		BookName:  book.BookName,
		ISBN:      book.ISBN,
		Deleted:   isDeleted,
		Added:     book.Added,
		Authors:   authors,
		Language:  book.Language,
		Year:      book.Year,
		Rate:      book.Rating,
		City:      book.City,
		Filename:  book.Filename,
		Keywords:  book.Keywords,
		Genres:    book.Genres,
		PubName:   book.PubName,
		Publisher: book.Publisher,
		Series:    book.Series,
		SerialNum: book.SerialNum,
	}
	//_, err = database.AddBook(newBook)
	return newBook, err
}

func processBooks(books *[]book.Book) (int, int) {
	/*
		Process list of books with error handling
	*/
	//var ctx = context.TODO()
	//var sem = semaphore.NewWeighted(maxWorkers)
	var wg sync.WaitGroup
	bChannel := make(chan models.BookModel, 1000)
	//booksBatch := make([]models.BookModel, 0)
	countBooks := len(*books)
	countError := 0


	for _, bookItem := range *books {
		wg.Add(1)
		//if err := sem.Acquire(ctx, 1); err != nil {
		//	log.Printf("Failed to acquire semaphore: %v", err)
		//	break
		//}

		go func(item book.Book) {
			//defer sem.Release(1)
			defer wg.Done()
			newBook, err := ProcessBook(item)
			if err != nil {
				fmt.Printf("Error adding book %s. Error message: %s\n", item.BookId, err)
			}
			if newBook != nil {
				bChannel <- *newBook
			}
		}(bookItem)
	}

	go func() {
		wg.Wait()
		close(bChannel)
	}()


	for item := range bChannel {
		fmt.Println(item)
		//booksBatch = append(booksBatch, item)
		//if len(booksBatch) == 100 {
		//	//database.AddBooks()
		//	fmt.Println(booksBatch)
		//	fmt.Println("100")
		//	booksBatch = booksBatch[:0]
		//}
	}

	//if err := sem.Acquire(ctx, int64(maxWorkers)); err != nil {
	//	log.Printf("Failed to acquire semaphore: %v", err)
	//}
	return countBooks - countError, countError
}

func ProcessIndexFile(fileName string) {
	books, err := extractors.ProcessIndexFile(fileName)
	if err != nil {
		fmt.Printf("Error processing file %s. Error message: %v\n", fileName, err)
	} else {
		processBooks(&books)
	}
}

func getFileFromPath(path string, ignoreFiles []string) ([]string, error) {
	/*
		Walk directory, find .inp files for processing
	*/
	files := make([]string, 0)

	err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == IndexFileExtension {
			fileName := filepath.Base(path)
			if !checkItemInArray(fileName, &ignoreFiles) {
				files = append(files, path)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return files, nil

}

func ProcessIndexPath(path string, ignoreFiles []string) []string {
	/*
		Walk path directory, find .inp files for processing
	*/
	processedFiles := make([]string, 0)
	files, err := getFileFromPath(path, ignoreFiles)
	if err != nil {
		fmt.Printf("Prevent panic by handling failure accessing a path %q: %v\n", path, err)
		return nil
	}

	books := make([]book.Book, 0)
	for _, fileName := range files {
		books, err = extractors.ProcessIndexFile(fileName)
		if err != nil {
			fmt.Printf("Error processing file %s. Error message: %v\n", path, err)
		} else {
			start := time.Now()
			processBooks(&books)
			duration := time.Since(start)
			log.Printf("Duration loading file %s: %v", fileName, duration)
			processedFiles = append(processedFiles, fileName)
		}
	}
	return processedFiles
}
