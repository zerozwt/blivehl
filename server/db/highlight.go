package db

import (
	"fmt"

	"github.com/zerozwt/blivehl/server/bs"
	"go.etcd.io/bbolt"
)

func hlRoomBucketKey(roomId int) []byte {
	return []byte(fmt.Sprintf("hl_%d", roomId))
}

func hlLiveBucketKey(liveId int64) []byte {
	return []byte(fmt.Sprint(liveId))
}

func SaveHightlight(req *bs.CommitHighlightRequest) error {
	return gDB.Update(func(tx *bbolt.Tx) error {
		roomBucket, _ := tx.CreateBucketIfNotExists(hlRoomBucketKey(req.RoomID))
		liveBucket, _ := roomBucket.CreateBucketIfNotExists(hlLiveBucketKey(req.LiveID))
		liveBucket.Put(i64reverseKey(req.Time), encodeValue(bs.TimelineEntry{Time: req.Time, Comment: req.Comment}))
		return nil
	})
}

func QueryHightlight(roomId int, liveId int64) ([]bs.TimelineEntry, error) {
	ret := []bs.TimelineEntry{}

	err := gDB.View(func(tx *bbolt.Tx) error {
		roomBucket := tx.Bucket(hlRoomBucketKey(roomId))
		if roomBucket == nil {
			return nil
		}

		liveBucket := roomBucket.Bucket(hlLiveBucketKey(liveId))
		if liveBucket == nil {
			return nil
		}

		liveBucket.ForEach(func(k, v []byte) error {
			item := bs.TimelineEntry{}
			err := decodeValue(v, &item)
			if err != nil {
				return err
			}
			ret = append(ret, item)
			return nil
		})
		return nil
	})

	return ret, err
}
