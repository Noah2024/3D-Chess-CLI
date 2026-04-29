package save

import (
	"3DC/util/metadata"
	"3DC/util/must"
	"os"
	"path/filepath"

	"github.com/kelindar/bitmap"
)

func SaveGame(bmps map[string]bitmap.Bitmap, location string) {
	os.Mkdir(location, 0644) //Owner can rxw but everyone else can only r
	metadata.CreateSaveMetaData(filepath.Join(location, "meta"))
	for key, bm := range bmps {
		fileLoc := filepath.Join(location, string(key))
		file := must.Must(os.Create(fileLoc))
		bm.WriteTo(file)
	}
}
