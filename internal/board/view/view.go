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

// Takes a single bitmap and adds the associated character to board array for specified y-level
func buildBoardLayer(layerSlice *[8][8]string, bm bitmap.Bitmap, vis string, yLevel int) {
	defer wg.Done()
	bm.Range(func(index uint32) {
		X, Y, Z := bitutil.UintToVec(index)
		if Y == yLevel {
			layerSlice[Z-1][X-1] = vis
		}
	})
}

// Internal function call to read the json storing board and output
func ViewLayer(yLevel int) {
	//Will allow for variable input later
	allPieces, _ := bitutil.LoadGame("data/output")
	var sliceOfBoard [8][8]string

	for meta, bm := range allPieces {
		wg.Add(1)
		go buildBoardLayer(&sliceOfBoard, bm, meta, yLevel)
	}
	wg.Wait()
	fmt.Printf("Layer : %d \n", yLevel)
	fmt.Println("╔══════════════════╗")
	for _, V := range sliceOfBoard {
		fmt.Print("║")
		for _, K := range V {
			if K == "" {
				fmt.Print(" -")
			} else {
				fmt.Print(" " + K)
			}
		}
		fmt.Print("  ║")
		fmt.Println()
	}
	fmt.Println("╚══════════════════╝")
	fmt.Println("  A B C D E F G H ")
}

func ViewAllLayers() {
	numLayers := int(BoardSize / LayerSize)
	for i := 0; i < numLayers; i++ {
		ViewLayer(i)
	}
}
