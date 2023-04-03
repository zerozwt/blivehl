package db

import (
	"bytes"
	"encoding/binary"

	jsoniter "github.com/json-iterator/go"
	"go.etcd.io/bbolt"
)

var gDB *bbolt.DB

func InitDB(file string) error {
	var err error
	gDB, err = bbolt.Open(file, 0644, nil)
	return err
}

func encodeValue(value any) []byte {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	ret, _ := json.Marshal(value)
	return ret
}

func decodeValue(data []byte, value any) error {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	return json.Unmarshal(data, value)
}

func i64reverseKey(value int64) []byte {
	buf := bytes.Buffer{}
	binary.Write(&buf, binary.BigEndian, value)
	data := buf.Bytes()
	for idx := range data {
		data[idx] = ^data[idx]
	}
	return data
}
