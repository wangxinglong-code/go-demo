package common

import (
	"bytes"
	"strconv"
	"testing"
)

//go test -bench=Add* -benchmem -run=none
func BenchmarkAddStringWithBuffer(b *testing.B) {
	hello := "hello"
	world := "world"
	for i := 0; i < 1000; i++ {
		var buffer bytes.Buffer
		buffer.WriteString(hello)
		buffer.WriteString(",")
		buffer.WriteString(world)
		buffer.WriteString("aaa !")
		buffer.WriteString(strconv.Itoa(10000))
		_ = buffer.String()
	}
}

func BenchmarkAddStringWithBufferJoin(b *testing.B) {
	for i := 0; i < 1000; i++ {
		BufferJoin([]string{"hello", ",", "world", "aaa !", strconv.Itoa(10000)})
	}
}
