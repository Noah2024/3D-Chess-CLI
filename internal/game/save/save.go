// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

// save provides the functionality to save the current game state to a specified location.
package save

import (
	"3DC/config"
	"3DC/util/logger"
	"3DC/util/metadata"
	"3DC/util/must"
	"os"
	"path/filepath"

	"github.com/kelindar/bitmap"
)

// Saves entire board state
func SaveGame(bmps map[string]bitmap.Bitmap, meta metadata.MetaData, location string) error {
	os.Mkdir(location, 0o755) //Owner can rwx but everyone else can only r and x
	metadata.SaveMetaData(meta, location)
	for key, bm := range bmps {
		fileLoc := filepath.Join(location, string(key))
		file := must.Must(os.Create(fileLoc))
		_, err := bm.WriteTo(file)
		if err != nil {
			logger.Error("Unexpected Error in saving game type %s")
			return err
		}
	}
	return nil
}

// Saves state for only one pieceType (lowkey need a better name)
func SavePieceType(vis string, bm bitmap.Bitmap) error {
	// fmt.Println(vis)
	fileLoc := filepath.Join(config.CurrentGame, vis)
	// fmt.Println(fileLoc)
	file := must.Must(os.Create(fileLoc))
	_, err := bm.WriteTo(file)
	if err != nil {
		logger.Error("Unexpected Error in saving piece type %s")
		return err
	}
	return err
}
