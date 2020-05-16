package main

type LatLon struct {
	lat, lon float64
}

func getElevations(storage DemStorage, latlons []LatLon) ([]float64, error) {
	type OrderedLatLon struct {
		i      int
		latlon LatLon
	}
	elevations := make([]float64, len(latlons))
	tasks := make(map[TileIndex][]OrderedLatLon)
	for i, latlon := range latlons {
		tileIndex := tileIndexFromLatLon(latlon)
		tasks[tileIndex] = append(tasks[tileIndex], OrderedLatLon{i, latlon})
	}
	for tileIndex, task := range tasks {
		tile, err := storage.getDemTile(tileIndex)
		switch err {
		case TileNotFound:
			for _, orderedLatLon := range task {
				elevations[orderedLatLon.i] = NoValue
			}
		case nil:
			for _, orderedLatLon := range task {
				elevations[orderedLatLon.i] = tile.getInterpolated(orderedLatLon.latlon)
			}
		default:
			return nil, err
		}
	}
	return elevations, nil
}
