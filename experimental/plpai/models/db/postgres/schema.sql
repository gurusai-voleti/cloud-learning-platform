-- Knowledge Graph tables
CREATE TABLE knowledge_nodes (
    id VARCHAR(255) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    type VARCHAR(50) NOT NULL,
    domain VARCHAR(100) NOT NULL,
    prerequisites TEXT[],
    resources TEXT[],
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE knowledge_edges (
    id VARCHAR(255) PRIMARY KEY,
    source_id VARCHAR(255) NOT NULL REFERENCES knowledge_nodes(id),
    target_id VARCHAR(255) NOT NULL REFERENCES knowledge_nodes(id),
    relationship VARCHAR(50) NOT NULL,
    weight FLOAT NOT NULL,
    created_at TIMESTAMP NOT NULL
);

-- Skills Graph tables
CREATE TABLE skill_nodes (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(50) NOT NULL,
    level INTEGER NOT NULL,
    assessment_criteria TEXT[],
    required_knowledge TEXT[],
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE skill_edges (
    id VARCHAR(255) PRIMARY KEY,
    source_id VARCHAR(255) NOT NULL REFERENCES skill_nodes(id),
    target_id VARCHAR(255) NOT NULL REFERENCES skill_nodes(id),
    relationship VARCHAR(50) NOT NULL,
    weight FLOAT NOT NULL,
    created_at TIMESTAMP NOT NULL
);

-- Learning tables
CREATE TABLE learner_profiles (
    id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    knowledge_state JSONB NOT NULL DEFAULT '{}',
    skill_state JSONB NOT NULL DEFAULT '{}',
    learning_style JSONB NOT NULL DEFAULT '{}',
    goals JSONB NOT NULL DEFAULT '[]',
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE learning_paths (
    id VARCHAR(255) PRIMARY KEY,
    learner_id VARCHAR(255) NOT NULL REFERENCES learner_profiles(id),
    nodes JSONB NOT NULL,
    current_node INTEGER NOT NULL DEFAULT 0,
    adaptivity_rules JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE learning_resources (
    id VARCHAR(255) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,
    content TEXT NOT NULL,
    target_knowledge TEXT[],
    target_skills TEXT[],
    difficulty INTEGER NOT NULL,
    estimated_duration INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Assessment tables
CREATE TABLE assessments (
    id VARCHAR(255) PRIMARY KEY,
    type VARCHAR(50) NOT NULL,
    measured_knowledge TEXT[],
    measured_skills TEXT[],
    criteria JSONB NOT NULL DEFAULT '{}',
    adaptivity_rules JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE TABLE learner_progress (
    id VARCHAR(255) PRIMARY KEY,
    learner_id VARCHAR(255) NOT NULL REFERENCES learner_profiles(id),
    assessment_id VARCHAR(255) NOT NULL REFERENCES assessments(id),
    knowledge_results JSONB NOT NULL DEFAULT '{}',
    skill_results JSONB NOT NULL DEFAULT '{}',
    timestamp TIMESTAMP NOT NULL,
    feedback TEXT
);

-- Indexes for better query performance
CREATE INDEX idx_knowledge_nodes_domain ON knowledge_nodes(domain);
CREATE INDEX idx_skill_nodes_category ON skill_nodes(category);
CREATE INDEX idx_learning_paths_learner ON learning_paths(learner_id);
CREATE INDEX idx_learner_progress_learner ON learner_progress(learner_id); 