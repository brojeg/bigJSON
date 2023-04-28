package benchmark

import (
	"bigJson/internal/controller"
	"os"
	"strconv"
	"testing"

	jsoniter "github.com/json-iterator/go"
)

func BenchmarkBufferSize(b *testing.B) {
	bufferSizes := []int{1024, 2048, 4096, 8192, 16384, 32768}

	for _, bufferSize := range bufferSizes {
		b.Run("BufferSize"+strconv.Itoa(bufferSize), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				file, err := os.Open("../bigJSON.json")
				if err != nil {
					b.Fatal("Error opening file:", err)
				}
				defer file.Close()

				iter := jsoniter.Parse(jsoniter.ConfigFastest, file, bufferSize)
				controller.Process(iter, "", "", "", []string{"tomato"})
				file.Close()
			}
		})
	}
}
