// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

package list_test

import (
	"3DC/config"
	"3DC/internal/game/list"
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

//Test is AI Generated

func TestListGames(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	config.DataDir = tempDir

	if err := os.WriteFile(filepath.Join(tempDir, "game1"), []byte{}, 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tempDir, "game2"), []byte{}, 0644); err != nil {
		t.Fatal(err)
	}

	var output bytes.Buffer

	// Act
	list.ListGames(&output)

	// Assert
	got := output.String()

	if !strings.Contains(got, "All Boards avilable to load") {
		t.Errorf("expected header in output\n%s", got)
	}

	if !strings.Contains(got, "----------") {
		t.Errorf("expected separator in output\n%s", got)
	}

	if !strings.Contains(got, "game1") {
		t.Errorf("expected game1 to be listed\n%s", got)
	}

	if !strings.Contains(got, "game2") {
		t.Errorf("expected game2 to be listed\n%s", got)
	}
}
