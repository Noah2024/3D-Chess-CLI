// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

// dialog package provides a simple way to prompt users for confirmation in the terminal.
// It includes a function to display a yes/no question and wait for the user's response, returning a boolean value based on the input.
// It is currently little used as I have yet to decide if I want to be constantly promoting users.
// But im leaning twords no.
package dialog

//Lowkey can chatgpt gen this for me
//But works perfectly

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Confirm prompts the user with a yes/no question and blocks execution
// until a valid response is given.
//
// Returns true for yes, false for no.
func Confirm(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", prompt)

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("\nError reading input, please try again.")
			continue
		}

		input = strings.TrimSpace(strings.ToLower(input))

		switch input {
		case "y", "yes":
			fmt.Println(input)
			return true
		case "n", "no":
			fmt.Println("Operation canceled by user")
			return false
		default:
			fmt.Println("Please enter 'y' or 'n'.")
		}
	}
}
