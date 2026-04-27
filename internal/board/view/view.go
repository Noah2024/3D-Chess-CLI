package view

import (
	"3DC/config"
	"3DC/util/bitutil"
	"fmt"
	"sync"

	"github.com/kelindar/bitmap"
)

const (
	BoardSize = config.BoardSize
	LayerSize = config.LayerSize
	LineSize  = config.LineSize
	SpaceSize = config.SpaceSize
)

const datapath = "../../data"

var wg sync.WaitGroup

// May need to change this depdent on size
var yLevel = 3
var sliceOfBoard [8][8]string

// Takes a single bitmap and adds the necssary visual rune to the slice of board array for later display
func buildBoardSlice(bm bitmap.Bitmap, vis string) {
	fmt.Print("Go Routine?")
	defer wg.Done()
	bm.Range(func(index uint32) {
		X, Y, Z := bitutil.UintToVec(index)
		if Y == yLevel {
			sliceOfBoard[Z-1][X-1] = vis
		}
	})
}

// Internal function call to read the json storing board and output
func ViewBoard() {

	//Will allow for variable input later
	allPieces, _ := bitutil.LoadGame("data/output")
	fmt.Print(len(allPieces))

	for meta, bm := range allPieces {
		wg.Add(1)
		go buildBoardSlice(bm, meta)
	}
	wg.Wait()

	fmt.Println("SLICE OF BOARD BELOW")
	fmt.Print(sliceOfBoard)
}
