package raw

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Color struct{ R, G, B float32 }

func WriteRawImage(indata []float64, width, height int) []byte {
	pixels := make([]Color, len(indata)/4)
	j := 0
	for i := 0; i < len(indata); i += 4 {
		pixels[j].R = float32(indata[i])
		pixels[j].G = float32(indata[i+1])
		pixels[j].B = float32(indata[i+2])
		// ignore alpha
		j++
	}

	var byteBuffer bytes.Buffer

	fileFormatVersionMajor := 1
	fileFormatVersionMinor := 0

	writeBinaryInt32(&byteBuffer, int32(fileFormatVersionMajor))
	writeBinaryInt32(&byteBuffer, int32(fileFormatVersionMinor))
	writeBinaryInt32(&byteBuffer, int32(width))
	writeBinaryInt32(&byteBuffer, int32(height))

	binary.Write(&byteBuffer, binary.BigEndian, pixels)

	return byteBuffer.Bytes()
}

func writeBinaryInt32(buffer *bytes.Buffer, value int32) {
	if err := binary.Write(buffer, binary.BigEndian, value); err != nil {
		fmt.Println(err)
	}
}
