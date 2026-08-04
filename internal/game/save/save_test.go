// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

// integration test for saving games
package save_test

import (
	"3DC/internal/game/save"
	"os"
	"path/filepath"
	"testing"

	"github.com/kelindar/bitmap"
)

func TestSaveGame(t *testing.T) {

	t.Run("TestValidSaveGame", func(t *testing.T) {
		dir := t.TempDir()
		location := filepath.Join(dir, "CurrentGame")

		bmps := map[string]bitmap.Bitmap{
			"♚": bitmap.Bitmap{}, //NNo information is really needed here
		}

		saveErr := save.SaveGame(bmps, location)
		if saveErr != nil {
			t.Errorf("Could not save game because of error %s", saveErr)
		}

		if _, err := os.Stat(location); err != nil {
			t.Errorf("Save was successful but no file was created? %s", saveErr)
		}
	})

	t.Run("TestSaveIndivudalPiece", func(t *testing.T) {
		dir := t.TempDir()
		location := filepath.Join(dir, "CurrentGame")

		bmps := map[string]bitmap.Bitmap{
			"♚": bitmap.Bitmap{}, //NNo information is really needed here
		}

		//Was already tested, so we use without testing now
		saveErr := save.SaveGame(bmps, location)
		if saveErr != nil {
			t.Errorf("Could not save game because of error %s", saveErr)
		}

		err := save.SavePieceType("♚", bmps["♚"])
		if err != nil {
			t.Errorf("Inital save was succesful, but indivual piece failed %s", err)
		}

	})

}
