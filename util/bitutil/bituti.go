package bitutil

import (
	"3DC/config"
	"3DC/util/logger"
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

const (
	BoardSize = config.BoardSize
	LayerSize = config.LayerSize
	LineSize  = config.LineSize
	SpaceSize = config.SpaceSize
)

func VecToUint(x, y, z int) uint32 {
	return uint32(x + (y-1)*int(LayerSize) + (z-1)*int(LineSize))
}

// Decodes uint32 position into integer x,y,z position
func UintToVec(space uint32) (int, int, int) {
	if space < 1 || space > 512 {
		logger.Error(fmt.Sprintf("uint32 %d out of range for board size %d ", space, BoardSize))
		panic("See above error")
	}

	// index = x + y*8 + z*64 essentially decoding this
	//Step by step removing the largest term at a time
	space-- // convert to 0-based
	y := space / LayerSize
	space %= LayerSize
	z := space / LineSize
	x := space % LineSize

	return int(x + 1), int(y + 1), int(z + 1)
}

func SaveGame(bmps map[rune]bitmap.Bitmap, root string) {
	os.Mkdir(root, 0744) //Owner can rxw but everyone else can only r

	for key, bm := range bmps {
		fileLoc := root + "/" + string(key)
		file := must.Must(os.Create(fileLoc))
		bm.WriteTo(file)
	}
}

func LoadGame(location string) (data map[string]bitmap.Bitmap, err error) {

	result := make(map[string]bitmap.Bitmap)

	entries := must.Must(os.ReadDir(location))

	for _, entry := range entries {
		//Will handcode some way to skip metadata latr
		if entry.IsDir() {
			continue
		}
		// fmt.Printf("%d, %d", location
		// , entry.Name())
		file := must.Must(os.Open(location + "/" + entry.Name()))

		// logger.Warn(strconv.Itoa(len(file)))
		bm := must.Must(bitmap.ReadFrom(file))
		// fmt.Println(rune(entry.Name()[0]))
		result[entry.Name()] = bm
	}
	return result, nil

}
