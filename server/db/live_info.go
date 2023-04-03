package db

import (
	"bytes"
	"encoding/binary"
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"github.com/zerozwt/blivehl/server/bs"
	"go.etcd.io/bbolt"
)

func liveListBucketKey(roomId int) []byte {
	return []byte(fmt.Sprintf("live_%d", roomId))
}

func liveItemDataKey(liveId int64) []byte {
	buf := bytes.Buffer{}
	binary.Write(&buf, binary.BigEndian, liveId)
	data := buf.Bytes()
	for idx := range data {
		data[idx] = ^data[idx]
	}
	return data
}

func SaveLiveInfo(roomId int, info *bs.BasicLiveInfo) error {
	return gDB.Update(func(tx *bbolt.Tx) error {
		bucket, _ := tx.CreateBucketIfNotExists(liveListBucketKey(roomId))
		json := jsoniter.ConfigCompatibleWithStandardLibrary
		data, _ := json.Marshal(info)
		bucket.Put(liveItemDataKey(info.LiveID), data)
		return nil
	})
}
