package handler

import (
	"github.com/zerozwt/blivehl/server/bs"
	"github.com/zerozwt/blivehl/server/engine"
	"github.com/zerozwt/blivehl/server/service"
)

func init() {
	engine.RegisterApi("/basic_info", basicInfo)
	engine.RegisterApi("/room_list", recentRooms)
}

func basicInfo(req *bs.BasicInfoRequest) (*bs.BasicInfoResponse, error) {
	return service.GetRoomInfoService().GetBasicInfo(req)
}

func recentRooms(req *bs.RoomListRequest) (*bs.RoomListResponse, error) {
	return service.GetRoomInfoService().GetRecentRooms()
}

func InitHandlers() {}
