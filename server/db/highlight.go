package db

import (
	"fmt"

	"github.com/zerozwt/blivehl/server/bs"
	"github.com/zerozwt/blivehl/server/engine"
	"github.com/zerozwt/blivehl/server/utils"
	"go.etcd.io/bbolt"
)

func hlRoomBucketKey(roomId int) []byte {
	return []byte(fmt.Sprintf("hl_%d", roomId))
}

func hlLiveBucketKey(liveId int64) []byte {
	return []byte(fmt.Sprint(liveId))
}

func SaveHightlight(ctx *engine.Context, req *bs.CommitHighlightRequest) error {
	return gDB.Update(func(tx *bbolt.Tx) error {
		roomBucket, _ := tx.CreateBucketIfNotExists(hlRoomBucketKey(req.RoomID))
		liveBucket, _ := roomBucket.CreateBucketIfNotExists(hlLiveBucketKey(req.LiveID))
		userID, _ := utils.GetCtxUser(ctx)
		userBucket, _ := liveBucket.CreateBucketIfNotExists([]byte(userID))
		userBucket.Put(i64reverseKey(req.Time), encodeValue(bs.TimelineEntry{Time: req.Time, Comment: req.Comment}))
		return nil
	})
}

func QueryHightlight(ctx *engine.Context, roomId int, liveId int64) ([]bs.TimelineEntry, error) {
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

		userID, _ := utils.GetCtxUser(ctx)
		userBucket := liveBucket.Bucket([]byte(userID))
		if userBucket == nil {
			return nil
		}

		userBucket.ForEach(func(k, v []byte) error {
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

func QueryAllTimeline(roomId int, liveId int64) ([]bs.AdminTimelineEntry, error) {
	ret := []bs.AdminTimelineEntry{}
	err := gDB.View(func(tx *bbolt.Tx) error {
		roomBucket := tx.Bucket(hlRoomBucketKey(roomId))
		if roomBucket == nil {
			return nil
		}

		liveBucket := roomBucket.Bucket(hlLiveBucketKey(liveId))
		if liveBucket == nil {
			return nil
		}

		liveBucket.ForEachBucket(func(k []byte) error {
			bucket := liveBucket.Bucket(k)
			if bucket == nil {
				return nil
			}
			bucket.ForEach(func(k, v []byte) error {
				item := bs.TimelineEntry{}
				err := decodeValue(v, &item)
				if err != nil {
					return err
				}
				adminItem := bs.AdminTimelineEntry{Author: string(k)}
				adminItem.Time = item.Time
				adminItem.Comment = item.Comment
				ret = append(ret, adminItem)
				return nil
			})
			return nil
		})

		return nil
	})
	return ret, err
}
