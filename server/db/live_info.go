package db

import (
	"fmt"

	"github.com/zerozwt/blivehl/server/bs"
	"go.etcd.io/bbolt"
)

func liveListBucketKey(roomId int) []byte {
	return []byte(fmt.Sprintf("live_%d", roomId))
}

func liveItemDataKey(liveId int64) []byte {
	return i64reverseKey(liveId)
}

func SaveLiveInfo(roomId int, info *bs.BasicLiveInfo) error {
	return gDB.Update(func(tx *bbolt.Tx) error {
		bucket, _ := tx.CreateBucketIfNotExists(liveListBucketKey(roomId))
		bucket.Put(liveItemDataKey(info.LiveID), encodeValue(info))
		return nil
	})
}

func QueryLiveInfo(roomId int, until int64, limit int) ([]*bs.BasicLiveInfo, error) {
	ret := []*bs.BasicLiveInfo{}

	err := gDB.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(liveListBucketKey(roomId))
		if bucket == nil {
			return nil
		}
		cursor := bucket.Cursor()
		kvInit := func() ([]byte, []byte) {
			if until == 0 {
				return cursor.First()
			}
			return cursor.Seek(i64reverseKey(until))
		}
		for k, v := kvInit(); k != nil && limit > 0; k, v = cursor.Next() {
			limit--
			item := bs.BasicLiveInfo{}
			if err := decodeValue(v, &item); err != nil {
				return err
			}
			ret = append(ret, &item)
		}
		return nil
	})

	return ret, err
}
