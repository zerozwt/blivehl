package db

import (
	"sort"
	"strconv"
	"time"

	"go.etcd.io/bbolt"
)

const (
	MAX_RECENT_LIVE_ROOM = 8
)

type LiveRoomInfo struct {
	UID            int64  `json:"uid"`
	RoomID         int    `json:"room_id"`
	StreamerName   string `json:"name"`
	StreamerIcon   string `json:"icon"`
	LastAccessTime int64  `json:"lat"`
}

func roomid2key(roomId int) []byte {
	return []byte(strconv.Itoa(roomId))
}

func putRoomInfo(bucket *bbolt.Bucket, info *LiveRoomInfo) {
	bucket.Put(roomid2key(info.RoomID), encodeValue(info))
}

func SaveRecentLiveRoom(uid int64, roomId int, name, icon string) error {
	return gDB.Update(func(tx *bbolt.Tx) error {
		infoMap := make(map[int]*LiveRoomInfo)

		bucket, _ := tx.CreateBucketIfNotExists([]byte("recent"))
		if err := bucket.ForEach(func(k, v []byte) error {
			item := LiveRoomInfo{}
			err := decodeValue(v, &item)
			if err != nil {
				return err
			}
			infoMap[item.RoomID] = &item
			return nil
		}); err != nil {
			return err
		}

		current := &LiveRoomInfo{
			UID:            uid,
			RoomID:         roomId,
			StreamerName:   name,
			StreamerIcon:   icon,
			LastAccessTime: time.Now().Unix(),
		}

		putRoomInfo(bucket, current)
		infoMap[roomId] = current

		roomIds := []int{}
		for k := range infoMap {
			roomIds = append(roomIds, k)
		}

		sort.Slice(roomIds, func(i, j int) bool {
			return infoMap[roomIds[i]].LastAccessTime > infoMap[roomIds[j]].LastAccessTime
		})

		if len(roomIds) > MAX_RECENT_LIVE_ROOM {
			for _, id := range roomIds[MAX_RECENT_LIVE_ROOM:] {
				bucket.Delete(roomid2key(id))
			}
		}
		return nil
	})
}

func GetRecentRooms() ([]*LiveRoomInfo, error) {
	ret := []*LiveRoomInfo{}
	err := gDB.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("recent"))
		if bucket == nil {
			return nil
		}
		return bucket.ForEach(func(k, v []byte) error {
			item := LiveRoomInfo{}
			err := decodeValue(v, &item)
			if err != nil {
				return err
			}
			ret = append(ret, &item)
			return nil
		})
	})
	if len(ret) > 0 {
		sort.Slice(ret, func(i, j int) bool {
			return ret[i].LastAccessTime > ret[j].LastAccessTime
		})
	}
	return ret, err
}
