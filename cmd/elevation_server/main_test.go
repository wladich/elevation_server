package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestFastFloatToString(t *testing.T) {
	data := []float64{0, -10.565884654, -1.0001, -0.999999, -0.32, 0, 0.00001, 0.5654, 1.0001, 3}
	for _, v := range data {
		got := fastFloatToString(v)
		expected := fmt.Sprintf("%.2f", v)
		if expected != got {
			t.Errorf("%v: expected=%v, got=%v", v, expected, got)
		}

	}
}

type responseWriterResult struct {
	status int
	body   string
}

type mockResponseWriter struct {
	result *responseWriterResult
}

func newMockResonseWriter() mockResponseWriter {
	var writer mockResponseWriter
	writer.result = &responseWriterResult{}
	return writer
}

func (writer mockResponseWriter) Header() http.Header {
	return http.Header{}
}

func (writer mockResponseWriter) Write(data []byte) (int, error) {
	writer.result.body = string(data)
	return len(data), nil
}

func (writer mockResponseWriter) WriteHeader(statusCode int) {
	writer.result.status = statusCode
}

//func BenchmarkServer(b *testing.B) {
//	dataFile := "path/to/dem_tiles"
//	testsDir := "/home/w/phenom_mnt/projects/elevation_server/benchmark/data2/"
//	var err error
//	demStorage, err = dem.NewReader(dataFile)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer demStorage.Close()
//	files, err := ioutil.ReadDir(testsDir)
//	if err != nil {
//		panic(err)
//	}
//	for _, f := range files {
//		requestBody, err := ioutil.ReadFile(path.Join(testsDir, f.Name()))
//		if err != nil {
//			panic(err)
//		}
//		bodyReader := ioutil.NopCloser(bytes.NewReader(requestBody))
//		request := http.Request{Method: "POST", URL: &url.URL{Path: "/"}, Body: bodyReader}
//		writer := newMockResonseWriter()
//		handleRequest(writer, &request)
//		if writer.result.status != 0 {
//			panic(errors.New(fmt.Sprintf("request status is %v", writer.result.status)))
//		}
//	}
//}
