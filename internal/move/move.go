// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

// Package move contins the logic for moving pieces. Using the planes from dataplanes.go to generate valid moves then validates them for broken rules
// This is the largest and most compelex file in the project, due to it being the core gameplay element EVERYTHING ELSE is in support of the logic containted here
package move

import (
	"3DC/config"
	"3DC/internal/game/load"
	"3DC/internal/game/save"
	"3DC/util/bitutil"
	"3DC/util/dataplane"
	"3DC/util/logger"
	"3DC/util/metadata"
	"sync"

	"fmt"

	"github.com/kelindar/bitmap"
)

var friendPieces bitmap.Bitmap
var enemyPieces bitmap.Bitmap
var allPieces bitmap.Bitmap
var pieceLoadError error
var blackPawns bitmap.Bitmap //Used to determine direction of pawns move dynamically at runtime
var wg sync.WaitGroup

//Note to self for future development 7/24/2026
//Im noticing some concurrency issues when it comes to seperating out move generation with goroutines
//This was a given with the bishop but I've begun to see it elsewhere as well
//In the future I my plan is the leave move generation as sequential, then when I need to generate muliple moves
//Simply seperate those out and put them in parallel (such as for determining checking)

// BIG NOTE TO SELF 7/13/2026
// uintLoc is ONE indexed, meanwhile the bitmap  is ZERO indexed

func removeFriends(allPossibleMoves bitmap.Bitmap) bitmap.Bitmap {
	// fmt.Printf("ALL POSSIBLE MOVES BEFORE %064b\n", allPossibleMoves)//For Debug
	result := allPossibleMoves.Clone(nil)
	result.Xor(friendPieces)
	result.And(allPossibleMoves)
	// fmt.Printf("ALL POSSIBLE MOVES AFTER %064b\n", allPossibleMoves)//For Debug
	return result
}

// Takes a given vector of a move and seperates it into right and left halves, before running checks on each half to determine if a piece is intersecting
// If no piece is present nothing happens and that half is jointed with the other, elsewise only the half upto and including the piece that is intersecting
// Friendly pieces are removed later
func restrictMoves(curtPieceUintLoc uint32, moveLine bitmap.Bitmap) bitmap.Bitmap {
	// ==============================================
	// Seperate both directions of attack vector
	// ==============================================

	// fmt.Printf("Move line %064b \n", moveLine)

	var leftMask bitmap.Bitmap
	var rightHalf bitmap.Bitmap
	rightHalf.Grow(config.BoardSize - 1) //Grow to normal size
	rightHalf.Ones()
	leftMask.Grow(config.BoardSize - 1)
	leftMask.Ones()
	leftMask.Filter(func(x uint32) bool {
		return x < curtPieceUintLoc+1 //Not sure why this has to be +1 exactly, but otherwise movement lines will be slighly off
	})

	rightHalf.Xor(leftMask) //Leaves ones only on the right half
	rightHalf.And(moveLine) //Cancels out all the movement from the left

	// Calculate left half
	var leftHalf bitmap.Bitmap
	leftHalf.Grow(config.BoardSize - 1) //Filled left half with ones
	leftHalf.Ones()
	leftHalf.Filter(func(x uint32) bool {
		return x < curtPieceUintLoc
	})
	leftHalf.And(moveLine) //Only ones left will be one in both

	// ==============================================
	// Determine friend or foe and mask accordingly
	// ==============================================

	// fmt.Printf("Right Before %064b \n", rightHalf)
	// fmt.Printf("Left Before %064b \n", leftHalf)

	rightHitPieces := rightHalf.Clone(nil)   //bitmap.Bitmap
	rightHitPieces.And(allPieces)            //Contians all the pieces if any in the right half
	foundPerson, rtn := rightHitPieces.Min() //Gets first piece to be in line with attack

	if rtn == true { //If there is an enemy
		var newRight bitmap.Bitmap
		newRight.Grow(config.BoardSize - 1)
		newRight.Ones()
		newRight.Filter(func(x uint32) bool {
			return x <= foundPerson //Team offset handles the
		})

		rightHalf.And(newRight)
	}

	leftHitPieces := leftHalf.Clone(nil)
	leftHitPieces.And(allPieces)                //Contians all the pieces if any in the left half
	foundPersonLeft, rtn := leftHitPieces.Max() //Gets first piece to be in line with attack

	if rtn == true { //If there is an enemy
		var newLeft bitmap.Bitmap
		var newLeftMask bitmap.Bitmap //Supposed to be all zeros
		newLeft.Grow(config.BoardSize - 1)
		newLeft.Ones()
		newLeftMask.Grow(config.BoardSize - 1) //Can only filter through set bits)
		newLeftMask.Ones()

		newLeftMask.Filter(func(x uint32) bool {
			return x >= foundPersonLeft
		})

		newLeft.And(newLeftMask)
		leftHalf.And(newLeft)
	}

	// ==============================================
	// Combine into final bitmap and return
	// ==============================================

	rightHalf.Or(leftHalf)
	return rightHalf
}

