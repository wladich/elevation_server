package main

import (
	"database/sql"
	"errors"
	"fmt"
	lz4 "github.com/hungys/go-lz4"
	_ "github.com/mattn/go-sqlite3"
)

type DemStorage struct {
	conn     *sql.DB
	tileStmt *sql.Stmt
}

var TileNotFound = errors.New("not found")

func uncompressLZ4(buf []byte) ([]byte, error) {
	uncompressed := make([]byte, 301*301*2)
	n, err := lz4.DecompressSafe(buf, uncompressed)
	if n != 301*301*2 {
		return nil, errors.New(fmt.Sprintf("Unexpected tile size: %v", n))
	}
	if err != nil {
		return nil, err
	}
	return uncompressed, nil
}

func (store *DemStorage) getDemTile(index TileIndex) (*Tile, error) {
	rows, err := store.tileStmt.Query(index.x, index.y)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		var buf []byte
		err = rows.Scan(&buf)
		if err != nil {
			return nil, err
		}
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
