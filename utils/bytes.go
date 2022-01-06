// Package utils provides some cross used information
package utils

import (
	"encoding/binary"
	"encoding/json"
)

// ByteToMap returns a map from bytes
func ByteToMap(bytes []byte) (map[string]interface{}, error) {
	val := make(map[string]interface{})
	if len(bytes) > 0 {
		err := json.Unmarshal(bytes, &val)
		if err != nil {
			return make(map[string]interface{}), err
		}
	}
	return val, nil
}

// ByteToUint64 return an uint64 from bytes
func ByteToUint64(bytes []byte) uint64 {
	num := uint64(0)
	if len(bytes) > 0 {
		num = binary.BigEndian.Uint64(bytes)
	}
	return num
}

// Itob returns an 8-byte big endian representation of v.
func Itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}
