package embedding

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	"cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	ITEMS_PER_REQUEST = 20
	TOKEN_LIMIT       = 20000
	DEFAULT_MODEL     = "text-embedding-004"
)

type VertexEmbedder struct {
	client    *aiplatform.PredictionClient
	projectID string
	location  string
	modelName string
}

func NewVertexEmbedder(ctx context.Context, projectID string) (*VertexEmbedder, error) {
	location := "us-central1"
	client, err := aiplatform.NewPredictionClient(ctx,
		option.WithEndpoint(fmt.Sprintf("%s-aiplatform.googleapis.com:443", location)))
	if err != nil {
		return nil, fmt.Errorf("failed to create Vertex AI client: %v", err)
	}

	modelName := os.Getenv("VERTEX_EMBEDDING_MODEL")
	if modelName == "" {
		modelName = DEFAULT_MODEL
	}

	return &VertexEmbedder{
		client:    client,
		projectID: projectID,
		location:  location,
		modelName: modelName,
	}, nil
}

func (v *VertexEmbedder) Close() {
	if v.client != nil {
		v.client.Close()
	}
}

func (v *VertexEmbedder) GenerateEmbeddings(ctx context.Context, texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, fmt.Errorf("no texts provided for embedding generation")
	}

	log.Printf("Generating embeddings for %d texts using model %s", len(texts), v.modelName)

	// Check total token length
	if !v.CheckTokenLimit(texts) {
		log.Printf("Warning: total tokens exceeds limit %d", TOKEN_LIMIT)
	}

	endpoint := fmt.Sprintf("projects/%s/locations/%s/publishers/google/models/%s",
		v.projectID,
		v.location,
		v.modelName)

	var allEmbeddings [][]float32
	var mu sync.Mutex
	var wg sync.WaitGroup
	errChan := make(chan error, 1)

	// Process in batches
	for i := 0; i < len(texts); i += ITEMS_PER_REQUEST {
		end := i + ITEMS_PER_REQUEST
		if end > len(texts) {
			end = len(texts)
		}

		batch := texts[i:end]
		wg.Add(1)

		go func(batchTexts []string) {
			defer wg.Done()

			// Create instances from texts
			instances := make([]*structpb.Value, len(batchTexts))
			for j, text := range batchTexts {
				instances[j] = structpb.NewStructValue(&structpb.Struct{
					Fields: map[string]*structpb.Value{
						"content": structpb.NewStringValue(text),
					},
				})
			}

			// Create the predict request
			req := &aiplatformpb.PredictRequest{
				Endpoint:  endpoint,
				Instances: instances,
			}

			// Call predict
			resp, err := v.client.Predict(ctx, req)
			if err != nil {
				select {
				case errChan <- fmt.Errorf("embedding generation failed: %v", err):
				default:
				}
				return
			}

			// Process response
			batchEmbeddings := make([][]float32, len(resp.Predictions))
			for j, prediction := range resp.Predictions {
				values := prediction.GetStructValue().
					Fields["embeddings"].GetStructValue().
					Fields["values"].GetListValue().Values

				embedding := make([]float32, len(values))
				for k, value := range values {
					embedding[k] = float32(value.GetNumberValue())
				}
				batchEmbeddings[j] = embedding
			}

			// Store embeddings thread-safely
			mu.Lock()
			allEmbeddings = append(allEmbeddings, batchEmbeddings...)
			mu.Unlock()
		}(batch)
	}

	// Wait for all goroutines
	wg.Wait()

	// Check for errors
	select {
	case err := <-errChan:
		return nil, err
	default:
		if len(allEmbeddings) == 0 {
			return nil, fmt.Errorf("no embeddings generated")
		}
		return allEmbeddings, nil
	}
}

// Helper method to estimate tokens in text
func (v *VertexEmbedder) EstimateTokens(text string) int {
	// Simple estimation - could be made more sophisticated
	return len(strings.Fields(text))
}

// Helper method to check if total tokens exceed limit
func (v *VertexEmbedder) CheckTokenLimit(texts []string) bool {
	var total int
	for _, text := range texts {
		total += v.EstimateTokens(text)
	}
	return total <= TOKEN_LIMIT
}
