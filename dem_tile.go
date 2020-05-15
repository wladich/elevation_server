package main

import (
	"encoding/binary"
	"math"
)

const NoValue = -32768
const HgtSize = 1201
const HgtSplitParts = 4
const TileSize = (HgtSize-1)/HgtSplitParts + 1

type TileIndex struct {
	x, y int
}

type Tile struct {
	data  []int16
	index TileIndex
}

func bytesToInt16(b []byte) []int16 {
	h := make([]int16, len(b)/2)
	for i := range h {
		h[i] = int16(binary.LittleEndian.Uint16(b[i*2 : (i+1)*2]))
	}
	return h
}

func (tile *Tile) getInterpolated(latlon LatLon) float64 {
	x := (latlon.lon*HgtSplitParts - float64(tile.index.x)) * (TileSize - 1)
	y := (latlon.lat*HgtSplitParts - float64(tile.index.y)) * (TileSize - 1)
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

func tileIndexFromLatLon(latlon LatLon) TileIndex {
	x := int(math.Floor(latlon.lon * HgtSplitParts))
	y := int(math.Floor(latlon.lat * HgtSplitParts))
	return TileIndex{x, y}
}
