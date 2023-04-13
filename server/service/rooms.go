package service

import (
	"github.com/zerozwt/blivehl/server/bs"
	"github.com/zerozwt/blivehl/server/db"
	"github.com/zerozwt/blivehl/server/logger"
)

type RoomService struct{}

var gRoomService *RoomService = &RoomService{}

func GetRoomInfoService() *RoomService {
	return gRoomService
}

func (s *RoomService) GetBasicInfo(req *bs.BasicInfoRequest) (*bs.BasicInfoResponse, error) {
	// get room info from fetcher
	info, err := GetRoomInfoFetcher().Get(req.RoomID)
	if err != nil {
		return nil, err
	}

	ret := &bs.BasicInfoResponse{}
	ret.Streamer.UID = info.Base.Uid
	ret.Streamer.Name = info.Liver.Base.Name
	ret.Streamer.Icon = info.Liver.Base.Icon

	// async save live room cover
	ch_cover := make(chan bool)
	go func() {
		defer close(ch_cover)
		SaveFileAndGetLocalUrl(ret.Streamer.UID, info.Base.Cover)
	}()

	// save stramer icon
	ret.Streamer.Icon, err = SaveFileAndGetLocalUrl(ret.Streamer.UID, ret.Streamer.Icon)
	if err != nil {
		return nil, err
	}

	// store room info into db (allow fail)
	err = db.SaveRecentLiveRoom(ret.Streamer.UID, req.RoomID, ret.Streamer.Name, ret.Streamer.Icon)
	if err != nil {
		logger.WARN("save streamer %d info to db failed: %v", ret.Streamer.UID, err)
	}

	<-ch_cover
	return ret, nil
}

func (s *RoomService) GetRecentRooms() (*bs.RoomListResponse, error) {
	ret := &bs.RoomListResponse{
		List: []bs.RecentLiveRoom{},
	}

	list, err := db.GetRecentRooms()
	if err != nil {
		return nil, err
	}

	for _, item := range list {
		ret.List = append(ret.List, bs.RecentLiveRoom{
			ID:   item.RoomID,
			Name: item.StreamerName,
			Icon: item.StreamerIcon,
		})
	}

	return ret, nil
}
