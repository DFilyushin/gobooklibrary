package extractors

import (
	"bufio"
	"github.com/DFilyushin/gobooklibrary/book"
	"os"
	"strings"
)

const ItemSeparatorInLine = string(4)
const ArraySeparator = ":"
const IndexFileExtension = ".inp"
const InLineItemSeparator = ","

func getAuthorsByString(authors string) []book.Author {
	result := make([]book.Author, 0)
	var firstName, lastName, middleName string

	items := strings.Split(authors, ArraySeparator)
	items = items[:len(items)-1]
	for _, item := range items {
		authorNames := strings.Split(item, InLineItemSeparator)
		lastName, firstName, middleName = authorNames[0], authorNames[1], authorNames[2]
		author := &book.Author{LastName: lastName, FirstName: firstName, MiddleName: middleName}
		result = append(result, *author)
	}
	return result
}

func getKeywordsByString(keysLine string) []string {
	if len(keysLine) == 0 {
		return nil
	}
	delimiters := []string {",", " "}
	for _, delimiter := range delimiters {
		if strings.Contains(keysLine, delimiter) {
			return strings.Split(keysLine, delimiter)
		}
	}
	return nil
}

func processLine(line string) (*book.Book, error) {
	data := strings.Split(line, ItemSeparatorInLine)
	authorLine, genreLine, title, series, serNum, fileName, _, bookId, bookDeleted, _, bookAdded, language, rate, keywords :=
		data[0], data[1], data[2], data[3], data[4], data[5], data[6], data[7], data[8], data[9], data[10], data[11], data[12], data[13]

	authors := getAuthorsByString(authorLine)
	genres := strings.Split(genreLine, ArraySeparator)
	genres = genres[:len(genres)-1]

	return &book.Book{
		BookId:    bookId,
		Authors:   authors,
		Genres:    genres,
		BookName:  title,
		Series:    series,
		SerialNum: serNum,
		Filename:  fileName,
		Deleted:   bookDeleted,
		Added:     bookAdded,
		Language:  language,
		Keywords:  getKeywordsByString(keywords),
		Rating:    rate,
	}, nil
}

func ProcessIndexFile(fileName string) ([]book.Book, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := bufio.NewScanner(f)
	r.Split(bufio.ScanLines)

	var bookItem *book.Book
	var books = make([]book.Book, 0)
	lineCounter := 1
	for r.Scan() {
		bookItem, err = processLine(r.Text())
		if err != nil {
			return nil, err
		} else {
			books = append(books, *bookItem)
		}
		lineCounter ++
	}
	return books, nil
}