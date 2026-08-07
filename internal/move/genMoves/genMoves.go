// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

// Generate moves contians all the functions which are used at runtime by the move function to generate all moves.
// Move generates the board state which is then passed into these functions, hence move needs to be run/loaded first otherwise
// BordState will have no data
package genMoves

import (
	"3DC/config"
	"3DC/internal/game/load"
	"3DC/internal/move/special"
	"3DC/util/bitutil"
	"3DC/util/dataplane"
	"sync"

	"github.com/kelindar/bitmap"
)

// global wait group is no bueno, this lead to raceconditions
// I will at some point take another pass at parallizing this code, but for now its aight
var wg sync.WaitGroup

// removeFriends takes a given bitmap of all possible moves and removes any moves that would result in a piece moving onto a friendly piece
func removeFriends(BoardState load.BoardState, allPossibleMoves bitmap.Bitmap) bitmap.Bitmap {
	// fmt.Printf("ALL POSSIBLE MOVES BEFORE %064b\n", allPossibleMoves)//For Debug
	result := allPossibleMoves.Clone(nil)
	result.Xor(BoardState.FriendPieces)
	result.And(allPossibleMoves)
	// fmt.Printf("ALL POSSIBLE MOVES AFTER %064b\n", allPossibleMoves)//For Debug
	return result
}

