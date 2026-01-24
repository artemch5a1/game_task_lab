package sqlite

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"example/web-service-gin/internal/domain/model"

	"github.com/google/uuid"
)

func TestSQLiteRepositories_GenresAndGames(t *testing.T) {
	ctx := context.Background()

	dbDir := t.TempDir()
	dbPath := filepath.Join(dbDir, "test.db")

	db, err := Open(ctx, Config{Path: dbPath})
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	genreRepo := NewGenreRepository(db.SQL)
	gameRepo := NewGameRepository(db.SQL)

	genre := &model.Genre{ID: uuid.New(), Title: "Action"}
	if _, err := genreRepo.Create(ctx, genre); err != nil {
		t.Fatalf("Create genre: %v", err)
	}

	releaseDate := time.Now().UTC().Truncate(time.Second)
	game := &model.Game{
		ID:          uuid.New(),
		Title:       "Test Game",
		Description: "desc",
		ReleaseDate: releaseDate,
		GenreID:     genre.ID,
	}
	if _, err := gameRepo.Create(ctx, game); err != nil {
		t.Fatalf("Create game: %v", err)
	}

	gotGame, err := gameRepo.FindByID(ctx, game.ID)
	if err != nil {
		t.Fatalf("FindByID game: %v", err)
	}
	if gotGame.GenreID != genre.ID {
		t.Fatalf("expected genreID %s, got %s", genre.ID, gotGame.GenreID)
	}

	genres, err := genreRepo.FindAll(ctx, 0, 0)
	if err != nil {
		t.Fatalf("FindAll genres: %v", err)
	}
	if len(genres) != 1 {
		t.Fatalf("expected 1 genre, got %d", len(genres))
	}
}

