package dem

import (
	"github.com/wladich/elevation_server/pkg/constants"
	"math"
	"os"
)

const NoValue = -32768

type LatLon struct {
	Lat, Lon float64
}

type TileIndex struct {
	X, Y int
}

type TileRawData [constants.TileBytes]byte
type TileData [constants.TilePointsN]int16

type TileRaw struct {
	Data  TileRawData
	Index TileIndex
}

type tileFileIndexRecord struct {
	Offset int64
	Size   int64
}

type tileFileIndex [360 * constants.HgtSplitParts][180 * constants.HgtSplitParts]tileFileIndexRecord

type storageAbstract struct {
	fData *os.File
	index *tileFileIndex
}

type Tile struct {
	data  TileData
	index TileIndex
}

func TileIndexFromLatLon(latlon LatLon) TileIndex {
	x := int(math.Floor(latlon.Lon * constants.HgtSplitParts))
	y := int(math.Floor(latlon.Lat * constants.HgtSplitParts))
	return TileIndex{x, y}
}