// Takes a given vector of a move and seperates it into right and left halves, before running checks on each half to determine if a piece is intersecting
// If no piece is present nothing happens and that half is jointed with the other, elsewise only the half upto and including the piece that is intersecting
// Friendly pieces are removed later
func restrictMoves(BoardState load.BoardState, curtPieceUintLoc uint32, moveLine bitmap.Bitmap) bitmap.Bitmap {
	// ==============================================
	// Seperate both directions of attack vector
	// ==============================================

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

	rightHitPieces := rightHalf.Clone(nil)        //bitmap.Bitmap
	rightHitPieces.And(BoardState.AllPieces)      //Contians all the pieces if any in the right half
	foundPerson, rtnRight := rightHitPieces.Min() //Gets first piece to be in line with attack

	if rtnRight == true { //If there is an enemy
		var newRight bitmap.Bitmap
		newRight.Grow(config.BoardSize - 1)
		newRight.Ones()
		newRight.Filter(func(x uint32) bool {
			return x <= foundPerson //Team offset handles the
		})

		rightHalf.And(newRight)
	}

	leftHitPieces := leftHalf.Clone(nil)
	leftHitPieces.And(BoardState.AllPieces)         //Contians all the pieces if any in the left half
	foundPersonLeft, rtnLeft := leftHitPieces.Max() //Gets first piece to be in line with attack

	if rtnLeft == true { //If there is an enemy
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
func generateRookMoves(BoardState load.BoardState, loc uint32, x int, y int, z int) bitmap.Bitmap { // Will parallelize with go rountine
	//Note to self, OK SO, the bitmaps when storing values store an ENTIRE BYTE at a time

	// logger.Debug(fmt.Sprintf("Generating all possible rook moves from :x: %d, y: %d, z: %d", x, y, z))
	forward := dataplane.XPlane[x].Clone(nil) //-2
	sideToSide := dataplane.YPlane[y].Clone(nil)
	upAndDown := dataplane.YPlane[y].Clone(nil)

	// wg.Go(
	// 	func() {
	forward.And(dataplane.ZPlane[z])
	forward = restrictMoves(BoardState, loc, forward) //x
	// 	},
	// )

	// wg.Go(
	// 	func() {
	sideToSide.And(dataplane.XPlane[x])                     //-2
	sideToSide = restrictMoves(BoardState, loc, sideToSide) //z
	// 	},
	// )

	// wg.Add(1)
	// wg.Go(
	// 	func() {
	upAndDown.And(dataplane.ZPlane[z])
	upAndDown = restrictMoves(BoardState, loc, upAndDown) //y
	// 	},
	// )

	// wg.Wait()
	forward.Or(upAndDown)
	forward.Or(sideToSide)
	// fmt.Printf("All Pieces %064b\n", allPieces)                //For Debug
	// fmt.Printf("All Allowed Move, forward)s %064b\n", forward) //For Debug
	forward = removeFriends(BoardState, forward)
	return forward
}

// generateBishopMoves contains the bitwise operations necessary to generate all possible moves for a bishop piece
// it takes x y and z integer cooridnates and outputs a size 511 bitmap all ones of which represent possible moves
// inputs: x, y, z int | outputs: bitmap.Bitmap
func generateBishopMoves(BoardState load.BoardState, loc uint32, x int, y int, z int) bitmap.Bitmap {
	// x, y, z = x-1, y-1, z-1 //positions must be zero indexed for indexing dataplanes

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
	realRight = restrictMoves(BoardState, loc, realRight)
	realLeft = restrictMoves(BoardState, loc, realLeft)

	//Again becuase of the descructive nature of .And the ordering of these operations is VERY IMPORTANT
	//Done in the wrong order one plane could be destroyed before it can be used in a different operation
	//This means that ALL of these must take place in a single thread

	bottomLeft.And(cardinalLeft)
	bottomLeft = restrictMoves(BoardState, loc, bottomLeft) //

	bottomRight.And(cardinalRight)
	bottomRight = restrictMoves(BoardState, loc, bottomRight)

	//After non descructivly using cardinal right and left above
	//We can use them descrustivly to get the top movements

	cardinalRight.And(dataplane.ZY45Plane[-z+y+7].Clone(nil)) // Top Right
	cardinalLeft.And(dataplane.ZY45Plane[-z+y+7].Clone(nil))  // Top Left
	cardinalRight = restrictMoves(BoardState, loc, cardinalRight)
	cardinalLeft = restrictMoves(BoardState, loc, cardinalLeft)

	//Then Or them all together to form the bishop's moves

	cardinalRight.Or(cardinalLeft)
	cardinalRight.Or(bottomLeft)
	cardinalRight.Or(bottomRight)
	cardinalRight.Or(realLeft)
	cardinalRight.Or(realRight)

	// fmt.Printf("All Pieces %064b\n", allPieces)
	// fmt.Printf("All Allowed Move, forward)s %064b\n", cardinalRight) //For Debug
	cardinalRight = removeFriends(BoardState, cardinalRight)
	return cardinalRight
}

// generateQueenMoves is just a logal OR between the bishop and rook moves, as the queen can move like both pieces
func generateQueenMoves(BoardState load.BoardState, loc uint32, x int, y int, z int) bitmap.Bitmap {
	// Must not zero index here becuase otherwise that would throw off the move generation
	// from the generators below
	var bishopMoves bitmap.Bitmap
	var rookMoves bitmap.Bitmap

	bishopMoves = generateBishopMoves(BoardState, loc, x, y, z)
	rookMoves = generateRookMoves(BoardState, loc, x, y, z)

	rookMoves.Or(bishopMoves)
	return rookMoves
}

// ==============================================
// General Note:
// The knight, pawn, and king are all hand coded and validated moves for the knight (becuase I can't use a cheeky lil bitmap for it)
// ==============================================

// generateKnightMoves is a hand coded and validated function for generating all possible moves for a knight piece
func generateKnightMoves(BoardState load.BoardState, loc uint32, x int, y int, z int) bitmap.Bitmap {
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
		// wg.Go(func() {
		X, Y, Z := x+comb[0], y+comb[1], z+comb[2]

		if X > 7 || Y > 7 || Z > 7 {
			continue
		}
		if X < 0 || Y < 0 || Z < 0 {
			continue
		}
		result.Set(bitutil.VecToUint(X, Y, Z))
		// })
	}

	// wg.Wait()
	result = removeFriends(BoardState, result)
	// fmt.Printf("All Pieces %064b\n", load.BS.AllPieces)
	// fmt.Printf("Result %064b\n", result)
	return result
}

// generateKingMoves is a hand coded and validated function for generating all possible moves for a king piece
func generateKingMoves(BoardState load.BoardState, loc uint32, x int, y int, z int) bitmap.Bitmap {
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
		// wg.Go(func() {
		X, Y, Z := x+comb[0], y+comb[1], z+comb[2]

		if X > 7 || Y > 7 || Z > 7 {
			continue
		}
		if X < 0 || Y < 0 || Z < 0 {
			continue
		}
		result.Set(bitutil.VecToUint(X, Y, Z))
		// })
	}

	// wg.Wait()
	result = removeFriends(BoardState, result)
	// fmt.Printf("All Pieces %064b\n", allPieces)
	// fmt.Printf("Result %064b\n", result)
	return result
}

