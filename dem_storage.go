package main

import (
	"bytes"
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pierrec/lz4"
	"io"
)

type DemStorage struct {
	conn     *sql.DB
	tileStmt *sql.Stmt
}

var TileNotFound = errors.New("not found")

func uncompressLZ4(buf []byte) ([]byte, error) {
	lz4Reader := lz4.NewReader(bytes.NewReader(buf))
	var uncompressed bytes.Buffer
	_, err := io.Copy(&uncompressed, lz4Reader)
	if err != nil {
		return nil, err
	}
	return uncompressed.Bytes(), nil
}

func (store *DemStorage) getDemTile(index TileIndex) (*Tile, error) {
	rows, err := store.tileStmt.Query(index.x, index.y)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		var buf []byte
		rows.Scan(&buf)
		uncompressed, err := uncompressLZ4(buf)
		if err != nil {
			return nil, err
		}
		return &Tile{data: bytesToInt16(uncompressed), index: index}, nil
	} else {
		err = rows.Err()
		if err == nil {
			err = TileNotFound
		}
		return nil, err
	}
}

func openDemStorage(path string) (store DemStorage, err error) {
	store = DemStorage{}
	store.conn, err = sql.Open("sqlite3", path)
	if err != nil {
		return
	}
	store.tileStmt, err = store.conn.Prepare("SELECT tile_data FROM dem_tiles WHERE x=? AND y=?")
	return
}
