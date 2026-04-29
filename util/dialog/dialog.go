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
