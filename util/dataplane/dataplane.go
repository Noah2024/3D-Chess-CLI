// Planes along with pieces can move are precomputed here for use in move.go
package dataplane

//Stored here for debug
//TO DO LIST
//Figure out how dependencies will work for this system
//Refactor config to take advantage of dataplanes
// fmt.Printf("%064b\n", uint64(18446744073709551614))
import (
	"3DC/config"
	"3DC/util/bitutil"
	"fmt"

	"github.com/kelindar/bitmap"
)

var YPlane = [8]bitmap.Bitmap{
	{18446744073709551615, 0, 0, 0, 0, 0, 0, 0},
	{0, 18446744073709551615, 0, 0, 0, 0, 0, 0},
	{0, 0, 18446744073709551615, 0, 0, 0, 0, 0},
	{0, 0, 0, 18446744073709551615, 0, 0, 0, 0},
	{0, 0, 0, 0, 18446744073709551615, 0, 0, 0},
	{0, 0, 0, 0, 0, 18446744073709551615, 0, 0},
	{0, 0, 0, 0, 0, 0, 18446744073709551615, 0},
	{0, 0, 0, 0, 0, 0, 0, 18446744073709551615},
}

var ZPlane = [8]bitmap.Bitmap{

	{255, 255, 255, 255, 255, 255, 255, 255},
	{65280, 65280, 65280, 65280, 65280, 65280, 65280, 65280},
	{16711680, 16711680, 16711680, 16711680, 16711680, 16711680, 16711680, 16711680},
	{4278190080, 4278190080, 4278190080, 4278190080, 4278190080, 4278190080, 4278190080, 4278190080},
	{1095216660480, 1095216660480, 1095216660480, 1095216660480, 1095216660480, 1095216660480, 1095216660480, 1095216660480},
	{280375465082880, 280375465082880, 280375465082880, 280375465082880, 280375465082880, 280375465082880, 280375465082880, 280375465082880},
	{71776119061217280, 71776119061217280, 71776119061217280, 71776119061217280, 71776119061217280, 71776119061217280, 71776119061217280, 71776119061217280},
	{18374686479671623680, 18374686479671623680, 18374686479671623680, 18374686479671623680, 18374686479671623680, 18374686479671623680, 18374686479671623680, 18374686479671623680},
}

var XPlane = [8]bitmap.Bitmap{
	{72340172838076673, 72340172838076673, 72340172838076673, 72340172838076673, 72340172838076673, 72340172838076673, 72340172838076673, 72340172838076673},
	{144680345676153346, 144680345676153346, 144680345676153346, 144680345676153346, 144680345676153346, 144680345676153346, 144680345676153346, 144680345676153346},
	{289360691352306692, 289360691352306692, 289360691352306692, 289360691352306692, 289360691352306692, 289360691352306692, 289360691352306692, 289360691352306692},
	{578721382704613384, 578721382704613384, 578721382704613384, 578721382704613384, 578721382704613384, 578721382704613384, 578721382704613384, 578721382704613384},
	{1157442765409226768, 1157442765409226768, 1157442765409226768, 1157442765409226768, 1157442765409226768, 1157442765409226768, 1157442765409226768, 1157442765409226768},
	{2314885530818453536, 2314885530818453536, 2314885530818453536, 2314885530818453536, 2314885530818453536, 2314885530818453536, 2314885530818453536, 2314885530818453536},
	{4629771061636907072, 4629771061636907072, 4629771061636907072, 4629771061636907072, 4629771061636907072, 4629771061636907072, 4629771061636907072, 4629771061636907072},
	{9259542123273814144, 9259542123273814144, 9259542123273814144, 9259542123273814144, 9259542123273814144, 9259542123273814144, 9259542123273814144, 9259542123273814144},
}

// A quick function to check if various uint32 indexes are continaed in a given bitmap
// Used for testing the validity of generated dataplanes in development
func testIndex(bm bitmap.Bitmap, nums []uint32) {
	for _, num := range nums {
		fmt.Println(bm.Contains(num))
	}
}

// Don't even say a god damn word on the atrocious time complexity of this function.
// Its only meant to be run once EVER, its saying in the source code for reference
func GenerateAllPlanes() {

	var xzPlane [int(config.LineSize)]bitmap.Bitmap //Set Z
	for i := 0; i < int(config.LineSize); i++ {
		plane := GeneratePlane(func(x, y, z int) bool {
			return y == i //Zero indexed
		})
		xzPlane[i] = plane
	}

	var xyPlane [int(config.LineSize)]bitmap.Bitmap // Set Y
	for i := 0; i < int(config.LineSize); i++ {
		plane := GeneratePlane(func(x, y, z int) bool {
			return z == i //Zero indexed
		})
		xyPlane[i] = plane
	}

	var zyPlane [int(config.LineSize)]bitmap.Bitmap // Set X
	for i := 0; i < int(config.LineSize); i++ {
		plane := GeneratePlane(func(x, y, z int) bool {
			return x == i //Zero indexed
		})
		zyPlane[i] = plane
	}

	tmp := xyPlane[3].Clone(nil) // Z
	fmt.Printf("tmp: %064b\n", tmp)

	tmp2 := xzPlane[3].Clone(nil) // Y
	fmt.Printf("tmp2: %064b\n", tmp2)

	tmp3 := zyPlane[3].Clone(nil) // X
	fmt.Printf("tmp3: %064b\n", tmp3)

	fmt.Println("Please god")
	// tmp2, _ := tmp.Max()
	// fmt.Println(tmp2 + 1)
	// testIndex(tmp, []uint32{0, 1, 2, 3, 4, 5, 6, 7, 8, 63, 64, 65})

	//Temporary Visual Output
	// fmt.Println(tmp)
	// fmt.Println(xzPlane) //Y
	// fmt.Println("{{{}}}")
	// fmt.Println(xyPlane) //Z
	// fmt.Println("{{{}}}")
	// fmt.Println(zyPlane) //X
}

func GeneratePlane(fn func(x, y, z int) bool) bitmap.Bitmap {
	var bm bitmap.Bitmap
	bm.Grow(511) //again, bitmap is zero indexed
	// bm.Ones()
	//Must use zero indexed because were talking about the indexing of the bitmap itself
	//Rather than simply indexing the board, which for other reasons starts at one
	for y := 0; y < int(config.LineSize); y++ {
		for z := 0; z < int(config.LineSize); z++ {
			for x := 0; x < int(config.LineSize); x++ {
				if fn(x, y, z) {
					// fmt.Printf("Letting on %s, %s, %s \n", x, y, z)
					// idx := uint32(((x + 1) + (y)*int(config.LayerSize) + (z)*int(config.LineSize)) - 1)
					idx := bitutil.VecToUint(x+1, y+1, z+1) - 1

					bm.Set(idx)
				}
			}
		}
	}
	// fmt.Println(bm)
	return bm
}

func TestDataPlane(fx int, fy int, fz int) {
	// XZ plane at y=4
	GenerateAllPlanes()
}
