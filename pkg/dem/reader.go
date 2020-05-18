package dem

import (
	"errors"
	"fmt"
	"github.com/edsrzf/mmap-go"
	"github.com/hungys/go-lz4"
	"github.com/wladich/elevation_server/pkg/constants"
	"io"
	"os"
	"reflect"
	"unsafe"
)

type StorageReader storageAbstract

func NewReader(path string) (*StorageReader, error) {
	var storage StorageReader
	idxPath := path + ".idx"
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	storage.fData = f

	f, err = os.Open(idxPath)
	if err != nil {
		return nil, err
	}
	storage.fIdx = f

	storage.indexMmap, err = mmap.Map(storage.fIdx, mmap.RDONLY, 0)
	if err != nil {
		return nil, err
	}
	mmapData := (*reflect.SliceHeader)(unsafe.Pointer(&storage.indexMmap)).Data
	storage.index = (*tileFileIndex)(unsafe.Pointer(mmapData))
	return &storage, nil
}

func decompressTile(compressed []byte) (*TileRawData, error) {
	var tileData TileRawData
	n, err := lz4.DecompressSafe(compressed, tileData[:])
	if n != constants.TileBytes {
		return nil, errors.New(fmt.Sprintf("Unexpected tile size: %v", n))
	}
	if err != nil {
		return nil, err
	}
	return &tileData, nil
}

func (storage StorageReader) GetTile(index TileIndex) (*Tile, error) {
	x := index.X + 180*constants.HgtSplitParts
	y := index.Y + 90*constants.HgtSplitParts
	if x < 0 || y < 0 || x > len(storage.index) || y > len(storage.index[x]) {
		return nil, nil
	}
	tileFileIndex := storage.index[x][y]
	if tileFileIndex.size == 0 {
		return nil, nil
	}
	if _, err := storage.fData.Seek(tileFileIndex.offset, io.SeekStart); err != nil {
		return nil, err
	}
	compressed := make([]byte, tileFileIndex.size)
	n, err := storage.fData.Read(compressed)
	if err != nil {
		return nil, err
	}
	if int64(n) != tileFileIndex.size {
		return nil, errors.New("tile data incomplete")
	}
	tileData, err := decompressTile(compressed)
	if err != nil {
		return nil, err
	}
	tile := tileFromRaw(TileRaw{*tileData, index})
	return &tile, nil
}

func (storage StorageReader) Close() error {
	err1 := storage.indexMmap.Unmap()
	err2 := storage.fIdx.Close()
	err3 := storage.fData.Close()
	storage.index = nil
	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}
	return err3
}
