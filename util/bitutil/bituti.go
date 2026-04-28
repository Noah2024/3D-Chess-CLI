package bitutil

import (
	"3DC/config"
	"3DC/util/logger"
	"3DC/util/must"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/kelindar/bitmap"
)

//Below are the defualt values with which metadata is created

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
	EnPassant uint8

	//Time of last game save
	Time int64
}

// Default metadata declaration
var MetaDataVersion uint8 = 1
var Config uint8 = 1
var Turn uint8 = 1
var CastleRights uint8 = 0b11000011
var EnPessentRights uint8 = 0b00000000

// Creates directory at location and saves metadata as collection of bitmaps
func CreateSaveMetaData(location string) error {
	os.Mkdir(location, 0644)
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, MetaData{
		Version:   MetaDataVersion,
		Config:    Config,
		Turn:      Turn,
		Castle:    CastleRights,
		EnPassant: EnPessentRights,
		Time:      time.Now().Unix(),
	})

	err := os.WriteFile(filepath.Join(location, "meta.bin"), buf.Bytes(), 0644)
	must.Must("", err)
	return nil
}

// Loads the bin file from the given location and reads it into predefined metdata struct
func LoadMetaData(location string) (MetaData, error) {
	var data MetaData

	filePath := filepath.Join(location, "meta.bin")

	b := must.Must(os.ReadFile(filePath))

	buf := bytes.NewReader(b)

	err := binary.Read(buf, binary.LittleEndian, &data)
	must.Must("", err)

	return data, nil
}

func DistplayMetaData(meta MetaData) {
	fmt.Println("Meta Data")
	fmt.Println("----------")
	fmt.Printf("Version: %d\n", meta.Version)
	fmt.Printf("Config: %d \n", meta.Config)
	fmt.Printf("Turn: %d \n", meta.Turn)
	fmt.Printf("Castle: %d \n", meta.Castle)
	fmt.Printf("EnPessent %d\n", meta.EnPassant)
	fmt.Printf("Saved at %s\n", time.Unix(meta.Time, 0).UTC().Format(time.RFC3339))
	fmt.Println("----------")
}

const (
	BoardSize = config.BoardSize
	LayerSize = config.LayerSize
	LineSize  = config.LineSize
	SpaceSize = config.SpaceSize
)

func VecToUint(x, y, z int) uint32 {
	return uint32(x + (y-1)*int(LayerSize) + (z-1)*int(LineSize))
}

// Decodes uint32 position into integer x,y,z position
func UintToVec(space uint32) (int, int, int) {
	if space < 1 || space > 512 {
		logger.Error(fmt.Sprintf("uint32 %d out of range for board size %d ", space, BoardSize))
		panic("See above error")
	}

	// index = x + y*8 + z*64 essentially decoding this
	//Step by step removing the largest term at a time
	space-- // convert to 0-based
	y := space / LayerSize
	space %= LayerSize
	z := space / LineSize
	x := space % LineSize

	return int(x + 1), int(y + 1), int(z + 1)
}

func SaveGame(bmps map[rune]bitmap.Bitmap, location string) {
	os.Mkdir(location, 0744) //Owner can rxw but everyone else can only r
	CreateSaveMetaData(filepath.Join(location, "meta"))
	for key, bm := range bmps {
		fileLoc := location + "/" + string(key)
		file := must.Must(os.Create(fileLoc))
		bm.WriteTo(file)
	}
}

// Loads dictionary mapping display char to the bitmap corresponding with that piece
// Also returns
func LoadGame(location string) (data map[string]bitmap.Bitmap, err error) {

	result := make(map[string]bitmap.Bitmap)

	entries := must.Must(os.ReadDir(location))

	for _, entry := range entries {
		if entry.IsDir() {
			//meta.bin is stored in /meta dir
			//This is so we can easily skip it when loading the piece bitmaps
			continue
		}
		file := must.Must(os.Open(location + "/" + entry.Name()))

		bm := must.Must(bitmap.ReadFrom(file))
		result[entry.Name()] = bm
	}
	return result, nil

}
