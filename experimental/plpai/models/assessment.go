package models

import (
	"time"
)

// Assessment represents an evaluation activity
type Assessment struct {
	ID                string                 `json:"id"`
	Type              string                 `json:"type"` // quiz, project, observation
	MeasuredKnowledge []string               `json:"measuredKnowledge"`
	MeasuredSkills    []string               `json:"measuredSkills"`
	Criteria          map[string]interface{} `json:"criteria"`
	AdaptivityRules   map[string]interface{} `json:"adaptivityRules"`
	CreatedAt         time.Time              `json:"createdAt"`
	UpdatedAt         time.Time              `json:"updatedAt"`
}

// LearnerProgress represents assessment results
type LearnerProgress struct {
	ID               string            `json:"id"`
	LearnerID        string            `json:"learnerId"`
	AssessmentID     string            `json:"assessmentId"`
	KnowledgeResults map[string]Result `json:"knowledgeResults"` // KnowledgeNode ID -> Result
	SkillResults     map[string]Result `json:"skillResults"`     // SkillNode ID -> Result
	Timestamp        time.Time         `json:"timestamp"`
	Feedback         string            `json:"feedback"`
}

// Result represents the outcome of an assessment for a specific knowledge or skill
type Result struct {
	Score      float64   `json:"score"`      // 0-1
	Confidence float64   `json:"confidence"` // 0-1
	Notes      string    `json:"notes"`
	AssessedAt time.Time `json:"assessedAt"`
}
