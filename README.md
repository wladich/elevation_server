Building from source
====================
Tested to work on Ubuntu 20.04.1 in Docker container
```
apt update
apt install -y golang git
apt install -y liblz4-dev

export GOPATH=$HOME/go
mkdir -p $GOPATH/src/github.com/wladich
cd $GOPATH/src/github.com/wladich
git clone https://github.com/wladich/elevation_server.git

cd elevation_server/cmd/elevation_server
go build
cd ../make_data
go get
go build
```
Built executables are at paths:
* cmd/elevation_server/elevation_server
* cmd/make_data/make_data

Installation
============
Place binaries to directory of your choice, for example to `/usr/local/bin`.

Usage
=====
Prepare data
------------
1. Grab DEM files in hgt format with 3 arc-second resolution. The best available source is http://viewfinderpanoramas.org/dem3.html.
Place all *.hgt files to one directory.
2. run `make_data -hgt <PATH_TO_HGT_FILES> -out dem_tiles`
3. You should find two files: `dem_tiles` and `dem_tiles.idx`

Start server
------------
Run `elevation_server -dem <PATH_TO_dem_tiles>`

Other options:
```
  -dem string
        path to file with elevation tile
  -host string
        address to bind to (default "127.0.0.1")
  -port int
        port to listen (default 8080)
  -threads int
        maximum number of concurrently served requests (default 10)
```

API
---
The server uses HTTP protocol.
The only endpoint is "/", request type is POST.
The request consists of a list of latitude-longitude coordinate pairs.
Coordinates are floating point numbers in decimal format, separated with a space character.
First number in the pair is latitude, second is longitude.
Pairs are separated from each other with newline character (\n).

Response contains the list of elevations as floating point numbers in decimal format separated with newlines.


Error is being returned in following cases:
* request body size exceeds 250000 bytes
* request contains more than 10000 points
* request has invalid format

Example:
```
curl https://elevation.example.com --data-binary $'49.18148 16.35126\n49.13545 16.33615'
375.60
259.470
```
