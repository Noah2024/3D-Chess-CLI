// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

package load

import (
	"3DC/config"
	"3DC/util/logger"
	"3DC/util/must"
	"os"
	"unicode/utf8"

	"github.com/kelindar/bitmap"
)

type BoardState struct {
	FriendPieces        bitmap.Bitmap
	EnemyPieces         bitmap.Bitmap
	AllPieces           bitmap.Bitmap
	PieceLoadError      error
	FriendKing          bitmap.Bitmap
	FriendKingLoc       uint32
	SwapPawn            bool                     //Needed for checking becuase the pawn attacks differnet than it moves
	AllIndividualPieces map[string]bitmap.Bitmap //map[string]bitmap.Bitmap
	PieceInProcess      string                   //The piece currently being processed (used in protecting king in genMoves.restrictMoves)
}

// Loads dictionary, mapping display char to the bitmap corresponding with that piece
func LoadGame(fileLocation string) (data map[string]bitmap.Bitmap, err error) {

	result := make(map[string]bitmap.Bitmap)
	if _, err := os.Stat(fileLocation); os.IsNotExist(err) {
		logger.Warn("No game currently running, create one with '3DC game new'\n")
		return nil, err
	}
	entries := must.Must(os.ReadDir(fileLocation))

	for _, entry := range entries {
		if entry.IsDir() {
			//meta.bin is stored in /meta dir
			//This is so we can easily skip it when loading the piece bitmaps
			continue
		}
		file := must.Must(os.Open(fileLocation + "/" + entry.Name()))

		bm := must.Must(bitmap.ReadFrom(file))
		result[entry.Name()] = bm
	}
	return result, nil

}

// GenerateBoardState loops through all the bitmaps returned by LoadGame and using a reference piece determines friendly, & enemy pieces.
// Defined in load becuase it requires previously loaded data to function
func GenerateBoardState(loadedData map[string]bitmap.Bitmap, referencePiece string) (BoardState, error) {

	//Need to initalize the map
	rtn := BoardState{}
	rtn.AllIndividualPieces = make(map[string]bitmap.Bitmap) //Need to initalize this map
	rtn.PieceInProcess = referencePiece
	rtn.AllPieces.Grow(config.BoardSize - 1)
	rtn.FriendPieces.Grow(config.BoardSize - 1)

	r, _ := utf8.DecodeRuneInString(referencePiece)

	start, end := '♔', '♙'
	if r > '♙' {
		start, end = '♚', '♟'
	}

	for vis, bm := range loadedData {

		visAsRune, _ := utf8.DecodeRuneInString(vis)
		_, notEmpty := bm.Max()
		if !notEmpty {
			continue
		}

		if visAsRune >= start && visAsRune <= end {
			rtn.FriendPieces.Or(bm)
			if vis == "♔" || vis == "♚" {
				rtn.FriendKing = bm
			}
		} else {
			rtn.EnemyPieces.Or(bm)
		}

		// fmt.Println(vis)
		// fmt.Printf("BM FOR %064b\n", bm)
		rtn.FriendKingLoc, _ = rtn.FriendKing.Max() //Sets Friend King only once during move commnad
		rtn.AllIndividualPieces[vis] = bm
		rtn.AllPieces.Or(bm)
	}
	return rtn, nil
}

//Char Reference
// ChessPiece	Character	Unicode	 Go rune value (decimal)    Hex
// White King	♔	        U+2654	9812					   0x2654
// White Queen	♕	        U+2655	9813					   0x2655
// White Rook	♖			U+2656	9814					   0x2656
// White Bishop	♗			U+2657	9815					   0x2657
// White Knight	♘			U+2658	9816					   0x2658
// White Pawn	♙			U+2659	9817					   0x2659
// Black King	♚			U+265A	9818					   0x265A
// Black Queen	♛			U+265B	9819					   0x265B
// Black Rook	♜			U+265C	9820					   0x265C
// Black Bishop	♝			U+265D	9821					   0x265D
// Black Knight	♞			U+265E	9822					   0x265E
// Black Pawn	♟			U+265F	9823					   0x265F
