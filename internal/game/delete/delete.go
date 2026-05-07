package delete

import (
	"3DC/config"
	"3DC/util/dialog"
	"3DC/util/must"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

// Helper function which confirms users choice and does the deletion
func deleteDir(entry os.DirEntry) {
	if !dialog.Confirm(fmt.Sprintf("Are you sure you want to delete game '%v'?", entry.Name())) {
		return
	}
	err := os.RemoveAll(filepath.Join(config.DataDir, entry.Name()))
	must.Must("", err)
}

// Main delete function: loops the game directory to see if users choice matches any entry
func DeleteGame(name string) error {

	dir := must.Must(os.ReadDir(config.DataDir))

	index := -1
	if i, err := strconv.Atoi(name); err == nil {
		index = i
	}

	for i, entry := range dir {

		if index != -1 && i == index {
			deleteDir(entry)
			return nil
		}

		if entry.Name() == name {
			deleteDir(entry)
			return nil
		}

	}

	fmt.Println("Could not find specified game")
	return nil
}
