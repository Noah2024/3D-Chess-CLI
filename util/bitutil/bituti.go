package bitutil

import (
	"3DC/util/must"
	"fmt"
	"os"

	"github.com/kelindar/bitmap"
)

//Data to store
//Board Size
//Config
//Turn
//Castling Rights
//EnPessantRights
//BitMaps for each type of piece

//metadata schema
// 1 byte version
// 1 byte board size
// 1 byte number of pieces
//

func SaveGame(bmps map[string]bitmap.Bitmap, root string) {
	os.Mkdir(root, 0744) //Owner can rxw but everyone else can only r

	for key, bm := range bmps {
		fileLoc := root + "/" + key
		file := must.Must(os.Create(fileLoc))
		bm.WriteTo(file)
	}
}

func LoadGame(location string) (data []bitmap.Bitmap, err error) {

	var result []bitmap.Bitmap

	entries := must.Must(os.ReadDir(location))

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		fmt.Println(entry.Name())
	}
	return result, nil

}
