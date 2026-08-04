// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

package delete_test

import (
	"3DC/config"
	"3DC/internal/game/delete"
	"os"
	"path/filepath"
	"testing"
)

func TestDelete(t *testing.T) {
	t.Run("DeleteByName", func(t *testing.T) {
		// Arrange
		config.DataDir = t.TempDir()

		gamePath := filepath.Join(config.DataDir, "game1")
		if err := os.WriteFile(gamePath, []byte{}, 0644); err != nil {
			t.Fatal(err)
		}

		// Act
		err := delete.DeleteGame("game1")

		// Assert
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if _, err := os.Stat(gamePath); !os.IsNotExist(err) {
			t.Errorf("expected game to be deleted")
		}
	})

	t.Run("DeleteByIndex", func(t *testing.T) {
		// Arrange
		config.DataDir = t.TempDir()

		gamePath := filepath.Join(config.DataDir, "game1")
		if err := os.WriteFile(gamePath, []byte{}, 0644); err != nil {
			t.Fatal(err)
		}

		// Act
		err := delete.DeleteGame("0")

		// Assert
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if _, err := os.Stat(gamePath); !os.IsNotExist(err) {
			t.Errorf("expected indexed game to be deleted")
		}
	})

	t.Run("DeleteGameNotFound", func(t *testing.T) {
		// Arrange
		config.DataDir = t.TempDir()

		gamePath := filepath.Join(config.DataDir, "game1")
		if err := os.WriteFile(gamePath, []byte{}, 0644); err != nil {
			t.Fatal(err)
		}

		// Act
		err := delete.DeleteGame("does-not-exist")

		// Assert
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Ensure existing files were untouched
		if _, err := os.Stat(gamePath); os.IsNotExist(err) {
			t.Errorf("existing game was incorrectly deleted")
		}
	})

	t.Run("DeletesOnlyRequestedGame", func(t *testing.T) {
		// Arrange
		config.DataDir = t.TempDir()

		game1 := filepath.Join(config.DataDir, "game1")
		game2 := filepath.Join(config.DataDir, "game2")

		for _, path := range []string{game1, game2} {
			if err := os.WriteFile(path, []byte{}, 0644); err != nil {
				t.Fatal(err)
			}
		}

		// Act
		err := delete.DeleteGame("game1")

		// Assert
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if _, err := os.Stat(game1); !os.IsNotExist(err) {
			t.Errorf("expected game1 to be deleted")
		}

		if _, err := os.Stat(game2); os.IsNotExist(err) {
			t.Errorf("game2 was incorrectly deleted")
		}
	})
}
