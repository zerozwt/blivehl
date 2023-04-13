package service

import (
	"time"

	dm "github.com/zerozwt/BLiveDanmaku"
	"github.com/zerozwt/blivehl/server/utils"
)

var roomInfoFetcher *utils.Fetcher[int, dm.RoomInfo] = utils.NewFetcher(time.Second*3, dm.GetRoomInfo)

func GetRoomInfoFetcher() *utils.Fetcher[int, dm.RoomInfo] {
	return roomInfoFetcher
}
