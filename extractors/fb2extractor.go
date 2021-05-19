package extractors

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

type Book struct {
	Description  string `xml:"description"`
	Body string `xml:"body"`
	Binary []string `xml:"binary"`
}

func parseFb2File(fileName string)  {
	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	decoder := xml.NewDecoder(f)
	books := make([]Book, 0)

	var stack []string
	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			break
		}else if err != nil {
			fmt.Fprintf(os.Stderr, "xmlselect: %v\n", err)
			os.Exit(1)
		}
		switch tok := tok.(type) {
		case xml.StartElement:
			stack = append(stack, tok.Name.Local)
		case xml.EndElement:
			stack = stack[:len(stack)-1]
		case xml.CharData:
			if len(stack) > 0 {
				message := stack[len(stack)-1:]
				fmt.Printf("%s: %s\n", message, tok)
			}
		}
	}
	fmt.Println(books)
}
