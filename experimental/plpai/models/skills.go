package models

import (
	"time"
)

// SkillNode represents a specific skill or capability
type SkillNode struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	Category           string    `json:"category"` // cognitive, practical, social
	Level              int       `json:"level"`
	AssessmentCriteria []string  `json:"assessmentCriteria"`
	RequiredKnowledge  []string  `json:"requiredKnowledge"` // KnowledgeNode IDs
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
}

// SkillEdge represents a connection between two skill nodes
type SkillEdge struct {
	ID           string    `json:"id"`
	SourceID     string    `json:"sourceId"`
	TargetID     string    `json:"targetId"`
	Relationship string    `json:"relationship"` // builds_on, enables
	Weight       float64   `json:"weight"`       // 0-1
	CreatedAt    time.Time `json:"createdAt"`
}
