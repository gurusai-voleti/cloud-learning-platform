package chunking

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"document-processor/models"

	"cloud.google.com/go/storage"
)

const (
	DEFAULT_CHUNK_SIZE = 250
)

var (
	textDocExtensions   = []string{"txt", "md", "rst"}
	officeDocExtensions = []string{"doc", "docx", "ppt", "pptx", "xls", "xlsx"}
)

type DocumentChunker struct {
	ChunkSize int
	ctx       context.Context
}

func NewDocumentChunker(ctx context.Context) *DocumentChunker {
	return &DocumentChunker{
		ChunkSize: DEFAULT_CHUNK_SIZE,
		ctx:       ctx,
	}
}

// ProcessDocument downloads and chunks a single document
func (c *DocumentChunker) ProcessDocument(doc models.DataSourceFile, tempDir string) ([]string, error) {
	// Download document
	localPath := filepath.Join(tempDir, doc.Filename)
	if err := c.downloadFromGCS(doc.GCSPath, localPath); err != nil {
		return nil, fmt.Errorf("download failed: %v", err)
	}
	defer os.Remove(localPath)

	// Read document content based on type
	docTextList, err := c.readDocument(doc.Filename, localPath)
	if err != nil {
		return nil, fmt.Errorf("reading failed: %v", err)
	}

	// Clean and join text
	var cleanedText strings.Builder
	for _, text := range docTextList {
		cleanedText.WriteString(CleanText(text))
		cleanedText.WriteString("\n")
	}

	// Chunk the text
	chunks := c.ChunkText(cleanedText.String())

	return chunks, nil
}

func (c *DocumentChunker) ChunkText(text string) []string {
	sentences := splitIntoSentences(text)
	var chunks []string
	var currentChunk strings.Builder
	currentTokens := 0

	for _, sentence := range sentences {
		sentenceTokens := len(strings.Fields(sentence))

		// If adding this sentence would exceed chunk size, start new chunk
		if currentTokens+sentenceTokens > c.ChunkSize && currentChunk.Len() > 0 {
			chunks = append(chunks, strings.TrimSpace(currentChunk.String()))
			currentChunk.Reset()
			currentTokens = 0
		}

		currentChunk.WriteString(sentence)
		currentChunk.WriteString(" ")
		currentTokens += sentenceTokens
	}

	// Add final chunk if not empty
	if currentChunk.Len() > 0 {
		chunks = append(chunks, strings.TrimSpace(currentChunk.String()))
	}

	return chunks
}
func (c *DocumentChunker) readDocument(filename, filepath string) ([]string, error) {
	ext := strings.ToLower(path.Ext(filename))

	// Handle text files
	for _, textExt := range textDocExtensions {
		if ext == "."+textExt {
			content, err := ioutil.ReadFile(filepath)
			if err != nil {
				return nil, err
			}
			return []string{string(content)}, nil
		}
	}

	// Handle PDF files
	if ext == ".pdf" {
		return ReadPDF(filepath)
	}

	// Handle Office documents
	for _, officeExt := range officeDocExtensions {
		if ext == "."+officeExt {
			return ReadOfficeDoc(filepath)
		}
	}

	return nil, fmt.Errorf("unsupported file type: %s", ext)
}

func (c *DocumentChunker) downloadFromGCS(gcsPath, localPath string) error {
	// Parse bucket and object from gcs://bucket/path format
	parts := strings.SplitN(strings.TrimPrefix(gcsPath, "gs://"), "/", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid GCS path format: %s", gcsPath)
	}
	bucketName, objectPath := parts[0], parts[1]

	// Create storage client
	client, err := storage.NewClient(c.ctx)
	if err != nil {
		return fmt.Errorf("failed to create storage client: %v", err)
	}
	defer client.Close()

	// Download object
	bucket := client.Bucket(bucketName)
	obj := bucket.Object(objectPath)

	reader, err := obj.NewReader(c.ctx)
	if err != nil {
		return fmt.Errorf("failed to create object reader: %v", err)
	}
	defer reader.Close()

	// Create local file
	f, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %v", err)
	}
	defer f.Close()

	// Copy content
	if _, err := io.Copy(f, reader); err != nil {
		return fmt.Errorf("failed to download file: %v", err)
	}

	return nil
}

func splitIntoSentences(text string) []string {
	// Simple sentence splitting - can be made more sophisticated
	sentences := strings.FieldsFunc(text, func(r rune) bool {
		return r == '.' || r == '!' || r == '?'
	})

	// Clean each sentence
	var cleanSentences []string
	for _, sentence := range sentences {
		sentence = strings.TrimSpace(sentence)
		if sentence != "" {
			cleanSentences = append(cleanSentences, sentence+".")
		}
	}

	return cleanSentences
}

func CleanText(text string) string {
	// Remove excessive whitespace
	text = strings.Join(strings.Fields(text), " ")

	// Remove common PDF artifacts
	text = strings.ReplaceAll(text, "\x00", "")

	// Remove other common artifacts
	artifacts := []string{
		"",  // Common encoding artifact
		"•", // Bullet points
		"©", // Copyright symbol
	}
	for _, artifact := range artifacts {
		text = strings.ReplaceAll(text, artifact, "")
	}

	return strings.TrimSpace(text)
}
