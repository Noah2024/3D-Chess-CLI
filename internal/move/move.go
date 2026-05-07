package move

import (
	"3DC/internal/game/load"
	"3DC/util/bitutil"
	"3DC/util/logger"
	"fmt"
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
	// fmt.Printf("%v %v %v\n", int(loc[0]-'a'+1), int(loc[1]-'1'+1), int(loc[2]-'A')+1)
	x, z, y := int(loc[0]-'a'+1), int(loc[1]-'1'+1), int(loc[2]-'A'+1) //THIS ALSO NEEDS BETTER BOUNDS CHECKING

	return bitutil.VecToUint(x, y, z) //NEEDS BETTER BOUNDS CHECKING
}

// Determines piece type at a given location
func pieceType(loc uint32) {
	//Loading data from current game
	allPieces, _ := load.LoadGame("data/output")

	// meta := must.Must(metadata.LoadMetaData(filepath.Join("data/output", "meta")))
	for meta, bm := range allPieces {

		if bm.Contains(loc) {
			logger.Info(meta)
		}

	}
}

func Move(to string, from string) {

	logger.Info(fmt.Sprintf("%v", parseLoc(to)))
	// logger.Info(fmt.Sprintf("%v", parseLoc(from)))
	pieceType(parseLoc(to))
	logger.Warn("Aint no way bro")

}
