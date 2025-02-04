# Models

This folder contains the data model definitions for the Personalized Learning Platform.

## Core Models

### Knowledge Graph Models

1. **KnowledgeNode**
   - `id`: Unique identifier
   - `title`: Name of the knowledge concept
   - `description`: Detailed explanation
   - `type`: Type of knowledge (e.g., concept, fact, principle)
   - `domain`: Subject area/domain
   - `prerequisites`: Array of KnowledgeNode IDs required before learning this
   - `resources`: Array of learning resources associated with this node

2. **KnowledgeEdge**
   - `id`: Unique identifier
   - `sourceId`: ID of the prerequisite KnowledgeNode
   - `targetId`: ID of the dependent KnowledgeNode
   - `relationship`: Type of connection (e.g., "prerequisite", "related", "part_of")
   - `weight`: Strength of the relationship (0-1)

### Skills Graph Models

1. **SkillNode**
   - `id`: Unique identifier
   - `name`: Name of the skill
   - `description`: Detailed description
   - `category`: Skill category (e.g., cognitive, practical, social)
   - `level`: Difficulty/complexity level
   - `assessmentCriteria`: Array of criteria to evaluate this skill
   - `requiredKnowledge`: Array of KnowledgeNode IDs needed for this skill

2. **SkillEdge**
   - `id`: Unique identifier
   - `sourceId`: ID of the prerequisite SkillNode
   - `targetId`: ID of the dependent SkillNode
   - `relationship`: Type of connection (e.g., "builds_on", "enables")
   - `weight`: Strength of the relationship (0-1)

### Learning Models

1. **LearnerProfile**
   - `id`: Unique identifier
   - `userId`: Reference to user
   - `knowledgeState`: Map of KnowledgeNode IDs to mastery levels
   - `skillState`: Map of SkillNode IDs to proficiency levels
   - `learningStyle`: Learning preferences and style indicators
   - `goals`: Array of learning objectives

2. **LearningPath**
   - `id`: Unique identifier
   - `learnerId`: Reference to LearnerProfile
   - `nodes`: Ordered array of learning steps
   - `currentNode`: Current position in the path
   - `adaptivityRules`: Rules for path modification based on progress

3. **LearningResource**
   - `id`: Unique identifier
   - `title`: Resource title
   - `type`: Resource type (video, text, interactive, etc.)
   - `content`: Content or reference to content
   - `targetKnowledge`: Array of KnowledgeNode IDs
   - `targetSkills`: Array of SkillNode IDs
   - `difficulty`: Difficulty level
   - `estimatedDuration`: Time to complete in minutes

### Assessment Models

1. **Assessment**
   - `id`: Unique identifier
   - `type`: Assessment type (quiz, project, observation, etc.)
   - `measuredKnowledge`: Array of KnowledgeNode IDs being assessed
   - `measuredSkills`: Array of SkillNode IDs being assessed
   - `criteria`: Assessment criteria and rubrics
   - `adaptivityRules`: Rules for difficulty adjustment

2. **LearnerProgress**
   - `id`: Unique identifier
   - `learnerId`: Reference to LearnerProfile
   - `assessmentId`: Reference to Assessment
   - `knowledgeResults`: Map of KnowledgeNode IDs to assessment results
   - `skillResults`: Map of SkillNode IDs to assessment results
   - `timestamp`: When the assessment was completed
   - `feedback`: Detailed feedback and recommendations

This schema separates knowledge and skills into distinct graphs while maintaining their relationships through references. The knowledge graph represents conceptual understanding, while the skills graph represents practical abilities. The learning and assessment models tie everything together to create personalized learning experiences.

## Testing

The models package includes comprehensive tests that use a local PostgreSQL instance via testcontainers. 

### Prerequisites

Before running tests, ensure you have:

1. Docker installed and running
2. Go 1.16 or later
3. `fswatch` (optional, for watch mode)

### Running Tests

A test script is provided at `scripts/test.sh` with the following options:

```bash
# Run all tests
./scripts/test.sh

# Run tests with verbose output
./scripts/test.sh -v

# Run tests with race detection
./scripts/test.sh -r

# Generate test coverage report
./scripts/test.sh -c

# Test specific package
./scripts/test.sh -p ./models/db/postgres

# Watch mode - automatically run tests on file changes
./scripts/test.sh -w

# Combine multiple options
./scripts/test.sh -v -r -c -p ./models/db/postgres
```

### Coverage Reports

When using the `-c` flag, two coverage reports are generated:
- `coverage.out`: Raw coverage data
- `coverage.html`: HTML coverage report with line-by-line analysis

### Watch Mode

To use watch mode (`-w` flag), you'll need to install `fswatch`:

For macOS:
```bash
brew install fswatch
```

For Ubuntu/Debian:
```bash
apt-get install fswatch
```

The watch mode will automatically re-run tests whenever source files change, making it useful during development.

### Test Structure

The tests are organized to verify:
- Data model integrity
- Database operations (CRUD)
- Relationships between entities
- Edge cases and error conditions
- Concurrent operations (with -r flag)

Each model component (Knowledge, Skills, Learning, Assessment) has its own test suite that can be run independently using the `-p` flag.


