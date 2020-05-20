package dem

import (
	"math"
	"unsafe"
)

func (tile *Tile) GetInterpolated(latlon LatLon) float64 {
	x := (latlon.Lon*HgtSplitParts - float64(tile.index.X)) * (TileSize - 1)
	y := (latlon.Lat*HgtSplitParts - float64(tile.index.Y)) * (TileSize - 1)
	indX1 := int(math.Floor(x))
	indY1 := int(math.Floor(y))
	indX2 := indX1 + 1
	indY2 := indY1 + 1
	dx := x - float64(indX1)
	dy := y - float64(indY1)
	v1 := tile.data[indX1+indY1*TileSize]
	v2 := tile.data[indX2+indY1*TileSize]
	v3 := tile.data[indX1+indY2*TileSize]
	v4 := tile.data[indX2+indY2*TileSize]
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