// Hand coded and validated moves for the knight (becuase I can't use a cheeky lil bitmap for it)
func generatePawnMoves(BoardState load.BoardState, loc uint32, x int, y int, z int) bitmap.Bitmap {
	// x, y, z = x-1, y-1, z-1 //positions must be zero indexed for indexing da
	var result bitmap.Bitmap
	result.Grow(config.BoardSize - 1)

	zOffset := 1                                            //Using all black pieces to determine which way to move the pawn
	var canEnPessentFrom = dataplane.WhiteEnPessentPlane    //Plane from which can be enpessented from
	var enemyEnPessentable = BoardState.Meta.WhiteEnPessent //EnPessent move avilable if it exists

	if BoardState.AllIndividualPieces["♙"].Contains(loc) {
		zOffset = -1
		canEnPessentFrom = dataplane.BlackEnPessentPlane
		enemyEnPessentable = BoardState.Meta.BlackEnPessent
	}

	//AI used to speed up the processes of finding all valid permutations
	var normalMoves = [][]int{
		{0, 0, zOffset},
		{0, 1, zOffset},
		{0, -1, zOffset},
	}

	var doubleMoves = [][]int{
		{0, 0, zOffset * 2},
		{0, 1, zOffset * 2},
		{0, -1, zOffset * 2},
	}

	var attackingMoves = [][]int{
		{1, 0, zOffset},
		{1, 1, zOffset},
		{1, -1, zOffset},
		{-1, 0, zOffset},
		{-1, 1, zOffset},
		{-1, -1, zOffset},
	}

	canDouble := DoubleMove(BoardState, zOffset, loc)
	for i, comb := range normalMoves {
		// wg.Go(func() {
		legalDouble := true
		doubleX, doubleY, doubleZ := x+doubleMoves[i][0], y+doubleMoves[i][1], z+doubleMoves[i][2]

		X, Y, Z := x+comb[0], y+comb[1], z+comb[2]

		if X > 7 || Y > 7 || Z > 7 {
			continue
		}
		if X < 0 || Y < 0 || Z < 0 {
			continue
		}
		if doubleX > 7 || doubleY > 7 || doubleZ > 7 {
			legalDouble = false
		}
		if doubleX < 0 || doubleY < 0 || doubleZ < 0 {
			legalDouble = false
		}
		uin := bitutil.VecToUint(X, Y, Z)
		if !BoardState.AllPieces.Contains(uin) && !BoardState.SwapPawn { //Normal moves can only be made if there are no pieces there, ANY
			result.Set(uin)               //Plus if were not looking for checks
			if canDouble && legalDouble { //If the double move exists
				doubleuin := bitutil.VecToUint(doubleX, doubleY, doubleZ)
				result.Set(doubleuin)

			}
		}
	}
	// fmt.Println(BoardState.ReferencePiece)
	for _, comb := range attackingMoves {
		// wg.Go(func() {
		X, Y, Z := x+comb[0], y+comb[1], z+comb[2]

		if X > 7 || Y > 7 || Z > 7 {
			continue
		}
		if X < 0 || Y < 0 || Z < 0 {
			continue
		}
		uin := bitutil.VecToUint(X, Y, Z)

		if special.CouldMakeEnPessent(canEnPessentFrom, loc, X, int(enemyEnPessentable)) {
			result.Set(uin)
		}
		if BoardState.EnemyPieces.Contains(uin) || BoardState.SwapPawn { //Attacking moves can only be made if there ARE enemy pieces there
			result.Set(uin)
		}
		// })
	}

	// wg.Wait()
	return result
}

// Dynamically determines if a given piece can double move
func DoubleMove(BoardState load.BoardState, zOffSet int, locTocheck uint32) bool {
	var planeToCheck bitmap.Bitmap
	if zOffSet == -1 {
		planeToCheck = dataplane.WhiteDoubleMovePlane
	} else {
		planeToCheck = dataplane.BlackDoubleMovePlane
	}

	//Check if this piece is on the double move plane &
	if planeToCheck.Contains(locTocheck) {
		return true
	}

	return false
}

// moveMap matches a pieces visual representation to the function that generates all possible moves for that piece
// inputs: string | outputs: function(int, int, int) bitmap.Bitmap
var MoveMap = map[string]func(load.BoardState, uint32, int, int, int) bitmap.Bitmap{
	"♙": generatePawnMoves,
	"♘": generateKnightMoves,
	"♗": generateBishopMoves,
	"♖": generateRookMoves,
	"♕": generateQueenMoves,
	"♔": generateKingMoves,
	"♟": generatePawnMoves,
	"♞": generateKnightMoves,
	"♝": generateBishopMoves,
	"♜": generateRookMoves,
	"♛": generateQueenMoves,
	"♚": generateKingMoves,
}
