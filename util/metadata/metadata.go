// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3

// metadata package provides functionality for creating, saving, and loading metadata related to the game state. It defines the MetaData struct, which holds information about the game version, configuration, turn, castling rights, en passant rights, and the time of the last save. The package includes functions to create and save metadata to a specified location, load metadata from a file, and display the metadata in a human-readable format.
package metadata

//Will need ot update metadata to use io.writer intead of fmt

import (
	"3DC/util/logger"
	"3DC/util/must"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// MetaData struct holds information about the game state, including version, configuration, turn, castling rights, en passant rights, and the time of the last save.
type MetaData struct {
	// Version contains data about the version of data stored
	Version uint8

	// Config about in what way the data is meant to be read (in standard 8x8x8 or somthing else)
	Config uint8

	Turn uint8

	// Holds single byte represting who has castling rights
	// 1(white queenside)1(white kingside) 0000(padding) (1 blackkingside)(black queenside)
	Castle uint8

	// Becuase enpessent rights exist for only one turn this represents the rights of whoevers turn it is
	// Determined at the end of the previous players turn
	// Single byte representing who has enpessent rights for this next turn
	WhiteEnPessent uint8
	BlackEnPessent uint8

	//Time of last game save
	LastSaved int64

	Created int64
}

// Default metadata declaration
var MetaDataVersion uint8 = 1
var Config uint8 = 1
var Turn uint8 = 0
var CastleRights uint8 = 0b11000011
var WhiteEnPessent uint8 = 9 //0b0
var BlackEnPessent uint8 = 9 //0b0

// Creates directory at location and saves metadata as collection of bitmaps //CreateSaveMetaData
func CreateDefaultMetaData() MetaData {
	return MetaData{
		Version:        MetaDataVersion,
		Config:         Config,
		Turn:           Turn,
		Castle:         CastleRights,
		WhiteEnPessent: WhiteEnPessent,
		BlackEnPessent: BlackEnPessent,
		LastSaved:      time.Now().Unix(),
		Created:        time.Now().Unix(),
	}
}

// Loads the bin file from the given location and reads it into predefined metdata struct
func LoadMetaData(location string) (MetaData, error) {
	var data MetaData

	filePath := filepath.Join(location, "meta", "meta.bin")
	b := must.Must(os.ReadFile(filePath))
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.LittleEndian, &data)
	must.Must("", err)

	return data, nil
}

func SaveMetaData(data MetaData, location string) {

	filePath := filepath.Join(location, "meta", "meta.bin")
	metaDir := filepath.Join(location, "meta")

	_, err := os.Stat(metaDir)
	if err != nil {
		logger.Debug("No metadata currently exists for this game, creating it now")
		os.Mkdir(metaDir, 0o755)
	}

	data.LastSaved = time.Now().Unix()

	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, data)
	err = os.WriteFile(filePath, buf.Bytes(), 0o755)
	must.Must("", err)
}

// Displays the metadata in a human-readable format, including version, configuration, turn, castling rights, en passant rights, and the time of the last save.
// Called mostly by the view function
func DistplayMetaData(meta MetaData) {
	fmt.Println("Meta Data")
	fmt.Println("----------")
	fmt.Printf("Version: %d\n", meta.Version)
	fmt.Printf("Config: %d \n", meta.Config)
	fmt.Printf("Turn: %d \n", meta.Turn)
	fmt.Printf("Castle: %d \n", meta.Castle)
	fmt.Printf("WhiteEnPessent %d\n", meta.WhiteEnPessent)
	fmt.Printf("BlackEnPessent %d\n", meta.BlackEnPessent)
	fmt.Printf("Game Created %s\n", time.Unix(meta.Created, 0).UTC().Format(time.RFC3339))
	fmt.Printf("Last Saved %s\n", time.Unix(meta.LastSaved, 0).UTC().Format(time.RFC3339))
	fmt.Println("----------")
}
