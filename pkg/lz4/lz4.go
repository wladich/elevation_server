package lz4

// #cgo CFLAGS: -O3 -I/usr/include
// #include "lz4.h"
// #include "lz4hc.h"
// #cgo LDFLAGS: -llz4
import "C"
import (
	"errors"
	"unsafe"
)

func byteSliceToCharPointer(b []byte) *C.char {
	if len(b) == 0 {
		return (*C.char)(unsafe.Pointer(nil))
	}
	return (*C.char)(unsafe.Pointer(&b[0]))
}

func CompressBound(size int) int {
	return int(C.LZ4_compressBound(C.int(size)))
}

func CompressHigh(source []byte, compressionLevel int) []byte {
	dest := make([]byte, CompressBound(len(source)))
	n := C.LZ4_compress_HC(
		byteSliceToCharPointer(source),
		byteSliceToCharPointer(dest),
		C.int(len(source)),
		C.int(len(dest)),
		C.int(compressionLevel))
	if n == 0 {
		panic("Unexpected error while compressing with LZ4")
	}
	return dest[:n]
}

func Decompress(source, dest []byte) (int, error) {
	n := C.LZ4_decompress_safe(
		byteSliceToCharPointer(source),
		byteSliceToCharPointer(dest),
		C.int(len(source)),
		C.int(len(dest)))
	if n < 0 {
		return int(n), errors.New("error decompressing lz4 block")
	}
	return int(n), nil
}
