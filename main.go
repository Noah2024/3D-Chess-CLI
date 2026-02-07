package main

import (
	"fmt"
	"strings"
	"unicode"
	"bufio"
	"os"
)

//HANDLING CLI CREATION AND SHIT
type CLI struct {
	userInpt string
	active bool
}

func (cli *CLI) handleClose(){
	cli.active = false;
}

func (cli *CLI) handleCommand(cmd string, game *GAME){
	fmt.Printf("IMPORTANT: %s \n", cmd)
	parsedCmd := strings.Fields(cmd) 

	switch parsedCmd[0] {
		case "exit":
			cli.handleClose()
		case "viewBoard":
			game.viewBoard()
		case "move":
			if len(parsedCmd) == 2 {
				game.move(parsedCmd[1])
			}else{
				fmt.Printf("move command takes exactly one argument '%s'\n", cmd)
			}
		default:
			fmt.Printf("Unknown Command '%s'\n", cmd)
	}
}
//CLI CREATION AND SHIT DONE
//Note to self: doing anything in which the number of pieces is different than 16, is not a great idea, FOR NOW
//STARTING GAME LOGIC
type intvec struct{
	intX int
	inty int 
	intz int
}

type GAME struct {
	turn int
	allPieces [32]PIECE
	white [16]PIECE
	black [16]PIECE
	boardSize intvec
	board [8][8]rune
}

func (game *GAME)initPieces()  {
	fmt.Println("Init Pieces")
	x,y,i := 0,0,0 //I is the piece index
	for _, character := range "RNBQKBNR/PPPPPPPP/8/8/8/8/pppppppp/rnbqkbnr"{
		if character == '/' {
			y++
			x = 0
			continue
		}else if unicode.IsDigit(character){
			x += int(character - '0')
		}else{
			game.allPieces[i] = PIECE{
			name: character,
			pos: intvec{x,y,0},
			}
			game.board[x][y] = character
			i++
		}
		x+=1
	}
}
//VIEW BOARD MUST BE CHANGED TO MOVE TO 3D
func (game *GAME) viewBoard(){
	for _, rank := range game.board{
		fmt.Println(rank)
	}	
}

func (game *GAME) move(cmd string){
	fmt.Println("MOVING TOO")
} 


type PIECE struct {
	name rune
	pos intvec
	id int
}
//GAME LOGIC DONE 
//START MISC BITS
func main(){
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("\033[H\033[2J")

	cli :=  CLI{
		userInpt: "",
		active: true,
	}

	game := GAME{
		turn: 0,
		boardSize: intvec{8,8,0},
		
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

		// fmt.Scanln(&cli.userInpt)
		cli.handleCommand(cli.userInpt, &game)

	}
}



