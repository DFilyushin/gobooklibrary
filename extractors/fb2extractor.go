package extractors

import (
	"encoding/xml"
	"fmt"
	"github.com/DFilyushin/gobooklibrary/helpers"
	"golang.org/x/net/html/charset"
	"io"
	"os"
	"strings"
)

const CoverTag = "coverpage"
const ImageTag = "image"
const BinaryTag = "binary"

func ExtractImageFromFb2(fileName string) (string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return "", err
	}

	defer f.Close()

	decoder := xml.NewDecoder(f)
	decoder.CharsetReader = charset.NewReaderLabel

	var processed bool
	var processedBinary bool
	var tagValue string
	var coverFileName string

	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			return "", err
		}

		switch tok := tok.(type) {
		case xml.StartElement:
			if tok.Name.Local == CoverTag {
				processed = true
			}
			if tok.Name.Local == ImageTag && processed {
				coverFileName = tok.Attr[0].Value[1:]
			}
			if tok.Name.Local == BinaryTag {
				for _, item := range tok.Attr {
					if item.Name.Local == "id" && item.Value == coverFileName {
						processedBinary = true
					}
				}
			}
		case xml.EndElement:
			if tok.Name.Local == CoverTag {
				processed = false
			}
		case xml.CharData:
			if processedBinary {
				tagValue = string(tok)
				fmt.Println(tagValue)
				processedBinary = false
			}
		}
	}

	return "", nil
}

func ParseFb2File(fileName string, tags []string) (map[string]string, error) {
	result := make(map[string]string)

	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	decoder := xml.NewDecoder(f)
	decoder.CharsetReader = charset.NewReaderLabel

	var tagValues []string
	var tagOpened bool
	var tagValue string

	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		switch tok := tok.(type) {
		case xml.StartElement:
			if helpers.CheckItemInArray(tok.Name.Local, &tags) {
				tagValues = tagValues[:0] //clear storage
				tagOpened = true
			}
		case xml.EndElement:
			if helpers.CheckItemInArray(tok.Name.Local, &tags) {
				result[tok.Name.Local] = strings.Join(tagValues, " ")
				tagOpened = false
			}
		case xml.CharData:
			if tagOpened {
				tagValue = string(tok)
				tagValues = append(tagValues, strings.TrimSpace(tagValue))
			}
		}
	}
	return result, nil
}
