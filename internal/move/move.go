package move

import (
	"3DC/config"
	"3DC/internal/game/load"
	"3DC/internal/game/save"
	"3DC/util/bitutil"
	"3DC/util/logger"
	"3DC/util/metadata"
	"fmt"

	"github.com/kelindar/bitmap"
)

var m = map[string]int{
	"Alice": 25,
	"Bob":   30,
}

//Taking input from
//X               Z               Y
//a b c d e f g - 1 2 3 4 5 6 7 8 - A B C D E F G

// turns the user friendly a1A to uint32 location vector
func parseLoc(loc string) uint32 {

	if len(loc) != 3 {
		logger.Error(fmt.Sprintf("Could not parse location '%v' - invalid length of string", loc))
	}
	x, z, y := int(loc[0]-'a'+1), int(loc[1]-'1'+1), int(loc[2]-'A') //THIS ALSO NEEDS BETTER BOUNDS CHECKING

	return bitutil.VecToUint(x, y, z)
}

// Determines piece type at a given location // Needs a better name
func pieceType(loc uint32) (string, bitmap.Bitmap) {
	//Loading data from current game
	allPieces, _ := load.LoadGame(config.CurrentGame)

	//Contains is simd vectorized, I don't feel the need to optimize this search
	for meta, bm := range allPieces {

		if bm.Contains(loc) {
			logger.Info(meta)
			return meta, bm
		}

	}
	return "", bitmap.Bitmap{}
}

// This no work, recheck how the bitmaps contain each of the pieces
func Move(from string, to string) {

	uLocFrom := parseLoc(from)
	visFrom, bmFrom := pieceType(uLocFrom)
	if visFrom == "" {
		logger.Error(fmt.Sprintf("Could not find piece at location [%v]", from))
	}

	// logger.Info(fmt.Sprintf("%v", parseLoc(to)))
	uintLocTo := parseLoc(to)
	visTo, bmTo := pieceType(uintLocTo)

	//Updates bitmap of piece being moved
	bmFrom.Remove(uLocFrom)
	bmFrom.Set(uintLocTo)

	//Updates bitmap (if it exists) of piece being taken
	bmTo.Remove(uintLocTo)

	// logger.Warn(fmt.Sprintf("%v", visFrom))
	// logger.Warn(fmt.Sprintf("%v", visTo))

	metadata.CreateSaveMetaData(config.CurrentGame)
	save.SavePieceType(visFrom, bmFrom)

	if visTo != "" {
		save.SavePieceType(visTo, bmTo)
	}

	logger.Warn("Aint no way bro")
}
