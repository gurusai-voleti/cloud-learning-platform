package chunking

import (
	"fmt"
	"log"
	"path"
	"strings"

	"github.com/unidoc/unioffice/document"
	"github.com/unidoc/unioffice/presentation"
	"github.com/unidoc/unioffice/spreadsheet"
)

func ReadOfficeDoc(filepath string) ([]string, error) {
	ext := strings.ToLower(path.Ext(filepath))
	var docTextList []string

	switch ext {
	case ".docx":
		// Handle Word documents
		doc, err := document.Open(filepath)
		if err != nil {
			return nil, fmt.Errorf("error opening Word document: %v", err)
		}

		var text strings.Builder
		for _, para := range doc.Paragraphs() {
			// Extract text from each run in the paragraph
			for _, run := range para.Runs() {
				text.WriteString(run.Text())
			}
			text.WriteString("\n")
		}

		if text.Len() > 0 {
			docTextList = append(docTextList, text.String())
		}

	case ".pptx", ".ppt", ".pptm":
		// Handle PowerPoint documents
		pres, err := presentation.Open(filepath)
		if err != nil {
			return nil, fmt.Errorf("error opening PowerPoint document: %v", err)
		}

		var text strings.Builder
		for _, slide := range pres.Slides() {
			text.WriteString(slide.ExtractText().Text())
			text.WriteString("\n")
			text.WriteString("\n--- Slide Break ---\n")
		}

		if text.Len() > 0 {
			docTextList = append(docTextList, text.String())
		}

	case ".xlsx", ".xls":
		// Handle Excel documents
		xlsx, err := spreadsheet.Open(filepath)
		if err != nil {
			return nil, fmt.Errorf("error opening Excel document: %v", err)
		}

		var text strings.Builder
		for _, sheet := range xlsx.Sheets() {
			text.WriteString(fmt.Sprintf("--- Sheet: %s ---\n", sheet.Name()))

			for _, row := range sheet.Rows() {
				for _, cell := range row.Cells() {
					text.WriteString(cell.GetString())
					text.WriteString("\t")
				}
				text.WriteString("\n")
			}
			text.WriteString("\n")
		}

		if text.Len() > 0 {
			docTextList = append(docTextList, text.String())
		}

	default:
		return nil, fmt.Errorf("unsupported office document type: %s", ext)
	}

	if len(docTextList) == 0 {
		return nil, fmt.Errorf("no text content found in document")
	}

	log.Printf("Finished reading office document")
	return docTextList, nil
}
