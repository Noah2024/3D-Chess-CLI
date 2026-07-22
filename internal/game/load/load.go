// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

package load

import (
	"3DC/util/logger"
	"3DC/util/must"
	"os"
	"unicode/utf8"

	"github.com/kelindar/bitmap"
)

// Loads dictionary mapping display char to the bitmap corresponding with that piece
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

// GetFriendsAndEnemies loops through all the bitmaps returned by LoadGame and using a reference piece determines friendly, & enemy pieces.
// Defined in load becuase it is indepdendly loading in the board to make this determination
func GetFriendsAndEnemies(fileLocation string, referencePiece string) (bitmap.Bitmap, bitmap.Bitmap, bitmap.Bitmap, error) {
	// result := make(map[string]bitmap.Bitmap)
	if _, err := os.Stat(fileLocation); os.IsNotExist(err) {
		logger.Warn("No game currently running, create one with '3DC game new'\n")
		return nil, nil, nil, err
	}
	entries := must.Must(os.ReadDir(fileLocation))

	var friendPieces bitmap.Bitmap
	var enemyPieces bitmap.Bitmap
	var allPieces bitmap.Bitmap

	r, _ := utf8.DecodeRuneInString(referencePiece)

	start, end := '♔', '♙'
	if r > '♙' {
		start, end = '♚', '♟'
	}

	for _, entry := range entries {
		if entry.IsDir() {
			//meta.bin is stored in /meta dir
			//This is so we can easily skip it when loading the piece bitmaps
			continue
		}
		vis := entry.Name()
		visAsRune, _ := utf8.DecodeRuneInString(vis)

		file := must.Must(os.Open(fileLocation + "/" + entry.Name()))
		bm := must.Must(bitmap.ReadFrom(file))

		if visAsRune >= start && visAsRune <= end {
			friendPieces.Or(bm)
		} else {
			enemyPieces.Or(bm)
		}
		// fmt.Println(vis)
		// fmt.Printf("BM FOR %064b\n", bm)
		allPieces.Or(bm)

	}
	return friendPieces, enemyPieces, allPieces, nil

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
