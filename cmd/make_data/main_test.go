package main

import "testing"

//func TestReadHgtFile(t *testing.T) {
//	tile, err := readHgtFile("path/to/hgts/N50E037.hgt.bz2")
//	if err != nil {
//		t.Fatal(err)
//	}
//	fmt.Print(tile[:10])
//}

func TestHgtIndexFromName(t *testing.T) {
	testData := []struct{
		name string
		lat, lon int
	}{
		{"S50W037.hgt.bz2", -50, -37},
		{"N50W037.hgt.bz2", 50, -37},
		{"S50E037.hgt.bz2", -50, 37},
		{"N50E037.hgt.bz2", 50, 37},

		{"S90E037.hgt", -90, 37},
		{"N89E037.hgt.gz", 89, 37},
		{"N50E179.hgt.bz2", 50, 179},
		{"N50W180.hgt.bz2", 50, -180},
	}
	for _, test := range testData {
		ind, err := hgtIndexFromName(test.name)
		if err != nil {
			t.Fatal(err)
		}
		if ind != (HgtIndex{test.lat, test.lon}) {
			t.Fatalf("wrong index for %v: %v", test, ind)
		}
	}
}
