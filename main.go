package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

// HANDLING CLI CREATION AND SHIT
type CLI struct {
	userInpt string
	active   bool
}

func (cli *CLI) handleClose() {
	cli.active = false
}

func (cli *CLI) handleCommand(cmd string, game *GAME) {
	parsedCmd := strings.Fields(cmd)

	switch parsedCmd[0] {
	case "exit", "e":
		cli.handleClose()
	case "ls":
		game.viewBoard()
	case "mvp":
		if len(parsedCmd) == 2 {
			game.move(parsedCmd[1])
		} else {
			fmt.Printf("move command takes exactly one argument '%s'\n", cmd)
		}
	case "get":
		if len(parsedCmd) == 2 {
			game.getSpace(parsedCmd[1])
		} else {
			fmt.Printf("get command takes exactly one argument '%s'\n", cmd)
		}
	// case: "debug":
	// 	game.
	default:
		fmt.Printf("Unknown Command '%s'\n", cmd)
	}
}

// CLI CREATION AND SHIT DONE
// Note to self: doing anything in which the TOTAL number of pieces is different than 32, is not a great idea, FOR NOW
// STARTING GAME LOGIC
type intvec struct {
	x int
	y int
	z int
}

type GAME struct {
	turn      int
	allPieces [32]PIECE
	boardSize intvec
	board     [8][8]rune
	colorMap  [8][8]rune
}

type PIECE struct {
	name rune
	pos  intvec
}

func (game *GAME) initPieces() {
	fmt.Println("Init Pieces")
	//Prolly a more efficent way to initalize the pieces and board at the same time
	//But for now this should work
	for i := 0 + 1; i < game.boardSize.x; i += 1 {
		for j := 0 + 1; j < game.boardSize.y; j += 1 {
			game.board[i][j] = 0
		}
	}
	x, y, i := 0, 0, 0 //I is the piece index
	for _, character := range "RNBQKBNR/PPPPPPPP/8/8/8/8/pppppppp/rnbqkbnr" {
		if character == '/' {
			y++
			x = 0
			continue
		} else if unicode.IsDigit(character) {
			x += int(character - '0')
		} else {
			game.allPieces[i] = PIECE{
				name: character,
				pos:  intvec{x, y, 0},
			}
			game.board[y][x] = character
			i++
		}
		x += 1
	}
}

// VIEW BOARD MUST BE CHANGED TO MOVE TO 3D
func (game *GAME) viewBoard() {
	fmt.Println("    A B C D E F G H   ")
	for i, rank := range game.board {
		fmt.Print(string('0'+i) + " | ")
		for j, space := range rank {
			//I don't like preforming this many checks when rendering the board, but im not sure what other alternatives I have
			if game.colorMap[j][i] == '1' {
				fmt.Printf("\x1b[32m%c\x1b[0m", space+'0')
				return
			}
			if space == 0 {
				fmt.Print("0") //Must hadle 0s differently cuase otherwise they print as blank spaces string(space + '0')
			} else {
				fmt.Print(string(space))
			}
			fmt.Print(" ")
		}
		fmt.Print("|")
		fmt.Println()
	}

}

// NOTE TO SELF< MUST CHANGE BOUNDRIES ON FUNCS BEFORE MAKING 3D
func (game *GAME) decodeGmPos(pos string) (int, int, error) {
	if len(pos) != 2 {
		return 0, 0, fmt.Errorf("invalid postiiton length")
	}
	col := rune(pos[0])
	row := rune(pos[1])

	if col < 'a' || col > 'z' {
		return 0, 0, fmt.Errorf("invlaid colum")
	}
	if row < '0' || row > '9' {
		return 0, 0, fmt.Errorf("invlaid row")
	}

	x := int(row - '0')
	y := int(col - 'a')
	return x, y, nil
}

func (game *GAME) getTeam(piece rune) bool {
	if unicode.ToUpper(piece) == piece {
		return true
	} else {
		return false
	}
}

func (game *GAME) getSpace(pos string) {
	game.decodeGmPos(pos)
	x, y, error1 := game.decodeGmPos(pos)
	if error1 != nil {
		fmt.Printf("Error in reading space %s \n", pos)
	} else {
		fmt.Printf("Space X:%s - Y:%s contains %c \n", x, y, game.board[x][y])
	}
}

func (game *GAME) addHilightSpace(x int, y int) {
	//If I want to I could later add differnet colors, but this is mostly just for debug purposes
	// x, y, error1 := game.decodeGmPos(pos)
	// if error1 != nil {
	// 	fmt.Println("Could not find space to highlight")
	// }
	game.colorMap[x][y] = 1 //"\x1b[32m%c\x1b[0m"
}

func (game *GAME) preventMoveThrough(start intvec, end intvec, xDif int, yDif int) bool {
	//May need to modify for the bishop, or at least do more testing later
	fmt.Println(game.colorMap)
	//Need to offset so it dosn't prevent moves with itself
	//However this needs to be offset in the direction of motion, hence the need for xDif and yDif
	var xOffSet int
	var yOffSet int
	if xDif == 0 {
		xOffSet = 0
	} else if xDif > 0 {
		xOffSet = 1
	} else {
		xOffSet = -1
	}
	if yOffSet == 0 {
		yOffSet = 0
	} else if yDif > 0 {
		yOffSet = 1
	} else {
		yOffSet = -1
	}
	fmt.Println("Offsets are: ", xOffSet, yOffSet)
	for i := start.x + xOffSet; i <= end.x; i += 1 {
		for j := start.y + yOffSet; j <= end.y; j += 1 {
			game.addHilightSpace(i, j)
			if game.board[i][j] != 0 {
				fmt.Printf("PROBLEM AT %d, %d \n", j, i)
				fmt.Print(string(game.board[i][j]))
				return false
			}
			fmt.Printf("HILIGHT %c \n", game.colorMap[i][j])
		}
	}
	game.viewBoard()
	return true
}

