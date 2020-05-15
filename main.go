package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	math2 "math"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
	"strings"
)

const MaxInputSize = 250000
const MaxInputPoints = 10000

var demStorage *DemStorage

func fastFloatToString(f float64) string {
	var s string
	i := int(math2.Round(f * 100))
	if i >= 100 || i <= -100 {
		s = strconv.Itoa(i)
		l := len(s)
		return s[:l-2] + "." + s[l-2:]
	}
	sign := 1
	if i < 0 {
		sign = -1
		i *= -1
	}
	s = strconv.Itoa(i)
	if i < 10 {
		s = "0.0" + s
	} else {
		s = "0." + s
	}
	if sign == -1 {
		s = "-" + s
	}
	return s
}

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
			last := len(line) - 1
			if line[last] == '\n' {
				line = line[:last]
			}
			spaceIndex := strings.Index(line, " ")
			if spaceIndex == -1 {
				http.Error(resp, "Invalid request", http.StatusBadRequest)
				return
			}
			latStr := line[:spaceIndex]
			lonStr := line[spaceIndex+1:]
			lat, err := strconv.ParseFloat(latStr, 64)
			if err != nil {
				http.Error(resp, "Invalid request", http.StatusBadRequest)
				return
			}
			lon, err := strconv.ParseFloat(lonStr, 64)
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
		log.Printf("Failed to get elevation: %s", err)
		return
	}
	strElevations := make([]string, len(elevations))
	for i, elevation := range elevations {
		// TODO: reduce to 1 digit
		// TODO: handle no-data values
		strElevations[i] = fastFloatToString(elevation)
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
