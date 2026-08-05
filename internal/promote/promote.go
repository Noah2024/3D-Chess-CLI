// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

// special contains functions relating to certian special movement conditions in chess
// These include, Pawn Promotion, Pawn Double Move, En-Pessent, & Castling
package promote

import (
	"3DC/config"
	"3DC/internal/game/load"
	"3DC/internal/game/save"
	"3DC/util/dataplane"
	"3DC/util/logger"
	"3DC/util/metadata"
	"fmt"
	"strings"

	"github.com/kelindar/bitmap"
)

//Maps the word of a piece to its given visual represetnation
//the rune of the character is subtracted by 6 to get the corresponding visual for black pieces

var wordToPiece = map[string]string{
	"knight": "♞",
	"bishop": "♝",
	"rook":   "♜",
	"queen":  "♛",
}

// Will attempt to promote a pawn to the type given
func AttemptPawnPromotion(promotionTo string, teamStr string) {
	var team bool
	var visFrom string
	teamStr = strings.ToLower(teamStr)

	switch teamStr {
	case "white":
		{
			team = true
			visFrom = "♟"
		}
	case "black":
		{
			team = false
			visFrom = "♙"
		}
	default:
		{
			logger.Error(fmt.Sprintf("Only valid teams are 'black' & 'white' not '%s'", teamStr))
			return
		}
	}

	promotionTarget, contains := wordToPiece[strings.ToLower(promotionTo)]
	if !contains {
		logger.Error(fmt.Sprintf("'%s' is not a valid piece to promote to", promotionTo))
		return
	}

	allLoadedPieces, meta, loadErr := load.LoadGame(config.CurrentGame)
	if loadErr != nil {
		logger.Error(fmt.Sprintf("Could not load board state (ensure you have  a 'CurrentGame' folder in your data directory) %v", loadErr))
		return
	}

	uLoc, present := CanPromotePawn(team, allLoadedPieces)

	if !present {
		logger.Error(fmt.Sprintf("There is no pawn avilable to promote on team '%s'", teamStr))
		return
	}

	if !team { //Set visual correctly
		promotionTarget = string([]rune(promotionTarget)[0] - 6)
	}

	PromotePawn(meta, uLoc, uLoc, promotionTarget, visFrom, allLoadedPieces[visFrom], allLoadedPieces[promotionTarget])

}

// Queries if the ANY pawn on this team can be promoted.
// Relies on the assumption that its not possible for more than one pawn to be promotable at any given time
func CanPromotePawn(Team bool, AllIndividualPieces map[string]bitmap.Bitmap) (uint32, bool) {
	var pieceToCheck string
	var planeToCheck bitmap.Bitmap
	if Team {
		pieceToCheck = "♟"
		planeToCheck = dataplane.WhitePromotionPlane
	} else {
		pieceToCheck = "♙"
		planeToCheck = dataplane.BlackPromotionPlane

	}

	allPiecesToCheck := AllIndividualPieces[pieceToCheck].Clone(nil)
	allPiecesToCheck.And(planeToCheck)

	uloc, present := allPiecesToCheck.Max()
	if present {
		return uloc, true
	}

	return 0, false
}

func PromotePawn(meta metadata.MetaData, uintLocFrom uint32, uintLocTo uint32, visTo string, visFrom string, bmFrom bitmap.Bitmap, bmTo bitmap.Bitmap) {

	bmFrom.Remove(uintLocFrom)
	bmTo.Set(uintLocTo)

	metadata.SaveMetaData(meta, config.CurrentGame)
	save.SavePieceType(visFrom, bmFrom)

	if visTo != "" {
		save.SavePieceType(visTo, bmTo)
	}
	logger.Info(fmt.Sprintf("Successfully Promoted pawn to '%s'", visTo))
}
