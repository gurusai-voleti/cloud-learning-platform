package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"cloud.google.com/go/firestore"

	"document-processor/chunking"
	"document-processor/embedding"
	"document-processor/models"
)

type documentChunker interface {
	ProcessDocument(doc models.DataSourceFile, tempDir string) ([]string, error)
}

type vertexEmbedder interface {
	GenerateEmbeddings(ctx context.Context, chunks []string) ([][]float32, error)
}

// DocumentProcessor orchestrates the document processing pipeline
type DocumentProcessor struct {
	ctx             context.Context
	firestoreClient *firestore.Client
	chunker         documentChunker
	embedder        vertexEmbedder
	projectID       string
	queryEngineID   string
	jobRef          *firestore.DocumentRef
}

// NewDocumentProcessor creates a new document processor instance
func NewDocumentProcessor(
	ctx context.Context,
	firestoreClient *firestore.Client,
	projectID string,
	queryEngineID string,
	jobRef *firestore.DocumentRef,
) (*DocumentProcessor, error) {

	// Initialize chunker
	chunker := chunking.NewDocumentChunker(ctx)

	// Initialize embedder
	embedder, err := embedding.NewVertexEmbedder(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize embedder: %v", err)
	}

	return &DocumentProcessor{
		ctx:             ctx,
		firestoreClient: firestoreClient,
		chunker:         chunker,
		embedder:        embedder,
		projectID:       projectID,
		queryEngineID:   queryEngineID,
		jobRef:          jobRef,
	}, nil
}

func main() {
	// Configure logger
	log.SetOutput(os.Stdout)

	// Get environment variables
	projectID := os.Getenv("GCP_PROJECT")
	if projectID == "" {
		log.Fatal("GCP_PROJECT environment variable not set")
	}

	jobID := os.Getenv("JOB_ID")
	if jobID == "" {
		log.Fatal("JOB_ID environment variable not set")
	}
	log.Printf("Processing job ID: %s", jobID)

	// Initialize context
	ctx := context.Background()

	// Initialize Firestore client
	firestoreClient, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer firestoreClient.Close()

	// Get job document
	docRef := firestoreClient.Collection("batch_jobs").Doc(jobID)
	jobDoc, err := docRef.Get(ctx)
	if err != nil {
		updateJobError(ctx, docRef, fmt.Errorf("failed to get job document: %v", err))
		log.Fatal(err)
	}

	// Parse job input
	var jobInput models.JobInput
	inputData, ok := jobDoc.Data()["input_data"].(string)
	if !ok {
		err := fmt.Errorf("failed to get input_data as string from job document")
		updateJobError(ctx, docRef, err)
		log.Fatal(err)
	}
	if err := json.Unmarshal([]byte(inputData), &jobInput); err != nil {
		updateJobError(ctx, docRef, fmt.Errorf("failed to decode job input: %v", err))
		log.Fatal(err)
	}

	// Initialize processor
	processor, err := NewDocumentProcessor(
		ctx,
		firestoreClient,
		projectID,
		jobInput.QueryEngineID,
		docRef,
	)
	if err != nil {
		updateJobError(ctx, docRef, fmt.Errorf("failed to initialize processor: %v", err))
		log.Fatal(err)
	}

	// Process documents
	results, err := processor.ProcessDocuments(jobInput.Documents)
	if err != nil {
		updateJobError(ctx, docRef, fmt.Errorf("failed to process documents: %v", err))
		log.Fatal(err)
	}

	// Update job with results
	resultData := map[string]interface{}{
		"processed_documents": results,
	}

	_, err = docRef.Update(ctx, []firestore.Update{
		{Path: "result_data", Value: resultData},
		{Path: "status", Value: "succeeded"},
		{Path: "completed_at", Value: time.Now()},
	})
	if err != nil {
		updateJobError(ctx, docRef, fmt.Errorf("failed to update job results: %v", err))
		log.Fatal(err)
	}

	log.Printf("Successfully processed %d documents", len(results))
}

