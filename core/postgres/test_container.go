package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const pgUser = "postgres"
const pgPassword = "postgres"
const pgDB = "test_db"

func CleanUpTestContainer(c context.Context, t *testing.T, container testcontainers.Container, db DB) {
	db.Close()
	if err := container.Terminate(c); err != nil {
		t.Logf("failed to terminate container: %v", err)
	}
}

func NewTestContainer(ctx context.Context, t *testing.T) (testcontainers.Container, DB) {
	container, ip, port := startPostgresContainer(ctx, t)

	db, err := New(ctx, "postgres://"+pgUser+":"+pgPassword+"@"+ip+":"+port.Port()+"/"+pgDB+"?sslmode=disable",
		1*time.Minute, 30*time.Second, 4, 2)
	if err != nil {
		_ = container.Terminate(ctx)
		t.Fatal(err)
	}

	// Ensure DB is ready
	if err = pingWithRetry(ctx, db, 8, 750*time.Millisecond); err != nil {
		CleanUpTestContainer(ctx, t, container, db)
		t.Fatal(err)
	}

	if err = db.Migrate(); err != nil {
		CleanUpTestContainer(ctx, t, container, db)
		t.Fatal(err)
	}

	return container, db
}

func pingWithRetry(ctx context.Context, db DB, retries int, delay time.Duration) error {
	pool := db.Pool()
	var err error
	for i := 0; i < retries; i++ {
		pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
		err = pool.Ping(pingCtx)
		cancel()
		if err == nil {
			return nil
		}
		time.Sleep(delay)
	}
	return err
}

func startPostgresContainer(ctx context.Context, t *testing.T) (testcontainers.Container, string, nat.Port) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:18-alpine",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
		Env: map[string]string{
			"POSTGRES_USER":     pgUser,
			"POSTGRES_PASSWORD": pgPassword,
			"POSTGRES_DB":       pgDB,
		},
	}

	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start postgres container: %s", err)
	}

	ip, err := postgresC.Host(ctx)
	if err != nil {
		if termErr := postgresC.Terminate(ctx); termErr != nil {
			t.Logf("failed to terminate postgres container after host error: %v", termErr)
		}
		t.Fatalf("failed to get container host: %s", err)
	}

	port, err := postgresC.MappedPort(ctx, "5432")
	if err != nil {
		if termErr := postgresC.Terminate(ctx); termErr != nil {
			t.Logf("failed to terminate postgres container after mapped port error: %v", termErr)
		}
		t.Fatalf("failed to get container port: %s", err)
	}

	return postgresC, ip, port
}