// Using != as exclusive or in this function
func (game *GAME) validateMoves(pieceToMove rune, mvFrom intvec, mvTo intvec) bool {
	team := game.getTeam(pieceToMove)
	xDif := mvTo.x - mvFrom.x
	yDif := mvTo.y - mvFrom.y

	switch unicode.ToUpper(pieceToMove) {
	case 'P':
		fmt.Println("Moving Pawn")
		dir := 1
		if team == false {
			dir = -1
		}
		fmt.Println(dir)
	case 'N':
		fmt.Println("Moving Knigt")
	case 'B':
		fmt.Println("Moving Bishop")
		if !((xDif != 0) == (yDif != 0)) {
			fmt.Println("Invalid Bishop Move")
			fmt.Println("The Bishop Cannont Move like that")
			return false
		}
		if !game.preventMoveThrough(mvFrom, mvTo, xDif, yDif) {
			fmt.Println("Invalid Bishop Move")
			fmt.Println("Can not move thorugh another piece")
			return false
		}
	case 'R':
		fmt.Println("Moving Rook")
		if !((xDif == 0) != (yDif == 0)) {
			fmt.Println("Invalid Rook Move")
			fmt.Println("The Rook Cannont Move like that")
			return false
		}
		if !game.preventMoveThrough(mvFrom, mvTo, xDif, yDif) {
			fmt.Println("Invalid Rook Move")
			fmt.Println("Can not move thorugh another piece")
			return false
		}
	case 'Q':
		fmt.Println("Moving Queen")
	case 'K':
		fmt.Println("Moving King")
	}

	return true
}

func (game *GAME) move(cmd string) {
	//Consider moving this data validation into teh decodeGmPos function
	if !strings.Contains(cmd, "-") {
		fmt.Println("IMPROPER FORMAT")
		return
	}
	parts := strings.Split(cmd, "-")
	mvFrmX, mvFrmY, error1 := game.decodeGmPos(parts[0])
	mvToX, mvToY, error2 := game.decodeGmPos(parts[1])

	if error1 != nil || error2 != nil {
		fmt.Println("IMPROPER FORMAT")
		return
	}

	fmt.Printf("MOVING FROM X: %d Y: %d \n", mvFrmX, mvFrmY)
	fmt.Printf("TO X: %d Y: %d \n", mvToX, mvToY)

	pieceToMove := game.board[mvFrmX][mvFrmY]
	if pieceToMove == 0 {
		fmt.Println("No piece exists")
		return
	}
	if game.validateMoves(pieceToMove, intvec{mvFrmX, mvFrmY, 0}, intvec{mvToX, mvToY, 0}) {
		game.board[mvToX][mvToY] = pieceToMove
		game.board[mvFrmX][mvFrmY] = 0
		fmt.Printf("Moved Piece %+v to X: %s, Y:%s \n", pieceToMove, mvToX, mvToY)

	} else {
		fmt.Println("Could Not Move Piece: Check above for reason")
	}

	//TO CHANGE
	//Create one validate move function, which will have all the nasty swtich statemnts
	//Then Everything here can be functinally simpler

	// switch unicode.ToUpper(pieceToMove.name){
	// case 'P':
	// 	fmt.Prnitln("Moving Pawn")
	// case 'N':
	// 	fmt.Prnitln("Moving Knigt")
	// case 'B':
	// 	fmt.Println("Moving Bishop")
	// case 'R':
	// 	fmt.Println("Moving Rook")
	// case 'Q':
	// 	fmt.Println("Moving Queen")
	// case 'K':
	// 	fmt.Println("Moving King")
	// }
}

// GAME LOGIC DONE
// START MISC BITS
func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("\033[H\033[2J")

	cli := CLI{
		userInpt: "",
		active:   true,
	}

	game := GAME{
		turn:      0,
		boardSize: intvec{8, 8, 0},
	}
	game.initPieces()

	// fmt.Println("test")
	fmt.Println("----------")
	fmt.Println("3D Chess CLI")
	fmt.Println("type 'help' for a list of commands or 'exit' to close")
	fmt.Println("----------")
	for cli.active {

		cli.userInpt, _ = reader.ReadString('\n')
		// if err != nil {
		// 	fmt.Println("Some ungodly error occured")
		// }

		cli.userInpt = strings.TrimRight(cli.userInpt, "\r\n")
		fmt.Print("\033[H\033[2J")

		// fmt.Scanln(&cli.userInpt)
		cli.handleCommand(cli.userInpt, &game)

	}
}

//To Do List
//Make the viewBoard command easier to read, cause its fucking awful right now
//Make sure bishop cant move through pieces
//Work out Knight (and subsquently pawn movment)
//Work out more consistent error messaging with piece movement
//Make the board view, one of letters
// Find out when to clear cmd line for easier reading
