package view

import (
	"3DC/config"
	"3DC/internal/game/load"
	"3DC/util/bitutil"
	"3DC/util/metadata"
	"3DC/util/must"
	"fmt"
	"path/filepath"
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
func ViewLayer(yLevel int, displayMetaData bool) {
	//Will allow for variable input later
	allPieces, _ := load.LoadGame("data/output")

	if displayMetaData == true {
		meta := must.Must(metadata.LoadMetaData(filepath.Join("data/output", "meta")))
		metadata.DistplayMetaData(meta)
	}

	var sliceOfBoard [8][8]string
	for meta, bm := range allPieces {
		wg.Add(1)
		go buildBoardLayer(&sliceOfBoard, bm, meta, yLevel)

	}

	wg.Wait()
	fmt.Printf("Layer : %c \n", rune('A'+yLevel))
	zInc := 1
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
		fmt.Printf("  ║ %v", zInc)
		zInc += 1
		fmt.Println()
	}
	fmt.Println("╚══════════════════╝")
	fmt.Println("  A B C D E F G H ")
}

func ViewAllLayers() {
	numLayers := int(BoardSize / LayerSize)
	meta := must.Must(metadata.LoadMetaData(filepath.Join("data/output", "meta")))
	metadata.DistplayMetaData(meta)

	for i := 0; i < numLayers; i++ {
		ViewLayer(i, false)
	}
}
