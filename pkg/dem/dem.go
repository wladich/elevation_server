package dem

import (
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

type TileRawData [TileBytes]byte
type TileData [TilePointsN]int16

type TileRaw struct {
	Data  TileRawData
	Index TileIndex
}

type tileFileIndexRecord struct {
	Offset int64
	Size   int64
}

type tileFileIndex [360 * HgtSplitParts][180 * HgtSplitParts]tileFileIndexRecord

type storageAbstract struct {
	fData *os.File
	index *tileFileIndex
}

type Tile struct {
	data  TileData
	index TileIndex
}

func TileIndexFromLatLon(latlon LatLon) TileIndex {
	x := int(math.Floor(latlon.Lon * HgtSplitParts))
	y := int(math.Floor(latlon.Lat * HgtSplitParts))
	return TileIndex{x, y}
}
