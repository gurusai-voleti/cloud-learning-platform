package db

import (
	"context"

	"personalized_learning_platform/models"
)

// Store defines the interface for database operations
type Store interface {
	// Knowledge operations
	CreateKnowledgeNode(ctx context.Context, node *models.KnowledgeNode) error
	GetKnowledgeNode(ctx context.Context, id string) (*models.KnowledgeNode, error)
	UpdateKnowledgeNode(ctx context.Context, node *models.KnowledgeNode) error
	DeleteKnowledgeNode(ctx context.Context, id string) error
	ListKnowledgeNodes(ctx context.Context, domain string) ([]*models.KnowledgeNode, error)

	// Knowledge Edge operations
	CreateKnowledgeEdge(ctx context.Context, edge *models.KnowledgeEdge) error
	GetKnowledgeEdges(ctx context.Context, nodeID string) ([]*models.KnowledgeEdge, error)
	DeleteKnowledgeEdge(ctx context.Context, id string) error

	// Skill operations
	CreateSkillNode(ctx context.Context, node *models.SkillNode) error
	GetSkillNode(ctx context.Context, id string) (*models.SkillNode, error)
	UpdateSkillNode(ctx context.Context, node *models.SkillNode) error
	DeleteSkillNode(ctx context.Context, id string) error
	ListSkillNodes(ctx context.Context, category string) ([]*models.SkillNode, error)

	// Skill Edge operations
	CreateSkillEdge(ctx context.Context, edge *models.SkillEdge) error
	GetSkillEdges(ctx context.Context, nodeID string) ([]*models.SkillEdge, error)
	DeleteSkillEdge(ctx context.Context, id string) error

	// Learner operations
	CreateLearnerProfile(ctx context.Context, profile *models.LearnerProfile) error
	GetLearnerProfile(ctx context.Context, id string) (*models.LearnerProfile, error)
	UpdateLearnerProfile(ctx context.Context, profile *models.LearnerProfile) error
	DeleteLearnerProfile(ctx context.Context, id string) error

	// Learning Path operations
	CreateLearningPath(ctx context.Context, path *models.LearningPath) error
	GetLearningPath(ctx context.Context, id string) (*models.LearningPath, error)
	UpdateLearningPath(ctx context.Context, path *models.LearningPath) error
	DeleteLearningPath(ctx context.Context, id string) error
	ListLearnerPaths(ctx context.Context, learnerID string) ([]*models.LearningPath, error)

	// Learning Resource operations
	CreateLearningResource(ctx context.Context, resource *models.LearningResource) error
	GetLearningResource(ctx context.Context, id string) (*models.LearningResource, error)
	UpdateLearningResource(ctx context.Context, resource *models.LearningResource) error
	DeleteLearningResource(ctx context.Context, id string) error
	ListLearningResources(ctx context.Context, filters map[string]interface{}) ([]*models.LearningResource, error)

	// Assessment operations
	CreateAssessment(ctx context.Context, assessment *models.Assessment) error
	GetAssessment(ctx context.Context, id string) (*models.Assessment, error)
	UpdateAssessment(ctx context.Context, assessment *models.Assessment) error
	DeleteAssessment(ctx context.Context, id string) error

	// Learner Progress operations
	CreateLearnerProgress(ctx context.Context, progress *models.LearnerProgress) error
	GetLearnerProgress(ctx context.Context, id string) (*models.LearnerProgress, error)
	ListLearnerProgress(ctx context.Context, learnerID string) ([]*models.LearnerProgress, error)
}
