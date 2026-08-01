// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

// Package move contins the logic for moving pieces. Using the planes from dataplanes.go to generate valid moves then validates them for broken rules
// This is the largest and most compelex file in the project, due to it being the core gameplay element EVERYTHING ELSE is in support of the logic containted here
package move

import (
	"3DC/config"
	"3DC/internal/game/load"
	"3DC/internal/game/save"
	"3DC/internal/move/checking"
	"3DC/internal/move/genMoves"
	"3DC/util/bitutil"
	"3DC/util/logger"
	"3DC/util/metadata"

	"fmt"

	"github.com/kelindar/bitmap"
)

var BoardState load.BoardState

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
func ParseLoc(loc string) (uint32, int, int, int) {

	if len(loc) != 3 {
		logger.Error(fmt.Sprintf("Could not parse location '%v' - invalid length of string", loc))
	}
	x, z, y := int(loc[0]-'a'+1), int(loc[1]-'1'+1), int(loc[2]-'A'+1)

	return bitutil.VecToUint(x, y, z), x, y, z //bitutil.VecToUint(x, y, z)
}

// Determines peice type
// inputs uint32 location | outputs: string, bitmap.Bitmap (bitmap )
func PieceType(allLoadedPieces map[string]bitmap.Bitmap, loc uint32) (string, bitmap.Bitmap) {

	//Contains is simd vectorized, I don't feel the need to optimize this search
	for meta, bm := range allLoadedPieces {

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
	uLocFrom, fX, fY, fZ := ParseLoc(from)
	logger.Debug(fmt.Sprintf("Move called from %v to %v", from, to))

	allLoadedPieces, loadErr := load.LoadGame(config.CurrentGame)
	fmt.Println(allLoadedPieces)

	if loadErr != nil {
		logger.Error(fmt.Sprintf("Could not load board state (ensure you have  a 'CurrentGame' folder in your data directory) %v", BoardState.PieceLoadError))
		return
	}

	visFrom, bmFrom := PieceType(allLoadedPieces, uLocFrom)
	if visFrom == "" {
		logger.Error(fmt.Sprintf("Could not find piece at location %v", from))
		return
	}

	BoardState, err := load.GenerateBoardState(allLoadedPieces, visFrom)

	if err != nil {
		logger.Error(fmt.Sprintf("Error in determing pieces team %v", BoardState.PieceLoadError))
		return
	}

	uintLocTo, _, _, _ := ParseLoc(to)
	visTo, bmTo := PieceType(allLoadedPieces, uintLocTo)

	//visFrom encodes the type of piece, and thus the move function we use to generate all possible moves
	moveFunction := genMoves.MoveMap[visFrom]

	if moveFunction == nil {
		logger.Error(fmt.Sprintf("Unknown piece [%v]", visFrom))
		return
	}

	//Run entire checking system
	inCheck, inCheckMate, allowedKingMoves, savingKingMoves := checking.IsKingInCheck(BoardState)
	//If in checkmate, exit immiedatley
	if inCheckMate {
		logger.Error("Game is over, this pieces team is in checkmate!")
		return
	}

	allMoves := moveFunction(BoardState, uLocFrom, fX, fY, fZ)

	//Make sure this piece isn't protecting the king
	protectingMoves, kingDanger := checking.KingInDanger(BoardState, uLocFrom)
	fmt.Println("King Danger: ", kingDanger)
	if kingDanger { //I don't think this is right
		allMoves.And(protectingMoves)
		fmt.Printf("Protecting Moves: %064b\n", protectingMoves)
		// logger.Error(fmt.Sprintf("Piece %v is protecting its king!", visFrom))
	}

	//Kept for eventual need at debug
	// fmt.Printf("Piece Moving %s: ", visFrom)
	// fmt.Printf("Piece Being Taken %s ", visTo)

	//Restrict King Moves for king
	if visFrom == "♚" || visFrom == "♔" {
		fmt.Println("Restricting kings moves")
		fmt.Printf("Pos: %064b\n", BoardState.FriendKing)
		allMoves.And(allowedKingMoves)
	} else { //If king is in check can this piece take it out of check
		if inCheck {
			fmt.Println("Piece needs take the king out of check if possible")
			allMoves.And(savingKingMoves)
			_, canUnceck := allMoves.Max()
			if !canUnceck {
				logger.Error("Your king is in check!")
				return
			}
		}
	}
	//Restrict King Moves for king

	//Final check to see if piece can or can not move
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
