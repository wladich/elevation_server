package main

import (
	"github.com/wladich/elevation_server/pkg/dem"
)

type orderedLatLon struct {
	i      int
	latlon dem.LatLon
}

func getElevations(storage dem.StorageReader, latlons []dem.LatLon) ([]float64, error) {
	elevations := make([]float64, len(latlons))
	tasks := make(map[dem.TileIndex][]orderedLatLon)
	for i, latlon := range latlons {
		tileIndex := dem.TileIndexFromLatLon(latlon)
		tasks[tileIndex] = append(tasks[tileIndex], orderedLatLon{i, latlon})
	}
	for tileIndex, task := range tasks {
		tile, err := storage.GetTile(tileIndex)
		if err != nil {
			return nil, err
		}
		if tile == nil {
			for _, orderedLatLon := range task {
				elevations[orderedLatLon.i] = dem.NoValue
			}
		} else {
			for _, orderedLatLon := range task {
				elevations[orderedLatLon.i] = tile.GetInterpolated(orderedLatLon.latlon)
			}
		}
	}
	return elevations, nil
}
