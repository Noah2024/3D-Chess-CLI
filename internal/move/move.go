// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

// Package move contins the logic for moving pieces. Using the planes from dataplanes.go to generate valid moves then validates them for broken rules
// This is the largest and most compelex file in the project, due to it being the core gameplay element EVERYTHING ELSE is in support of the logic containted here
package move

import (
	"3DC/config"
	"3DC/internal/game/load"
	"3DC/internal/game/save"
	"3DC/internal/move/genMoves"
	"3DC/util/bitutil"
	"3DC/util/logger"
	"3DC/util/metadata"

	"fmt"

	"github.com/kelindar/bitmap"
)

//Note to self for future development 7/24/2026
//Im noticing some concurrency issues when it comes to seperating out move generation with goroutines
//This was a given with the bishop but I've begun to see it elsewhere as well
//In the future I my plan is the leave move generation as sequential, then when I need to generate muliple moves
//Simply seperate those out and put them in parallel (such as for determining checking)

// BIG NOTE TO SELF 7/13/2026
// uintLoc is ONE indexed, meanwhile the bitmap  is ZERO indexed

//Taking input from
//X               Z               Y
//a b c d e f g - 1 2 3 4 5 6 7 8 - A B C D E F G

// parseLoc turns the user friendly notation input by the user (e.g., "a1A") to a uint32 index which represents that location in the bitmap
// inputs: string | outputs: uint32, x, y, z
func parseLoc(loc string) (uint32, int, int, int) {

	if len(loc) != 3 {
		logger.Error(fmt.Sprintf("Could not parse location '%v' - invalid length of string", loc))
	}
	x, z, y := int(loc[0]-'a'+1), int(loc[1]-'1'+1), int(loc[2]-'A'+1)

	return bitutil.VecToUint(x, y, z), x, y, z //bitutil.VecToUint(x, y, z)
}

// Determines peice type
// inputs uint32 location | outputs: string, bitmap.Bitmap (bitmap )
func pieceType(loc uint32) (string, bitmap.Bitmap) {

	allPieces, _ := load.LoadGame(config.CurrentGame)
	// fmt.Printf("All Pieces Alternate %064b\n", allPieces)

	//Contains is simd vectorized, I don't feel the need to optimize this search
	for meta, bm := range allPieces {

		if bm.Contains(loc) {
			// logger.Info(meta)
			return meta, bm
		}

	}
	return "", bitmap.Bitmap{}
}

// Move is the standard function to move one piece to another location.
// It takes the user friendly notation of the from and to locations, parses them into uint32 locations, and then checks if the move is valid for that piece type.
// If it is valid, it updates the bitmaps for both pieces and saves the game state.
// inputs: from string, to string | outputs: none
func MoveCommand(from string, to string) {
	uLocFrom, fX, fY, fZ := parseLoc(from)
	logger.Debug(fmt.Sprintf("Move called from %v to %v", from, to))

	visFrom, bmFrom := pieceType(uLocFrom)
	if visFrom == "" {
		logger.Error(fmt.Sprintf("Could not find piece at location %v", from))
		return
	}

	genMoves.FriendPieces, genMoves.EnemyPieces, genMoves.AllPieces, genMoves.BlackPawns, genMoves.PieceLoadError = load.GetFriendsAndEnemies(config.CurrentGame, visFrom)

	if genMoves.PieceLoadError != nil {
		logger.Error(fmt.Sprintf("Error in determing pieces team %v", genMoves.PieceLoadError))
		return
	}
	uintLocTo, _, _, _ := parseLoc(to)
	// fmt.Println("TO")
	// fmt.Printf("uLoc: %d | x: %d | y: %d | z: %d \n", uintLocTo, tX, tY, tZ)
	visTo, bmTo := pieceType(uintLocTo)

	//visFrom encodes the type of piece, and thus the move function we use to generate all possible moves
	moveFunction := genMoves.MoveMap[visFrom]

	if moveFunction == nil {
		logger.Error(fmt.Sprintf("Unknown piece [%v]", visFrom))
		return
	}

	allMoves := moveFunction(uLocFrom, fX, fY, fZ)

	//Kept for eventual need at debug
	// fmt.Printf("Piece Moving %s: ", visFrom)
	// fmt.Printf("Piece Being Taken %s ", visTo)

	if !(allMoves.Contains(uintLocTo)) {
		logger.Error(fmt.Sprintf("Piece %v cannot move in that way", visFrom))
		return
	}

	//Updates bitmap of piece being moved - does not validate if move is legal
	atomicMove(uLocFrom, uintLocTo, visTo, visFrom, bmFrom, bmTo)

	logger.Info("Piece Moved Successfully!")

}

// atomicMove instantly moves a piece from one location to another without any validation or checks.
// It is only used in practice from within a validated move funciton, and should only be used for debugging elsehwere
// So many variables are needed becuase no state is stored in the compiled binary itself, and thus the piece must be updated here for changes to take effect.
// inputs: from string, to string | outputs: none
func atomicMove(uintLocFrom uint32, uintLocTo uint32, visTo string, visFrom string, bmFrom bitmap.Bitmap, bmTo bitmap.Bitmap) {

	bmFrom.Remove(uintLocFrom)
	bmFrom.Set(uintLocTo)

	//Updates bitmap (if it exists) of piece being taken
	bmTo.Remove(uintLocTo)

	metadata.CreateSaveMetaData(config.CurrentGame)
	save.SavePieceType(visFrom, bmFrom)

	if visTo != "" {
		save.SavePieceType(visTo, bmTo)
	}
}
