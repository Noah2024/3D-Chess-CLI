package load

import (
	"3DC/util/logger"
	"3DC/util/must"
	"os"

	"github.com/kelindar/bitmap"
)

// Loads dictionary mapping display char to the bitmap corresponding with that piece
// Also returns
func LoadGame(location string) (data map[string]bitmap.Bitmap, err error) {

	result := make(map[string]bitmap.Bitmap)
	if _, err := os.Stat(location); os.IsNotExist(err) {
		logger.Warn("No game currently running, create one with '3DC game new'")
		return nil, err
	}
	entries := must.Must(os.ReadDir(location))

	for _, entry := range entries {
		if entry.IsDir() {
			//meta.bin is stored in /meta dir
			//This is so we can easily skip it when loading the piece bitmaps
			continue
		}
		file := must.Must(os.Open(location + "/" + entry.Name()))

		bm := must.Must(bitmap.ReadFrom(file))
		result[entry.Name()] = bm
	}
	return result, nil

}

//Loads current game to bitmap AND handles checking
func LoadCurrentGameToBitMap() {

}
