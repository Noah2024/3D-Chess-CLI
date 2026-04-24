package view

import (
	"3DC/util/logger"
	"3DC/util/must"
	"encoding/json"
	"os"
	"path/filepath"
)

const datapath = "../../data"

type board struct {
	Move    string `json:"move"`
	Rows    []string
	History string
}

// Internal function call to read the json storing board and output
func ViewBoard() {
	root := must.Must(os.Executable())
	datapath := filepath.Join(root, "../data/boards/default.json")
	logger.Warn(datapath)
	rawData := must.Must(os.ReadFile(datapath))
	//Consider adding a direct warning to the above if there is no data returned

	var data board
	json.Unmarshal(rawData, &data)

	print(data.Move)
	logger.Info(data.History)

}
