package main

import (
	"compress/bzip2"
	"compress/gzip"
	"errors"
	"flag"
	"github.com/cheggaaa/pb/v3"
	"github.com/wladich/elevation_server/pkg/constants"
	"github.com/wladich/elevation_server/pkg/dem"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

type HgtRawData [constants.HgtSize * constants.HgtSize * 2]byte

type HgtIndex struct {
	lat, lon int
}

type TileRawSet [constants.HgtSplitParts * constants.HgtSplitParts]dem.TileRaw

func readHgtFile(path string) (*HgtRawData, error) {
	ext := strings.ToLower(filepath.Ext(path))
	f, err := os.Open(path)
	var reader io.Reader
	if err != nil {
		return nil, err
	}
	defer f.Close()
	switch ext {
	case ".hgt":
		reader = f
	case ".bz2":
		reader = bzip2.NewReader(f)
	case ".gz":
		reader, err = gzip.NewReader(f)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("unknown hgt file extension")
	}
	rawData, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	if len(rawData) != constants.HgtSize*constants.HgtSize*2 {
		return nil, errors.New("invalid HGT size")
	}
	var result HgtRawData
	for i := 0; i < constants.HgtSize*constants.HgtSize; i++ {
		result[i*2], result[i*2+1] = rawData[i*2+1], rawData[i*2]
	}
	return &result, nil
}

func hgtIndexFromName(name string) (HgtIndex, error) {
	re := regexp.MustCompile("([NS])(\\d{2})([EW])(\\d{3})\\.HGT")
	m := re.FindStringSubmatch(strings.ToUpper(name))
	if m != nil {
		lat, err1 := strconv.Atoi(m[2])
		lon, err2 := strconv.Atoi(m[4])
		if err1 == nil && err2 == nil {
			if m[1] == "S" {
				lat = -lat
			}
			if m[3] == "W" {
				lon = -lon
			}
			return HgtIndex{lat: lat, lon: lon}, nil
		}
	}
	return HgtIndex{}, errors.New("invalid HGT file name")
}

func splitDem(index HgtIndex, data *HgtRawData) (tiles TileRawSet) {
	for tileDx := 0; tileDx < constants.HgtSplitParts; tileDx++ {
		for tileDy := 0; tileDy < constants.HgtSplitParts; tileDy++ {
			tileX := index.lon*constants.HgtSplitParts + tileDx
			tileY := index.lat*constants.HgtSplitParts + tileDy
			tileI := tileDy*constants.HgtSplitParts + tileDx
			tile := &tiles[tileI]
			tile.Index.X = tileX
			tile.Index.Y = tileY
			for row := 0; row < constants.TileSize; row++ {
				hgtRow := constants.HgtSize - 1 - tileDy*(constants.TileSize-1) - row
				hgtIndex := hgtRow*constants.HgtSize + tileDx*(constants.TileSize-1)
				n := copy(tile.Data[row*constants.TileSize*2:], data[hgtIndex*2:(hgtIndex+constants.TileSize)*2])
				if n != constants.TileSize*2 {
					panic("invalid number of bytes copied")
				}
			}
		}
	}
	return
}

func processHgt(filename string, hgtDir string, storage *dem.StorageWriter) (error) {
	hgt, err := readHgtFile(path.Join(hgtDir, filename))
	if err != nil {
		return err
	}
	index, err := hgtIndexFromName(filename)
	if err != nil {
		return err
	}
	for _, tile := range splitDem(index, hgt) {
		if err = storage.PutTile(tile); err != nil {
			return err
		}
	}
	return nil
}

func makeTiles(hgtDir, demStorageFile string, concurency int) {
	storage, err := dem.NewWriter(demStorageFile)
	if err != nil {
		panic(err)
	}

	files, err := ioutil.ReadDir(hgtDir)
	if err != nil {
		log.Fatal(err)
	}

	bar := pb.Full.Start(len(files))
	defer bar.Finish()

	jobs := make(chan string, concurency)

	results := make(chan error, concurency)

	go func() {
		for _, filename := range files {
			jobs <- filename.Name()
		}
		close(jobs)
	}()

	for i := 0; i < concurency; i++ {
		go func () {
			for filename := range jobs {
				results <- processHgt(filename, hgtDir, storage)
			}
		}()
	}

	for i := 0; i < len(files); i++ {
		if err := <- results; err != nil {
			panic(err)
		}
		bar.Increment()
	}

	err = storage.Close()
	if err != nil {
		panic(err)
	}
}

func main() {
	hgtDir := flag.String("hgt", "", "Directory with hgt files in .hgt, .bz2 or .gz format")
	demStorageFile := flag.String("out", "", "Output file name")
	flag.Parse()
	if *hgtDir == "" || *demStorageFile == "" {
		flag.Usage()
		os.Exit(1)
	}
	numCPUs := runtime.NumCPU()
	makeTiles(*hgtDir, *demStorageFile, numCPUs + 1)
}
