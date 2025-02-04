package models

import (
	"time"
)

// KnowledgeNode represents a single concept or piece of knowledge
type KnowledgeNode struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Type          string    `json:"type"` // concept, fact, principle
	Domain        string    `json:"domain"`
	Prerequisites []string  `json:"prerequisites"` // Array of KnowledgeNode IDs
	Resources     []string  `json:"resources"`     // Array of LearningResource IDs
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// KnowledgeEdge represents a connection between two knowledge nodes
type KnowledgeEdge struct {
	ID           string    `json:"id"`
	SourceID     string    `json:"sourceId"`
	TargetID     string    `json:"targetId"`
	Relationship string    `json:"relationship"` // prerequisite, related, part_of
	Weight       float64   `json:"weight"`       // 0-1
	CreatedAt    time.Time `json:"createdAt"`
}
