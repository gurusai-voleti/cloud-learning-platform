package main

import (
	"context"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"document-processor/models"
)

// Mock DocumentChunker
type mockDocumentChunker struct {
	mock.Mock
}

func (m *mockDocumentChunker) ProcessDocument(doc models.DataSourceFile, tempDir string) ([]string, error) {
	args := m.Called(doc, tempDir)
	return args.Get(0).([]string), args.Error(1)
}

// Mock VertexEmbedder
type mockVertexEmbedder struct {
	mock.Mock
}

func (m *mockVertexEmbedder) GenerateEmbeddings(ctx context.Context, chunks []string) ([][]float32, error) {
	args := m.Called(ctx, chunks)
	return args.Get(0).([][]float32), args.Error(1)
}

func TestNewDocumentProcessor(t *testing.T) {
	ctx := context.Background()

	// Setup Firestore emulator client
	firestoreClient, err := firestore.NewClient(ctx, "test-project")
	assert.NoError(t, err)
	defer firestoreClient.Close()

	projectID := "test-project"
	queryEngineID := "test-engine"
	jobRef := firestoreClient.Collection("batch_jobs").Doc("test-job")

	processor, err := NewDocumentProcessor(ctx, firestoreClient, projectID, queryEngineID, jobRef)

	assert.NoError(t, err)
	assert.NotNil(t, processor)
	assert.Equal(t, ctx, processor.ctx)
	assert.NotNil(t, processor.firestoreClient)
	assert.Equal(t, projectID, processor.projectID)
	assert.Equal(t, queryEngineID, processor.queryEngineID)
	assert.Equal(t, jobRef, processor.jobRef)
}

func TestProcessDocuments(t *testing.T) {
	ctx := context.Background()

	// Setup Firestore emulator client
	firestoreClient, err := firestore.NewClient(ctx, "fake-project")
	assert.NoError(t, err)
	defer firestoreClient.Close()

	mockChunker := &mockDocumentChunker{}
	mockEmbedder := &mockVertexEmbedder{}

	processor := &DocumentProcessor{
		ctx:             ctx,
		firestoreClient: firestoreClient,
		chunker:         mockChunker,
		embedder:        mockEmbedder,
		projectID:       "fake-project",
		queryEngineID:   "test-engine",
		jobRef:          firestoreClient.Collection("batch_jobs").Doc("test-job"),
	}

	testDocs := []models.DataSourceFile{
		{Filename: "test1.txt", URL: "http://example.com/test1.txt"},
		{Filename: "test2.txt", URL: "http://example.com/test2.txt"},
	}

	mockChunker.On("ProcessDocument", mock.Anything, mock.Anything).Return([]string{"chunk1", "chunk2"}, nil)
	mockEmbedder.On("GenerateEmbeddings", mock.Anything, mock.Anything).Return([][]float32{{1.0, 2.0}, {3.0, 4.0}}, nil)

	results, err := processor.ProcessDocuments(testDocs)

	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, "test1.txt", results[0].Filename)
	assert.Equal(t, "http://example.com/test1.txt", results[0].URL)
	assert.Equal(t, 2, results[0].ChunkCount)

	mockChunker.AssertExpectations(t)
	mockEmbedder.AssertExpectations(t)
}

func TestGenerateID(t *testing.T) {
	filename := "test file.txt"
	id := generateID(filename)

	assert.Contains(t, id, "test-file-")
	assert.Regexp(t, `^test-file-\d+$`, id)
}

func TestSanitizeFilename(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"test_file.txt", "test-file"},
		{"Test File 123.pdf", "Test-File-123"},
		{"!@#$%^&*.doc", "---"},
	}

	for _, tc := range testCases {
		result := sanitizeFilename(tc.input)
		assert.Equal(t, tc.expected, result)
	}
}

func TestStructToMap(t *testing.T) {
	type TestStruct struct {
		Name  string    `json:"name"`
		Age   int       `json:"age"`
		Date  time.Time `json:"date"`
		Empty string    `json:"empty,omitempty"`
	}

	testTime := time.Date(2023, 4, 15, 12, 0, 0, 0, time.UTC)
	testStruct := TestStruct{
		Name: "John Doe",
		Age:  30,
		Date: testTime,
	}

	result := structToMap(testStruct)

	assert.Equal(t, "John Doe", result["name"])
	assert.Equal(t, float64(30), result["age"])
	assert.Equal(t, testTime.Format(time.RFC3339), result["date"])
	_, hasEmpty := result["empty"]
	assert.False(t, hasEmpty)
}
