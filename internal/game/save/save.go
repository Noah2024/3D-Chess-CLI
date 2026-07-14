// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

package save

import (
	"3DC/config"
	"3DC/util/metadata"
	"3DC/util/must"
	"log"
	"os"
	"path/filepath"

	"github.com/kelindar/bitmap"
)

// Saves entire board state
func SaveGame(bmps map[string]bitmap.Bitmap, location string) {
	os.Mkdir(location, 0o755) //Owner can rwx but everyone else can only r and x
	metadata.CreateSaveMetaData(location)
	for key, bm := range bmps {
		fileLoc := filepath.Join(location, string(key))
		file := must.Must(os.Create(fileLoc))
		bm.WriteTo(file)
	}
}

// Saves state for only one pieceType (lowkey need a better name)
func SavePieceType(vis string, bm bitmap.Bitmap) {
	// fmt.Println(vis)
	fileLoc := filepath.Join(config.CurrentGame, vis)
	// fmt.Println(fileLoc)
	file := must.Must(os.Create(fileLoc))
	bm.WriteTo(file)
}

func SaveDebugBoard(str string, location string) {
	os.Mkdir(location, 0o755) //Owner can rwx but everyone else can only r and x
	fileLoc := filepath.Join(location, "DebugBoardState")
	// file := must.Must(os.Create(fileLoc))
	err := os.WriteFile(fileLoc, []byte(str), 0o755)
	if err != nil {
		log.Fatal(err)
	}
	// bm.WriteTo(file)
}
