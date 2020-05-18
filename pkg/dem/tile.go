package dem

import (
	"github.com/wladich/elevation_server/pkg/constants"
	"math"
	"unsafe"
)

func (tile *Tile) GetInterpolated(latlon LatLon) float64 {
	x := (latlon.Lon*constants.HgtSplitParts - float64(tile.index.X)) * (constants.TileSize - 1)
	y := (latlon.Lat*constants.HgtSplitParts - float64(tile.index.Y)) * (constants.TileSize - 1)
	indX1 := int(math.Floor(x))
	indY1 := int(math.Floor(y))
	indX2 := indX1 + 1
	indY2 := indY1 + 1
	dx := x - float64(indX1)
	dy := y - float64(indY1)
	v1 := tile.data[indX1+indY1*constants.TileSize]
	v2 := tile.data[indX2+indY1*constants.TileSize]
	v3 := tile.data[indX1+indY2*constants.TileSize]
	v4 := tile.data[indX2+indY2*constants.TileSize]
	if v1 == NoValue || v2 == NoValue || v3 == NoValue || v4 == NoValue {
		return NoValue
	}
	return float64(v1)*(1-dx)*(1-dy) +
		float64(v2)*dx*(1-dy) +
		float64(v3)*(1-dx)*dy +
		float64(v4)*dx*dy
}

func tileFromRaw(rawTile TileRaw) Tile {
	return Tile{index: rawTile.Index, data: *(*TileData)(unsafe.Pointer(&rawTile.Data))}
}