//Strange thing, becuase of how fucking fast this bitmap library is, I think the overhead on the creation of a wait group
// Is actually SLOWER than just doing all the operations sequentually, though given im only working with one piece and a few test cases
// And and operating on a differnece of a few miliseconds, im going to wait until after I implment more pieces and thus more test cases
//Until then this old version will stay here until I can determine a real runtime benefit to using the waitgroup
// func generateRookMoves(loc uint32, x int, y int, z int) bitmap.Bitmap { // Will parallelize with go rountine
// 	//Note to self, OK SO, the bitmaps when storing values store an ENTIRE BYTE at a time

// 	// logger.Debug(fmt.Sprintf("Generating all possible rook moves from :x: %d, y: %d, z: %d", x, y, z))
// 	forward := dataplane.XPlane[x-1].Clone(nil) //-2
// 	forward.And(dataplane.ZPlane[z-1])
// 	forward = restrictMoves(loc, forward) //x

// 	sideToSide := dataplane.YPlane[y-1].Clone(nil)
// 	sideToSide.And(dataplane.XPlane[x-1])       //-2
// 	sideToSide = restrictMoves(loc, sideToSide) //z

// 	upAndDown := dataplane.YPlane[y-1].Clone(nil)
// 	upAndDown.And(dataplane.ZPlane[z-1])
// 	upAndDown = restrictMoves(loc, upAndDown) //y

// 	forward.Or(upAndDown)
// 	forward.Or(sideToSide)
// 	// fmt.Printf("All Pieces %064b\n", allPieces) //For Debug
// 	// fmt.Printf("All Allowed Moves %064b\n", forward) //For Debug
// 	return forward
// }

// Go Routine Version

// generateRookMoves contains the bitwise operations necessary to generate all possible moves for a rook piece
// it takes x y and z integer cooridnates and outputs a size 511 bitmap all ones of which represent possible moves
// inputs: x, y, z int | outputs: bitmap.Bitmap
func generateRookMoves(loc uint32, x int, y int, z int) bitmap.Bitmap { // Will parallelize with go rountine
	//Note to self, OK SO, the bitmaps when storing values store an ENTIRE BYTE at a time

	// logger.Debug(fmt.Sprintf("Generating all possible rook moves from :x: %d, y: %d, z: %d", x, y, z))
	forward := dataplane.XPlane[x-1].Clone(nil) //-2
	sideToSide := dataplane.YPlane[y-1].Clone(nil)
	upAndDown := dataplane.YPlane[y-1].Clone(nil)

	wg.Go(
		func() {
			forward.And(dataplane.ZPlane[z-1])
			forward = restrictMoves(loc, forward) //x
		},
	)

	wg.Go(
		func() {
			sideToSide.And(dataplane.XPlane[x-1])       //-2
			sideToSide = restrictMoves(loc, sideToSide) //z
		},
	)

	// wg.Add(1)
	wg.Go(
		func() {
			upAndDown.And(dataplane.ZPlane[z-1])
			upAndDown = restrictMoves(loc, upAndDown) //y
		},
	)

	wg.Wait()
	forward.Or(upAndDown)
	forward.Or(sideToSide)
	// fmt.Printf("All Pieces %064b\n", allPieces)                //For Debug
	// fmt.Printf("All Allowed Move, forward)s %064b\n", forward) //For Debug
	forward = removeFriends(forward)
	return forward
}

