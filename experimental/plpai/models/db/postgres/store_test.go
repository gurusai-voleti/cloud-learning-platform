package postgres

import (
	"context"
	"testing"
	"time"

	"personalized_learning_platform/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKnowledgeNodeOperations(t *testing.T) {
	store, cleanup, err := createTestStore(t)
	require.NoError(t, err)
	defer cleanup()

	ctx := context.Background()

	// Test Create
	node := &models.KnowledgeNode{
		ID:            uuid.New().String(),
		Title:         "Test Knowledge",
		Description:   "Test Description",
		Type:          "concept",
		Domain:        "mathematics",
		Prerequisites: []string{},
		Resources:     []string{},
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}

	err = store.CreateKnowledgeNode(ctx, node)
	assert.NoError(t, err)

	// Test Get
	retrieved, err := store.GetKnowledgeNode(ctx, node.ID)
	assert.NoError(t, err)
	assert.Equal(t, node.Title, retrieved.Title)
	assert.Equal(t, node.Description, retrieved.Description)
	assert.Equal(t, node.Type, retrieved.Type)
	assert.Equal(t, node.Domain, retrieved.Domain)

	// Test List
	nodes, err := store.ListKnowledgeNodes(ctx, "mathematics")
	assert.NoError(t, err)
	assert.Len(t, nodes, 1)
	assert.Equal(t, node.ID, nodes[0].ID)

	// Test Update
	node.Title = "Updated Title"
	err = store.UpdateKnowledgeNode(ctx, node)
	assert.NoError(t, err)

	updated, err := store.GetKnowledgeNode(ctx, node.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Title", updated.Title)

	// Test Delete
	err = store.DeleteKnowledgeNode(ctx, node.ID)
	assert.NoError(t, err)

	_, err = store.GetKnowledgeNode(ctx, node.ID)
	assert.Error(t, err)
}

func TestKnowledgeEdgeOperations(t *testing.T) {
	store, cleanup, err := createTestStore(t)
	require.NoError(t, err)
	defer cleanup()

	ctx := context.Background()

	// Create two nodes first
	node1 := &models.KnowledgeNode{
		ID:        uuid.New().String(),
		Title:     "Prerequisite",
		Type:      "concept",
		Domain:    "mathematics",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	node2 := &models.KnowledgeNode{
		ID:        uuid.New().String(),
		Title:     "Target",
		Type:      "concept",
		Domain:    "mathematics",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	require.NoError(t, store.CreateKnowledgeNode(ctx, node1))
	require.NoError(t, store.CreateKnowledgeNode(ctx, node2))

	// Test Create Edge
	edge := &models.KnowledgeEdge{
		ID:           uuid.New().String(),
		SourceID:     node1.ID,
		TargetID:     node2.ID,
		Relationship: "prerequisite",
		Weight:       0.8,
		CreatedAt:    time.Now().UTC(),
	}

	err = store.CreateKnowledgeEdge(ctx, edge)
	assert.NoError(t, err)

	// Test Get Edges
	edges, err := store.GetKnowledgeEdges(ctx, node1.ID)
	assert.NoError(t, err)
	assert.Len(t, edges, 1)
	assert.Equal(t, edge.ID, edges[0].ID)
	assert.Equal(t, edge.Weight, edges[0].Weight)

	// Test Delete Edge
	err = store.DeleteKnowledgeEdge(ctx, edge.ID)
	assert.NoError(t, err)

	edges, err = store.GetKnowledgeEdges(ctx, node1.ID)
	assert.NoError(t, err)
	assert.Len(t, edges, 0)
}

func TestLearnerProfileOperations(t *testing.T) {
	store, cleanup, err := createTestStore(t)
	require.NoError(t, err)
	defer cleanup()

	ctx := context.Background()

	profile := &models.LearnerProfile{
		ID:             uuid.New().String(),
		UserID:         uuid.New().String(),
		KnowledgeState: map[string]float64{"knowledge1": 0.8},
		SkillState:     map[string]float64{"skill1": 0.7},
		LearningStyle:  map[string]interface{}{"preferred_medium": "video"},
		Goals:          []models.LearningGoal{},
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}

	// Test Create
	err = store.CreateLearnerProfile(ctx, profile)
	assert.NoError(t, err)

	// Test Get
	retrieved, err := store.GetLearnerProfile(ctx, profile.ID)
	assert.NoError(t, err)
	assert.Equal(t, profile.UserID, retrieved.UserID)
	assert.Equal(t, profile.KnowledgeState["knowledge1"], retrieved.KnowledgeState["knowledge1"])
	assert.Equal(t, profile.SkillState["skill1"], retrieved.SkillState["skill1"])

	// Test Update
	profile.KnowledgeState["knowledge1"] = 0.9
	err = store.UpdateLearnerProfile(ctx, profile)
	assert.NoError(t, err)

	updated, err := store.GetLearnerProfile(ctx, profile.ID)
	assert.NoError(t, err)
	assert.Equal(t, 0.9, updated.KnowledgeState["knowledge1"])

	// Test Delete
	err = store.DeleteLearnerProfile(ctx, profile.ID)
	assert.NoError(t, err)

	_, err = store.GetLearnerProfile(ctx, profile.ID)
	assert.Error(t, err)
}
