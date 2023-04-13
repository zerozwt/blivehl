package handler

import (
	"github.com/zerozwt/blivehl/server/bs"
	"github.com/zerozwt/blivehl/server/engine"
	"github.com/zerozwt/blivehl/server/service"
)

func init() {
	engine.RegisterApi("/live/prepare", getPrepareInfo)
	engine.RegisterApi("/live/list", getLiveList)
}

func getPrepareInfo(ctx *engine.Context, req *bs.PrepareRequest) (*bs.PrepareResponse, error) {
	return service.GetLiveInfoService().GetPrepareInfo(req)
}

func getLiveList(ctx *engine.Context, req *bs.LiveListRequest) (*bs.LiveListResponse, error) {
	return service.GetLiveInfoService().GetLiveList(req)
}
