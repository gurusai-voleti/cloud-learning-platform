package models

import (
	"time"
)

// LearnerProfile represents a user's learning state and preferences
type LearnerProfile struct {
	ID             string                 `json:"id"`
	UserID         string                 `json:"userId"`
	KnowledgeState map[string]float64     `json:"knowledgeState"` // KnowledgeNode ID -> mastery level
	SkillState     map[string]float64     `json:"skillState"`     // SkillNode ID -> proficiency level
	LearningStyle  map[string]interface{} `json:"learningStyle"`
	Goals          []LearningGoal         `json:"goals"`
	CreatedAt      time.Time              `json:"createdAt"`
	UpdatedAt      time.Time              `json:"updatedAt"`
}

// LearningGoal represents a specific learning objective
type LearningGoal struct {
	ID              string    `json:"id"`
	Description     string    `json:"description"`
	TargetSkills    []string  `json:"targetSkills"`    // SkillNode IDs
	TargetKnowledge []string  `json:"targetKnowledge"` // KnowledgeNode IDs
	Deadline        time.Time `json:"deadline"`
	Status          string    `json:"status"` // not_started, in_progress, completed
}

// LearningPath represents a personalized learning journey
type LearningPath struct {
	ID              string                 `json:"id"`
	LearnerID       string                 `json:"learnerId"`
	Nodes           []LearningStep         `json:"nodes"`
	CurrentNode     int                    `json:"currentNode"`
	AdaptivityRules map[string]interface{} `json:"adaptivityRules"`
	CreatedAt       time.Time              `json:"createdAt"`
	UpdatedAt       time.Time              `json:"updatedAt"`
}

// LearningStep represents a single step in a learning path
type LearningStep struct {
	ID               string `json:"id"`
	ResourceID       string `json:"resourceId"`
	Type             string `json:"type"`         // content, assessment, practice
	RequiredTime     int    `json:"requiredTime"` // minutes
	CompletionStatus string `json:"completionStatus"`
}

// LearningResource represents educational content
type LearningResource struct {
	ID                string    `json:"id"`
	Title             string    `json:"title"`
	Type              string    `json:"type"`    // video, text, interactive
	Content           string    `json:"content"` // URL or content reference
	TargetKnowledge   []string  `json:"targetKnowledge"`
	TargetSkills      []string  `json:"targetSkills"`
	Difficulty        int       `json:"difficulty"`        // 1-5
	EstimatedDuration int       `json:"estimatedDuration"` // minutes
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}
