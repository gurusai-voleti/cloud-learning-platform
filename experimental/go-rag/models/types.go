package models

import (
	"time"
)

// QueryDocument matches Python QueryDocument model
type QueryDocument struct {
	ID            string            `json:"id"`
	QueryEngineID string            `json:"query_engine_id"`
	QueryEngine   string            `json:"query_engine"`
	DocURL        string            `json:"doc_url"`
	IndexFile     string            `json:"index_file"`
	IndexStart    int               `json:"index_start"`
	IndexEnd      int               `json:"index_end"`
	Metadata      map[string]string `json:"metadata"`
	CreatedAt     time.Time         `json:"created_at"`
}

// QueryDocumentChunk matches Python QueryDocumentChunk model
type QueryDocumentChunk struct {
	ID              string    `json:"id"`
	QueryEngineID   string    `json:"query_engine_id"`
	QueryDocumentID string    `json:"query_document_id"`
	Index           int       `json:"index"`
	Modality        string    `json:"modality"`
	Page            int       `json:"page"`
	Text            string    `json:"text"`
	CleanText       string    `json:"clean_text"`
	CreatedAt       time.Time `json:"created_at"`
}

// JobInput represents the input data for the document processing job
type JobInput struct {
	GCSBucket     string           `json:"gcs_bucket"`
	Documents     []DataSourceFile `json:"documents"`
	QueryEngineID string           `json:"query_engine_id"`
}

// DataSourceFile matches the webscraper output format
type DataSourceFile struct {
	Filename    string `json:"filename"`
	URL         string `json:"url"`
	GCSPath     string `json:"gcs_path"`
	ContentType string `json:"content_type"`
}

// ProcessedDocument represents the result of processing a document
type ProcessedDocument struct {
	Filename     string   `json:"filename"`
	URL          string   `json:"url"`
	DocumentID   string   `json:"document_id"`
	ChunkCount   int      `json:"chunk_count"`
	EmbeddingIDs []string `json:"embedding_ids"`
}
