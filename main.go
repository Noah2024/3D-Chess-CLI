package main

import (
	"fmt"
	"strings"
	// "os"
)

//HANDLING CLI CREATION AND SHIT
type CLI struct {
	userInpt string
	active bool
}

func (cli *CLI) handleClose(){
	cli.active = false;
}

func (cli *CLI) handleCommand(cmd string){

	parsedCmd := strings.Fields(cmd) 

	switch parsedCmd[0] {
		case "exit":
			cli.handleClose()
		default:
			fmt.Printf("Unknown Command '%s'\n", cmd)
	}
}
//CLI CREATION AND SHIT DONE
//STARTING GAME LOGIC
type GAME struct {
	turn int
	white [16]PIECE
	black [16]PIECE
}

type PIECE struct {
	name string
	id int
}
//GAME LOGIC DONE 
//START MISC BITS
func main(){

	cli :=  CLI{
		userInpt: "",
		active: true,
	}

	// fmt.Println("test")
	for cli.active {
		fmt.Scan(&cli.userInpt)
		cli.handleCommand(cli.userInpt)
	}
}