func generateBishopMove(loc uint32, x int, y int, z int) bitmap.Bitmap {
	x, y, z = x-1, y-1, z-1 //positions must be zero indexed for indexing dataplanes

	//The indexing for each of these is computed using a formula based on how they were computed, go to dataplanes to check
	//And work it out for yourself until I have time to better document it
	// XY45NegPlane := dataplane.XY45NegPlane[-x-y+14].Clone(nil) //-14
	// XY45Plane := dataplane.XY45Plane[-x+y+7].Clone(nil)
	// XZ45NegPlane := dataplane.XZ45NegPlane[-x-z+14].Clone(nil) //
	// XZ45PosPlane := dataplane.XZ45PosPlane[-x+z+7].Clone(nil)  //
	// ZY45NegPlane := dataplane.ZY45NegPlane[-z-y+14].Clone(nil)
	// ZY45Plane := dataplane.ZY45Plane[-z+y+7].Clone(nil)

	// Cardinal Right and left are the direct diagnols you would see WITHOUT any dimension in the Y,
	// Essentially they are the diagnols you would see on a normal chess board looking down

	cardinalRight := dataplane.XZ45PosPlane[-x+z+7].Clone(nil)
	cardinalLeft := dataplane.XZ45NegPlane[-x-z+14].Clone(nil)

	//Real Right and Left are those cardinal directions cast onto the y axis where the piece is
	realRight := dataplane.YPlane[y].Clone(nil)
	realLeft := dataplane.YPlane[y].Clone(nil)

	bottomLeft := dataplane.XY45Plane[-x+y+7].Clone(nil)      // Side Left
	bottomRight := dataplane.XY45NegPlane[-x-y+14].Clone(nil) //-14

	//Beuase the cardinal right and left cut thorugh all dimensions however we need to
	//Use our current Y layer to get them AND becuase of the descructive nature of the .And operation
	//We need to do so in new bitmaps

	realRight.And(cardinalRight)
	realLeft.And(cardinalLeft)
	realRight = restrictMoves(loc, realRight)
	realLeft = restrictMoves(loc, realLeft)

	//Again becuase of the descructive nature of .And the ordering of these operations is VERY IMPORTANT
	//Done in the wrong order one plane could be destroyed before it can be used in a different operation
	//This means that ALL of these must take place in a single thread

	bottomLeft.And(cardinalLeft)
	bottomLeft = restrictMoves(loc, bottomLeft) //

	bottomRight.And(cardinalRight)
	bottomRight = restrictMoves(loc, bottomRight)

	//After non descructivly using cardinal right and left above
	//We can use them descrustivly to get the top movements

	cardinalRight.And(dataplane.ZY45Plane[-z+y+7].Clone(nil)) // Top Right
	cardinalLeft.And(dataplane.ZY45Plane[-z+y+7].Clone(nil))  // Top Left
	cardinalRight = restrictMoves(loc, cardinalRight)
	cardinalLeft = restrictMoves(loc, cardinalLeft)

	//Then Or them all together to form the bishop's moves

	cardinalRight.Or(cardinalLeft)
	cardinalRight.Or(bottomLeft)
	cardinalRight.Or(bottomRight)
	cardinalRight.Or(realLeft)
	cardinalRight.Or(realRight)

	// fmt.Printf("All Pieces %064b\n", allPieces)
	// fmt.Printf("All Allowed Move, forward)s %064b\n", cardinalRight) //For Debug
	cardinalRight = removeFriends(cardinalRight)
	return cardinalRight
}

