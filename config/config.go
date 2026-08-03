// Static config used to set certian game wide defaults at compile time.\n

// config is a package that provides static configuration for the 3DC application, including constants and variables that define the game's settings and directories.
// It is loaded at compile time and is used throughout the application to ensure consistent behavior and settings.
//
//	!!! After compilation its settings cannont be altered without recompiling the application. !!!
//
// It may be worth while to make a dynamic config file that can be altered by the user, but for now this is sufficent.
package config

import (
	"3DC/util/color"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

//Basic information application

const name = "3DC"
const year = "2026"
const version = "0.5.0"
const author = "Noah Yurasko"

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

// Board size is 8x8x8 THIS IS SIZE, NOT INDEX, SO 0-63 is valid range
const BoardSize uint32 = 512

// Layer size is 8x8 THIS IS SIZE, NOT INDEX, SO 0-63 is valid range
const LayerSize uint32 = 64

// Line size is 8 THIS IS SIZE, NOT INDEX, SO 0-7 is valid range
const LineSize uint32 = 8

// Sinle space size (please don't change this, I really honestly don't know what will happen if you do)
const SpaceSize uint32 = 1

// Direction which all data for the CLI application is stored. This includes game states, logs, and any other data that the application needs to persist between runs.
var DataDir string

// CurrentGame is the path to the current game state. It is used by the application to load and save the current game.
var CurrentGame string

// LogDir is the path to the directory where log files are stored. Each log file is named with the date it was created, and a new log file is created each day.
var LogDir string

// CurrentLog is the path to the current log file. It is used by the application to write log messages to the appropriate log file.
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

// Initalize all the requiset data directories for the applications function.
// Creaing them if necessary. This includes the data directory, the current game directory, and the log directory.
// It also creates a new log file for the current day if one does not already exist.
func init() {

	//Setting up the main user directory
	userDir := internalMust(os.UserConfigDir())

	//Data directory storing game states
	DataDir = filepath.Join(userDir, "3DC/DATA")
	err := os.MkdirAll(DataDir, 0o755)
	internalMust("", err)

	//Establishing where the Current Game is
	CurrentGame = filepath.Join(DataDir, "CurrentGame")

	//Defining where the path to the
	LogDir = filepath.Join(userDir, "3DC/LOG")
	err2 := os.MkdirAll(LogDir, 0o755)
	internalMust("", err2)

	//Creating the log file if it dosn't exist
	logName := time.Now().Format("2006-01-02")
	CurrentLog = filepath.Join(LogDir, logName+".log")

	//O_CREATE - Create file if dosn't exist
	//O_EXCL - Except fail if it already exsits
	//O_Wrongly - Opens file for writing (which I don't need rn, but maybe later)

	//Atomic create //Because otherwise the gap between opening the file and checking its existence
	//Could technically allow another program to create it
	f, err3 := os.OpenFile(CurrentLog, os.O_CREATE|os.O_EXCL, 0o755)
	if err3 != nil && f != nil {
		if !os.IsExist(err3) {
			//Could not find or create log file
			internalMust("", err3)
		}
	}
	defer f.Close()
}
