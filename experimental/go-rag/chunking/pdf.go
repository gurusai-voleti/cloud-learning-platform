package chunking

import (
	"fmt"
	"log"
	"os"

	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
)

func ReadPDF(filepath string) ([]string, error) {
	// Open the PDF file
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("error opening PDF file: %v", err)
	}
	defer file.Close()

	f, err := model.NewPdfReader(file)
	if err != nil {
		return nil, fmt.Errorf("error creating PDF reader: %v", err)
	}

	// Get total page count
	numPages, err := f.GetNumPages()
	if err != nil {
		return nil, fmt.Errorf("error getting page count: %v", err)
	}

	log.Printf("Reading PDF with %d pages", numPages)

	// Extract text from each page
	var docTextList []string
	for i := 1; i <= numPages; i++ {
		page, err := f.GetPage(i)
		if err != nil {
			log.Printf("Warning: error getting page %d: %v", i, err)
			continue
		}

		extract, err := extractor.New(page)
		if err != nil {
			log.Printf("Warning: error creating extractor for page %d: %v", i, err)
			continue
		}

		text, err := extract.ExtractText()
		if err != nil {
			log.Printf("Warning: error extracting text from page %d: %v", i, err)
			continue
		}

		if text != "" {
			docTextList = append(docTextList, text)
		}
	}

	if len(docTextList) == 0 {
		return nil, fmt.Errorf("no text content found in PDF")
	}

	log.Printf("Finished reading PDF file, extracted %d pages of text", len(docTextList))
	return docTextList, nil
}
