// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

package view

import (
	"3DC/config"
	"3DC/internal/game/load"
	"3DC/util/bitutil"
	"3DC/util/metadata"
	"3DC/util/must"
	"fmt"
	"io"
	"sync"

	"github.com/kelindar/bitmap"
)

const (
	BoardSize = config.BoardSize
	LayerSize = config.LayerSize
	LineSize  = config.LineSize
	SpaceSize = config.SpaceSize
)

var wg sync.WaitGroup

// May need to change this depdent on size

// Takes a single bitmap and adds the associated character to board array for specified y-level
func buildBoardLayer(layerSlice *[8][8]string, bm bitmap.Bitmap, vis string, yLevel int) {
	defer wg.Done()
	bm.Range(func(index uint32) {
		X, Y, Z := bitutil.UintToVec(index)
		if Y == yLevel {
			layerSlice[Z][X] = vis
		}
	})
}

func BuildLayer(allPieces map[string]bitmap.Bitmap, yLevel int) [8][8]string {
	//Will allow for variable input later

	var board [8][8]string
	for meta, bm := range allPieces {
		wg.Add(1)
		go buildBoardLayer(&board, bm, meta, yLevel)
	}

	wg.Wait()
	return board
}

// Internal function call to bitmap storing the
func PrintLayer(yLevel int, displayMetaData bool, w io.Writer) {
	allPieces, meta, _ := load.LoadGame(config.CurrentGame)

	sliceOfBoard := BuildLayer(allPieces, yLevel)

	if displayMetaData == true {
		metadata.DistplayMetaData(meta)
	}

	fmt.Fprintf(w, "Layer : %c \n", rune('A'+yLevel))
	zInc := 1
	fmt.Fprintln(w, "╔══════════════════╗")
	for _, V := range sliceOfBoard {
		fmt.Fprint(w, "║")
		for _, K := range V {
			if K == "" {
				fmt.Fprint(w, " -")
			} else {
				fmt.Fprint(w, " "+K)
			}
		}
		fmt.Fprintf(w, "  ║ %v", zInc)
		zInc += 1
		fmt.Fprintln(w)
	}
	fmt.Fprintln(w, "╚══════════════════╝")
	fmt.Fprintln(w, "  A B C D E F G H ")
}

func ViewAllLayers(w io.Writer) {
	numLayers := int(BoardSize / LayerSize)
	meta := must.Must(metadata.LoadMetaData(config.CurrentGame))
	metadata.DistplayMetaData(meta)

	for i := 0; i < numLayers; i++ {
		PrintLayer(i, false, w)
	}
}
