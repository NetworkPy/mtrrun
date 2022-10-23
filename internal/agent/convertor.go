package agent

import (
	"bytes"
	"encoding/binary"
)

// Convert string array to byte array.
// Unused function, but it implements for support labels in metrics
func convertStrToByte(str []string) ([]byte, error) {
	buf := new(bytes.Buffer)

	if len(str) == 0 {
		return buf.Bytes(), nil
	}

	for i := 0; i < len(str); i++ {
		byteStr := []byte(str[i])
		err := binary.Write(buf, binary.LittleEndian, uint64(len(byteStr)))

		if err != nil {
			return nil, err
		}

		_, err = buf.Write(byteStr)

		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

// Convert byte array to string array.
// Unused function, but it implements for support labels in metrics
func convertByteToStr(b []byte) []string {
	str := make([]string, 0)

	if len(b) == 0 {
		return str
	}

	digit := make([]byte, 8, 8)
	var count int

	for i := 0; i < len(b); i++ {
		digit[count] = b[i]
		count++

		if count == 8 {
			count = 0

			// Give length for next word
			l := int(binary.LittleEndian.Uint64(digit))

			// i - current position
			// l - word's length
			// added 1, so 'i' located in current element
			str = append(str, string(b[i+1:i+l+1]))

			i += l
		}
	}

	return str
}
