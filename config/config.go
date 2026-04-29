// Embedded config values
package config

import (
	"3DC/util/must"
	"os"
	"path/filepath"
)

// Defining size and shape of board
// Stored in Uints right now to make Uint -> Vec easier
// BUT it may be benificial later to store them as ints
// And to make Vec -> Uint easier
const BoardSize uint32 = 512
const LayerSize uint32 = 64
const LineSize uint32 = 8
const SpaceSize uint32 = 1

// For debug purposes you can manually set the Data Directory
// Otherwise it generates with os.UserConfigDir
var DataDir string
var CurrentGame = "data/CurrentGame"

func init() {

	if DataDir == "" {
		userDir := must.Must(os.UserConfigDir())
		DataDir = filepath.Join(userDir, "3DC/DATA")
		err := os.MkdirAll(DataDir, 0644)
		must.Must("", err)
	}
	CurrentGame = filepath.Join(DataDir, "CurrentGame")
}
