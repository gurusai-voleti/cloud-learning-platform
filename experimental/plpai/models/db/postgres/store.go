package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"personalized_learning_platform/models"

	"github.com/lib/pq"
)

// Store implements the db.Store interface for PostgreSQL
type Store struct {
	db *sql.DB
}

// New creates a new PostgreSQL store
func New(dataSourceName string) (*Store, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	return &Store{
		db: db,
	}, nil
}

// CreateKnowledgeNode implements the Store interface
func (s *Store) CreateKnowledgeNode(ctx context.Context, node *models.KnowledgeNode) error {
	query := `
		INSERT INTO knowledge_nodes (
			id, title, description, type, domain, prerequisites, resources, 
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := s.db.ExecContext(ctx, query,
		node.ID,
		node.Title,
		node.Description,
		node.Type,
		node.Domain,
		pq.Array(node.Prerequisites),
		pq.Array(node.Resources),
		node.CreatedAt,
		node.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error creating knowledge node: %w", err)
	}

	return nil
}

// GetKnowledgeNode implements the Store interface
func (s *Store) GetKnowledgeNode(ctx context.Context, id string) (*models.KnowledgeNode, error) {
	query := `
		SELECT id, title, description, type, domain, prerequisites, resources, 
			   created_at, updated_at
		FROM knowledge_nodes
		WHERE id = $1
	`

	node := &models.KnowledgeNode{}
	var prerequisites, resources []string

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&node.ID,
		&node.Title,
		&node.Description,
		&node.Type,
		&node.Domain,
		pq.Array(&prerequisites),
		pq.Array(&resources),
		&node.CreatedAt,
		&node.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("knowledge node not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("error getting knowledge node: %w", err)
	}

	node.Prerequisites = prerequisites
	node.Resources = resources
	return node, nil
}

// Helper function to handle JSON fields
func jsonField(v interface{}) ([]byte, error) {
	if v == nil {
		return []byte("{}"), nil
	}
	return json.Marshal(v)
}

// ListKnowledgeNodes implements the Store interface
func (s *Store) ListKnowledgeNodes(ctx context.Context, domain string) ([]*models.KnowledgeNode, error) {
	query := `
		SELECT id, title, description, type, domain, prerequisites, resources, 
			   created_at, updated_at
		FROM knowledge_nodes
		WHERE domain = $1
	`

	rows, err := s.db.QueryContext(ctx, query, domain)
	if err != nil {
		return nil, fmt.Errorf("error listing knowledge nodes: %w", err)
	}
	defer rows.Close()

	var nodes []*models.KnowledgeNode
	for rows.Next() {
		node := &models.KnowledgeNode{}
		var prerequisites, resources []string

		err := rows.Scan(
			&node.ID,
			&node.Title,
			&node.Description,
			&node.Type,
			&node.Domain,
			pq.Array(&prerequisites),
			pq.Array(&resources),
			&node.CreatedAt,
			&node.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning knowledge node: %w", err)
		}

		node.Prerequisites = prerequisites
		node.Resources = resources
		nodes = append(nodes, node)
	}

	return nodes, nil
}

// UpdateKnowledgeNode implements the Store interface
func (s *Store) UpdateKnowledgeNode(ctx context.Context, node *models.KnowledgeNode) error {
	query := `
		UPDATE knowledge_nodes 
		SET title = $1, description = $2, type = $3, domain = $4, 
			prerequisites = $5, resources = $6, updated_at = $7
		WHERE id = $8
	`

	result, err := s.db.ExecContext(ctx, query,
		node.Title,
		node.Description,
		node.Type,
		node.Domain,
		pq.Array(node.Prerequisites),
		pq.Array(node.Resources),
		time.Now().UTC(),
		node.ID,
	)
	if err != nil {
		return fmt.Errorf("error updating knowledge node: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("knowledge node not found: %s", node.ID)
	}

	return nil
}

// DeleteKnowledgeNode implements the Store interface
func (s *Store) DeleteKnowledgeNode(ctx context.Context, id string) error {
	query := `DELETE FROM knowledge_nodes WHERE id = $1`

	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting knowledge node: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("knowledge node not found: %s", id)
	}

	return nil
}

// CreateKnowledgeEdge implements the Store interface
func (s *Store) CreateKnowledgeEdge(ctx context.Context, edge *models.KnowledgeEdge) error {
	query := `
		INSERT INTO knowledge_edges (id, source_id, target_id, relationship, weight, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := s.db.ExecContext(ctx, query,
		edge.ID,
		edge.SourceID,
		edge.TargetID,
		edge.Relationship,
		edge.Weight,
		edge.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("error creating knowledge edge: %w", err)
	}

	return nil
}

// GetKnowledgeEdges implements the Store interface
func (s *Store) GetKnowledgeEdges(ctx context.Context, nodeID string) ([]*models.KnowledgeEdge, error) {
	query := `
		SELECT id, source_id, target_id, relationship, weight, created_at
		FROM knowledge_edges
		WHERE source_id = $1 OR target_id = $1
	`

	rows, err := s.db.QueryContext(ctx, query, nodeID)
	if err != nil {
		return nil, fmt.Errorf("error getting knowledge edges: %w", err)
	}
	defer rows.Close()

	var edges []*models.KnowledgeEdge
	for rows.Next() {
		edge := &models.KnowledgeEdge{}
		err := rows.Scan(
			&edge.ID,
			&edge.SourceID,
			&edge.TargetID,
			&edge.Relationship,
			&edge.Weight,
			&edge.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning knowledge edge: %w", err)
		}
		edges = append(edges, edge)
	}

	return edges, nil
}

// DeleteKnowledgeEdge implements the Store interface
func (s *Store) DeleteKnowledgeEdge(ctx context.Context, id string) error {
	query := `DELETE FROM knowledge_edges WHERE id = $1`

	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting knowledge edge: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("knowledge edge not found: %s", id)
	}

	return nil
}

