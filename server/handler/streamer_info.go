package handler

import (
	"github.com/zerozwt/blivehl/server/bs"
	"github.com/zerozwt/blivehl/server/engine"
	"github.com/zerozwt/blivehl/server/service"
)

func init() {
	engine.RegisterApi("/room/basic", basicInfo, loginChecker)
	engine.RegisterApi("/room/list", recentRooms, loginChecker)
}

func basicInfo(ctx *engine.Context, req *bs.BasicInfoRequest) (*bs.BasicInfoResponse, error) {
	return service.GetRoomInfoService().GetBasicInfo(ctx, req)
}

func recentRooms(ctx *engine.Context, req *bs.RoomListRequest) (*bs.RoomListResponse, error) {
	return service.GetRoomInfoService().GetRecentRooms(ctx)
}

func InitHandlers() {}
