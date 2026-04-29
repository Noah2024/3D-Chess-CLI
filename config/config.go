// Embedded config values
package config

import (
	"3DC/util/color"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Controls what log level the users sees.
// All logs are output to LOG dir no matter what
var LogLevel LOGLEVEL = 1

type LOGLEVEL int

const (
	debug = iota // 0
	info         // 1
	warn         // 2
	err          // 3
	fatal        // 4
)

// Defining size and shape of board
// Stored in Uints right now to make Uint -> Vec easier
// BUT it may be benificial later to store them as ints
// And to make Vec -> Uint easier
const BoardSize uint32 = 512
const LayerSize uint32 = 64
const LineSize uint32 = 8
const SpaceSize uint32 = 1

var DataDir string
var CurrentGame string
var LogDir string
var CurrentLog string

// Implementation of must function completly internal to the config file.
// Becuase log is not yet implemneted config needs its own internal must for displaying BAD BAD BAD errors
func internalMust[T any](val T, err error) T {
	if err != nil {
		fmt.Print(color.ColorText("!!FATAL ERROR IN STATIC CONFIG!! \n", color.Purple))
		fmt.Print(color.ColorText(err.Error(), color.Purple))
		os.Exit(1)
	}
	return val
}

// Initalize all the requiset data directories
func init() {

	//Setting up the main user directory
	userDir := internalMust(os.UserConfigDir())

	//Data directory storing game states
	DataDir = filepath.Join(userDir, "3DC/DATA")
	err := os.MkdirAll(DataDir, 0644)
	internalMust("", err)

	//Establishing where the Current Game is
	CurrentGame = filepath.Join(DataDir, "CurrentGame")

	//Defining where the path to the
	LogDir = filepath.Join(userDir, "3DC/LOG")
	err2 := os.MkdirAll(LogDir, 0644)
	internalMust("", err2)

	//Creating the log file if it dosn't exist
	logName := time.Now().Format("2006-01-02")
	CurrentLog = filepath.Join(LogDir, logName+".log")

	//O_CREATE - Create file if dosn't exist
	//O_EXCL - Except fail if it already exsits
	//O_Wrongly - Opens file for writing (which I don't need rn, but maybe later)

	//Atomic create
	f, err3 := os.OpenFile(CurrentLog, os.O_CREATE|os.O_EXCL, 644)
	if err3 != nil {
		if !os.IsExist(err3) {
			//Could not find or create log file
			internalMust("", err3)
		}
	}
	defer f.Close()
}