// CreateLearnerProfile implements the Store interface
func (s *Store) CreateLearnerProfile(ctx context.Context, profile *models.LearnerProfile) error {
	query := `
		INSERT INTO learner_profiles (
			id, user_id, knowledge_state, skill_state, learning_style, goals,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	knowledgeState, err := jsonField(profile.KnowledgeState)
	if err != nil {
		return fmt.Errorf("error marshaling knowledge state: %w", err)
	}

	skillState, err := jsonField(profile.SkillState)
	if err != nil {
		return fmt.Errorf("error marshaling skill state: %w", err)
	}

	learningStyle, err := jsonField(profile.LearningStyle)
	if err != nil {
		return fmt.Errorf("error marshaling learning style: %w", err)
	}

	goals, err := jsonField(profile.Goals)
	if err != nil {
		return fmt.Errorf("error marshaling goals: %w", err)
	}

	_, err = s.db.ExecContext(ctx, query,
		profile.ID,
		profile.UserID,
		knowledgeState,
		skillState,
		learningStyle,
		goals,
		profile.CreatedAt,
		profile.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("error creating learner profile: %w", err)
	}

	return nil
}

// GetLearnerProfile implements the Store interface
func (s *Store) GetLearnerProfile(ctx context.Context, id string) (*models.LearnerProfile, error) {
	query := `
		SELECT id, user_id, knowledge_state, skill_state, learning_style, goals,
			   created_at, updated_at
		FROM learner_profiles
		WHERE id = $1
	`

	profile := &models.LearnerProfile{}
	var knowledgeState, skillState, learningStyle, goals []byte

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&profile.ID,
		&profile.UserID,
		&knowledgeState,
		&skillState,
		&learningStyle,
		&goals,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("learner profile not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("error getting learner profile: %w", err)
	}

	if err := json.Unmarshal(knowledgeState, &profile.KnowledgeState); err != nil {
		return nil, fmt.Errorf("error unmarshaling knowledge state: %w", err)
	}
	if err := json.Unmarshal(skillState, &profile.SkillState); err != nil {
		return nil, fmt.Errorf("error unmarshaling skill state: %w", err)
	}
	if err := json.Unmarshal(learningStyle, &profile.LearningStyle); err != nil {
		return nil, fmt.Errorf("error unmarshaling learning style: %w", err)
	}
	if err := json.Unmarshal(goals, &profile.Goals); err != nil {
		return nil, fmt.Errorf("error unmarshaling goals: %w", err)
	}

	return profile, nil
}

// UpdateLearnerProfile implements the Store interface
func (s *Store) UpdateLearnerProfile(ctx context.Context, profile *models.LearnerProfile) error {
	query := `
		UPDATE learner_profiles 
		SET user_id = $1, knowledge_state = $2, skill_state = $3,
			learning_style = $4, goals = $5, updated_at = $6
		WHERE id = $7
	`

	knowledgeState, err := jsonField(profile.KnowledgeState)
	if err != nil {
		return fmt.Errorf("error marshaling knowledge state: %w", err)
	}

	skillState, err := jsonField(profile.SkillState)
	if err != nil {
		return fmt.Errorf("error marshaling skill state: %w", err)
	}

	learningStyle, err := jsonField(profile.LearningStyle)
	if err != nil {
		return fmt.Errorf("error marshaling learning style: %w", err)
	}

	goals, err := jsonField(profile.Goals)
	if err != nil {
		return fmt.Errorf("error marshaling goals: %w", err)
	}

	result, err := s.db.ExecContext(ctx, query,
		profile.UserID,
		knowledgeState,
		skillState,
		learningStyle,
		goals,
		time.Now().UTC(),
		profile.ID,
	)

	if err != nil {
		return fmt.Errorf("error updating learner profile: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("learner profile not found: %s", profile.ID)
	}

	return nil
}

// DeleteLearnerProfile implements the Store interface
func (s *Store) DeleteLearnerProfile(ctx context.Context, id string) error {
	query := `DELETE FROM learner_profiles WHERE id = $1`

	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting learner profile: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("learner profile not found: %s", id)
	}

	return nil
}
