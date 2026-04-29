package list

import (
	"3DC/config"
	"3DC/util/logger"
	"3DC/util/must"
	"fmt"
	"os"
)

func ListGames() {
	gamesDir := must.Must(os.ReadDir(config.DataDir))
	if len(gamesDir) == 0 {
		logger.Warn("Could not find any saved games in data/games folder")
		return
	}
	fmt.Print("All Boards avilable to load \n")
	fmt.Println("----------")
	for I, V := range gamesDir {
		fmt.Printf("%d: %v \n", I, V.Name())
	}
}