func generateQueenMove(loc uint32, x int, y int, z int) bitmap.Bitmap {
	// Must not zero index here becuase otherwise that would throw off the move generation
	// from the generators below
	var bishopMoves bitmap.Bitmap
	var rookMoves bitmap.Bitmap

	bishopMoves = generateBishopMove(loc, x, y, z)
	rookMoves = generateRookMoves(loc, x, y, z)

	rookMoves.Or(bishopMoves)
	return rookMoves
}

// Hand coded and validated moves for the knight (becuase I can't use a cheeky lil bitmap for it)
func generateKnightMove(loc uint32, x int, y int, z int) bitmap.Bitmap {
	// x, y, z = x-1, y-1, z-1 //positions must be zero indexed for indexing da

	var result bitmap.Bitmap
	result.Grow(config.BoardSize - 1)

	//AI used to speed up the processes of finding all valid permutations
	var allCombs = [][]int{
		// 0 in the middle (original orientation)
		{2, 0, 1},
		{2, 0, -1},
		{1, 0, 2},
		{1, 0, -2},
		{-1, 0, 2},
		{-1, 0, -2},
		{-2, 0, 1},
		{-2, 0, -1},

		// 0 in the first position
		{0, 2, 1},
		{0, 2, -1},
		{0, 1, 2},
		{0, 1, -2},
		{0, -1, 2},
		{0, -1, -2},
		{0, -2, 1},
		{0, -2, -1},

		// 0 in the third position
		{2, 1, 0},
		{2, -1, 0},
		{1, 2, 0},
		{1, -2, 0},
		{-1, 2, 0},
		{-1, -2, 0},
		{-2, 1, 0},
		{-2, -1, 0},
	}

	for _, comb := range allCombs {
		wg.Go(func() {
			X, Y, Z := x+comb[0], y+comb[1], z+comb[2]

			if X > 8 || Y > 8 || Z > 8 {
				return
			}
			if X < 1 || Y < 1 || Z < 1 {
				return
			}
			result.Set(bitutil.VecToUint(X, Y, Z))
		})
	}

	wg.Wait()
	result = removeFriends(result)
	// fmt.Printf("All Pieces %064b\n", allPieces)
	// fmt.Printf("Result %064b\n", result)
	return result
}

// Hand coded and validated moves for the knight (becuase I can't use a cheeky lil bitmap for it)
func generateKingMove(loc uint32, x int, y int, z int) bitmap.Bitmap {
	// x, y, z = x-1, y-1, z-1 //positions must be zero indexed for indexing da

	var result bitmap.Bitmap
	result.Grow(config.BoardSize - 1)

	//AI used to speed up the processes of finding all valid permutations
	var allCombs = [][]int{
		// 0 in the middle
		{1, 0, 1},
		{1, 0, -1},
		{-1, 0, 1},
		{-1, 0, -1},

		// 0 in the first position
		{0, 1, 1},
		{0, 1, -1},
		{0, -1, 1},
		{0, -1, -1},

		// 0 in the third position
		{1, 1, 0},
		{1, -1, 0},
		{-1, 1, 0},
		{-1, -1, 0},

		// All
		{1, 1, 1},
		{1, 1, -1},
		{1, -1, 1},
		{1, -1, -1},
		{-1, 1, 1},
		{-1, 1, -1},
		{-1, -1, 1},
		{-1, -1, -1},

		//Permutations w zero
		{1, 0, 0},
		{-1, 0, 0},
		{0, 1, 0},
		{0, -1, 0},
		{0, 0, 1},
		{0, 0, -1},
	}

	for _, comb := range allCombs {
		wg.Go(func() {
			X, Y, Z := x+comb[0], y+comb[1], z+comb[2]

			if X > 8 || Y > 8 || Z > 8 {
				return
			}
			if X < 1 || Y < 1 || Z < 1 {
				return
			}
			result.Set(bitutil.VecToUint(X, Y, Z))
		})
	}

	wg.Wait()
	result = removeFriends(result)
	// fmt.Printf("All Pieces %064b\n", allPieces)
	// fmt.Printf("Result %064b\n", result)
	return result
}

