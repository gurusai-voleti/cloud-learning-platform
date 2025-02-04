package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresContainer struct {
	Container testcontainers.Container
	URI       string
}

// setupTestContainer creates a new Postgres container for testing
func setupTestContainer(t *testing.T) (*PostgresContainer, error) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:13",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       "testdb",
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return nil, err
	}

	hostIP, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("postgres://test:test@%s:%s/testdb?sslmode=disable",
		hostIP, mappedPort.Port())

	return &PostgresContainer{
		Container: container,
		URI:       uri,
	}, nil
}

// setupTestDB initializes the test database with schema
func setupTestDB(t *testing.T, db *sql.DB) error {
	schemaPath := filepath.Join("schema.sql")
	schemaSQL, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("error reading schema file: %w", err)
	}

	_, err = db.Exec(string(schemaSQL))
	if err != nil {
		return fmt.Errorf("error executing schema: %w", err)
	}

	return nil
}

// createTestStore creates a new store instance for testing
func createTestStore(t *testing.T) (*Store, func(), error) {
	container, err := setupTestContainer(t)
	if err != nil {
		return nil, nil, err
	}

	store, err := New(container.URI)
	if err != nil {
		return nil, nil, err
	}

	err = setupTestDB(t, store.db)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		store.db.Close()
		container.Container.Terminate(context.Background())
	}

	return store, cleanup, nil
}
