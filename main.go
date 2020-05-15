package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	_ "net/http/pprof"
)

const MaxInputSize = 250000
const MaxInputPoints = 10000

var demStorage *DemStorage

func handleRequest(resp http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(resp, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if req.URL.Path != "/" {
		http.NotFound(resp, req)
		return
	}

	contentLengthHeaders := req.Header["Content-Length"]
	if len(contentLengthHeaders) > 0 {
		contentLength, err := strconv.ParseInt(contentLengthHeaders[0], 10, 32)
		if err == nil && contentLength > MaxInputSize {
			http.Error(resp, "Request too big", http.StatusRequestEntityTooLarge)
			return
		}
	}

	var latlons []LatLon
	inputLinesReader := bufio.NewReader(http.MaxBytesReader(resp, req.Body, MaxInputSize))
	for {
		line, readErr := inputLinesReader.ReadString('\n')
		if readErr != nil && readErr != io.EOF {
			http.Error(resp, "Request too big", http.StatusRequestEntityTooLarge)
			return
		}
		if len(line) > 0 {
			fields := strings.Fields(line)
			if len(fields) != 2 {
				http.Error(resp, "Invalid request", http.StatusBadRequest)
				return
			}
			lat, err := strconv.ParseFloat(fields[0], 64)
			if err != nil {
				http.Error(resp, "Invalid request", http.StatusBadRequest)
				return
			}
			lon, err := strconv.ParseFloat(fields[1], 64)
			if err != nil {
				http.Error(resp, "Invalid request", http.StatusBadRequest)
				return
			}
			latlons = append(latlons, LatLon{lat, lon})
		}
		if len(latlons) > MaxInputPoints {
			http.Error(resp, "Request too big", http.StatusRequestEntityTooLarge)
			return
		}
		if readErr == io.EOF {
			break
		}
	}
	elevations, err := getElevations(*demStorage, latlons)
	// TODO: log error
	if err != nil {
		http.Error(resp, "Server error", http.StatusInternalServerError)
		return
	}
	strElevations := make([]string, len(elevations))
	for i, elevation := range elevations {
		// TODO: reduce to 1 digit
		// TODO: handle no-data values
		strElevations[i] = fmt.Sprintf("%.2f", elevation)
	}
	result := strings.Join(strElevations, "\n")
	resp.Write([]byte(result))
}

func main() {
	port := flag.Int("port", 8080, "port to listen")
	host := flag.String("host", "127.0.0.1", "address to bind to")
	dataFile := flag.String("dem", "", "path to file with elevation tile")
	flag.Parse()
	if *dataFile == "" {
		flag.Usage()
		os.Exit(1)
	}
	storage, err := openDemStorage(*dataFile)
	if err != nil {
		log.Fatal(err)
	}
	demStorage = &storage
	http.HandleFunc("/", handleRequest)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", *host, *port), nil))
}
