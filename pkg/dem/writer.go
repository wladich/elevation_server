package dem

import (
	"encoding/gob"
	"errors"
	"github.com/pierrec/lz4"
	"github.com/wladich/elevation_server/pkg/constants"
	"io"
	"os"
	"sync"
)

type StorageWriter struct {
	storageAbstract
	lock sync.Mutex
	fIdx *os.File
}

func NewWriter(path string) (*StorageWriter, error) {
	var storage StorageWriter
	idxPath := path + ".idx"
	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	storage.fData = f

	f, err = os.OpenFile(idxPath, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}
	storage.fIdx = f

	storage.index = &tileFileIndex{}
	return &storage, nil
}

func (storage *StorageWriter) Close() error {
	encoder := gob.NewEncoder(storage.fIdx)
	err1 := encoder.Encode(*storage.index)
	err2 := storage.fIdx.Close()
	err3 := storage.fData.Close()
	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}
	return err3
}

func compressTile(tileData TileRawData) ([]byte, error) {
	compressed := make([]byte, lz4.CompressBlockBound(len(tileData)))
	n, err := lz4.CompressBlockHC(tileData[:], compressed, 0)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, errors.New("compressed data has 0 size")
	}
	return compressed[:n], nil
}

func (storage *StorageWriter) PutTile(tile TileRaw) error {
	compressed, err := compressTile(tile.Data)
	if err != nil {
		return err
	}
	storage.lock.Lock()
	defer storage.lock.Unlock()
	pos, err := storage.fData.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}
	_, err = storage.fData.Write(compressed)
	if err != nil {
		return err
	}
	x := tile.Index.X + 180*constants.HgtSplitParts
	y := tile.Index.Y + 90*constants.HgtSplitParts
	if x < 0 || y < 0 || x > len(storage.index) || y > len(storage.index[x]) {
		return errors.New("tile index out of range")
	}
	storage.index[x][y] = tileFileIndexRecord{pos, int64(len(compressed))}
	return nil
}
