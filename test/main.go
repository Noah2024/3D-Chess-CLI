package main

import "fmt"

func main() {
	// \x1b[31m is Red, \x1b[0m is Reset
	fmt.Printf("This is \x1b[31mRed\x1b[0m and this is normal.\n")
}
