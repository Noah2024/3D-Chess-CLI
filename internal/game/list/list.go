// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

// list provides the functionality to list all saved games in the data directory.
// utalizing the os.ReadDir function to read the contents of the data directory and print the names of all saved games to the console.
package list

import (
	"3DC/config"
	"3DC/util/logger"
	"3DC/util/must"
	"fmt"
	"io"
	"os"
)

func ListGames(w io.Writer) {
	gamesDir := must.Must(os.ReadDir(config.DataDir))
	if len(gamesDir) == 0 {
		logger.Warn("Could not find any saved games in data/games folder")
		return
	}

	fmt.Fprintln(w, "All Boards avilable to load")
	fmt.Fprintln(w, "----------")
	for i, v := range gamesDir {
		fmt.Fprintf(w, "%d: %s\n", i, v.Name())
	}
}