// Hand coded and validated moves for the knight (becuase I can't use a cheeky lil bitmap for it)
func generatePawnMove(loc uint32, x int, y int, z int) bitmap.Bitmap {
	// x, y, z = x-1, y-1, z-1 //positions must be zero indexed for indexing da

	var result bitmap.Bitmap
	result.Grow(config.BoardSize - 1)

	zOffset := 1
	if blackPawns.Contains(loc) {
		zOffset = -1
	}

	//AI used to speed up the processes of finding all valid permutations
	var normalMoves = [][]int{
		{0, 0, zOffset},
		{0, 1, zOffset},
		{0, -1, zOffset},
	}

	var attackingMoves = [][]int{
		{1, 0, zOffset},
		{1, 1, zOffset},
		{1, -1, zOffset},
		{-1, 0, zOffset},
		{-1, 1, zOffset},
		{-1, -1, zOffset},
	}

	for _, comb := range normalMoves {
		wg.Go(func() {
			X, Y, Z := x+comb[0], y+comb[1], z+comb[2]

			if X > 8 || Y > 8 || Z > 8 {
				return
			}
			if X < 1 || Y < 1 || Z < 1 {
				return
			}
			uin := bitutil.VecToUint(X, Y, Z)
			if !allPieces.Contains(uin) { //Normal moves can only be made if there are no pieces there, ANY
				result.Set(uin)
			}
		})
	}

	for _, comb := range attackingMoves {
		wg.Go(func() {
			X, Y, Z := x+comb[0], y+comb[1], z+comb[2]

			if X > 8 || Y > 8 || Z > 8 {
				return
			}
			if X < 1 || Y < 1 || Z < 1 {
				return
			}
			uin := bitutil.VecToUint(X, Y, Z)
			if enemyPieces.Contains(uin) { //Attacking moves can only be made if there ARE enemy pieces there
				result.Set(uin)
			}
		})
	}

	wg.Wait()
	// fmt.Printf("All Pieces %064b\n", allPieces)
	// fmt.Printf("Result %064b\n", result)
	return result
}

// moveMap matches a pieces visual representation to the function that generates all possible moves for that piece
// inputs: string | outputs: function(int, int, int) bitmap.Bitmap
var moveMap = map[string]func(uint32, int, int, int) bitmap.Bitmap{
	"♙": generatePawnMove,
	"♘": generateKnightMove,
	"♗": generateBishopMove,
	"♖": generateRookMoves,
	"♕": generateQueenMove,
	"♔": generateKingMove,
	"♟": generatePawnMove,
	"♞": generateKnightMove,
	"♝": generateBishopMove,
	"♜": generateRookMoves,
	"♛": generateQueenMove,
	"♚": generateKingMove,
}

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

	friendPieces, enemyPieces, allPieces, blackPawns, pieceLoadError = load.GetFriendsAndEnemies(config.CurrentGame, visFrom)

	if pieceLoadError != nil {
		logger.Error(fmt.Sprintf("Error in determing pieces team %v", pieceLoadError))
		return
	}
	uintLocTo, _, _, _ := parseLoc(to)
	// fmt.Println("TO")
	// fmt.Printf("uLoc: %d | x: %d | y: %d | z: %d \n", uintLocTo, tX, tY, tZ)
	visTo, bmTo := pieceType(uintLocTo)

	//visFrom encodes the type of piece, and thus the move function we use to generate all possible moves
	moveFunction := moveMap[visFrom]

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
