package sqlite

import (
	"context"
	"path/filepath"
	"testing"

	"example/web-service-gin/internal/domain/model"
	specifictype "example/web-service-gin/internal/domain/specific_type"

	"github.com/google/uuid"
)

func TestSQLiteUserRepository_CreateFindByUsername(t *testing.T) {
	ctx := context.Background()
	dbPath := filepath.Join(t.TempDir(), "test.db")

	db, err := Open(ctx, Config{Path: dbPath})
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	repo := NewUserRepository(db.SQL)

	u := &model.User{
		ID:       uuid.New(),
		Username: "alice",
		Password: "pass",
		UserRole: specifictype.RoleUser,
	}
	if _, err := repo.Create(ctx, u); err != nil {
		t.Fatalf("Create: %v", err)
	}

	got, err := repo.FindByUsername(ctx, "alice")
	if err != nil {
		t.Fatalf("FindByUsername: %v", err)
	}
	if got.Username != "alice" {
		t.Fatalf("expected username alice, got %q", got.Username)
	}
}

