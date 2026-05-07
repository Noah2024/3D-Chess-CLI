package save

import (
	"3DC/config"
	"3DC/util/metadata"
	"3DC/util/must"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kelindar/bitmap"
)

// Saves entire board state
func SaveGame(bmps map[string]bitmap.Bitmap, location string) {
	os.Mkdir(location, 0644) //Owner can rxw but everyone else can only r
	metadata.CreateSaveMetaData(location)
	for key, bm := range bmps {
		fileLoc := filepath.Join(location, string(key))
		file := must.Must(os.Create(fileLoc))
		bm.WriteTo(file)
	}
}

// Saves state for only one pieceType (lowkey need a better name)
func SavePieceType(vis string, bm bitmap.Bitmap) {
	fmt.Println(vis)
	fileLoc := filepath.Join(config.CurrentGame, vis)
	fmt.Println(fileLoc)
	file := must.Must(os.Create(fileLoc))
	bm.WriteTo(file)
}