// ProcessDocuments handles the main document processing pipeline
func (p *DocumentProcessor) ProcessDocuments(documents []models.DataSourceFile) ([]models.ProcessedDocument, error) {
	var results []models.ProcessedDocument

	// Create temp directory
	tempDir, err := os.MkdirTemp("", "doc-processor")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Process each document
	for _, doc := range documents {
		log.Printf("Processing document: %s", doc.Filename)
		// Generate chunks
		chunks, err := p.chunker.ProcessDocument(doc, tempDir)
		if err != nil {
			log.Printf("Error processing document %s: %v", doc.Filename, err)
			continue
		}

		// Generate embeddings
		embeddings, err := p.embedder.GenerateEmbeddings(p.ctx, chunks)
		if err != nil {
			log.Printf("Error generating embeddings for %s: %v", doc.Filename, err)
			continue
		}

		// Create and store document record
		queryDoc := models.QueryDocument{
			ID:            generateID(doc.Filename),
			QueryEngineID: p.queryEngineID,
			DocURL:        doc.URL,
			IndexFile:     doc.Filename,
			IndexStart:    0,
			IndexEnd:      len(chunks) - 1,
			CreatedAt:     time.Now(),
		}

		// Store document and chunks
		if err := p.storeDocumentAndChunks(queryDoc, chunks, embeddings); err != nil {
			log.Printf("Error storing document %s: %v", doc.Filename, err)
			continue
		}

		// Add to results
		result := models.ProcessedDocument{
			Filename:   doc.Filename,
			URL:        doc.URL,
			DocumentID: queryDoc.ID,
			ChunkCount: len(chunks),
		}
		results = append(results, result)
	}

	return results, nil
}

// storeDocumentAndChunks stores the document and its chunks in Firestore
func (p *DocumentProcessor) storeDocumentAndChunks(
	doc models.QueryDocument,
	chunks []string,
	embeddings [][]float32,
) error {
	batch := p.firestoreClient.Batch()

	// Store document
	docRef := p.firestoreClient.Collection("query_documents").Doc(doc.ID)
	batch.Set(docRef, doc)

	// Store chunks with embeddings
	for i, chunk := range chunks {
		chunkID := fmt.Sprintf("%s-chunk-%d", doc.ID, i)
		chunkRef := p.firestoreClient.Collection("query_document_chunks").Doc(chunkID)

		chunkDoc := models.QueryDocumentChunk{
			ID:              chunkID,
			QueryEngineID:   p.queryEngineID,
			QueryDocumentID: doc.ID,
			Index:           i,
			Modality:        "text",
			Text:            chunk,
			CleanText:       chunking.CleanText(chunk),
			CreatedAt:       time.Now(),
		}

		// Add embedding to chunk data
		chunkData := map[string]interface{}{
			"embedding": embeddings[i],
		}
		for k, v := range structToMap(chunkDoc) {
			chunkData[k] = v
		}

		batch.Set(chunkRef, chunkData)
	}

	// Commit the batch
	_, err := batch.Commit(p.ctx)
	return err
}

// Helper function to update job error status
func updateJobError(ctx context.Context, docRef *firestore.DocumentRef, err error) {
	_, updateErr := docRef.Update(ctx, []firestore.Update{
		{Path: "errors", Value: []string{err.Error()}},
		{Path: "status", Value: "failed"},
		{Path: "completed_at", Value: time.Now()},
	})
	if updateErr != nil {
		log.Printf("Failed to update job error status: %v", updateErr)
	}
}

// Helper function to generate document ID
func generateID(filename string) string {
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("%s-%d", sanitizeFilename(filename), timestamp)
}

// Helper function to sanitize filename for ID generation
func sanitizeFilename(filename string) string {
	// Remove extension and special characters
	name := strings.TrimSuffix(filename, filepath.Ext(filename))
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			return r
		}
		return '-'
	}, name)
}

// Helper function to convert struct to map
func structToMap(obj interface{}) map[string]interface{} {
	data, _ := json.Marshal(obj)
	var result map[string]interface{}
	json.Unmarshal(data, &result)
	return result
}
