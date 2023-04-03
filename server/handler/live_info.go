package handler

import (
	"github.com/zerozwt/blivehl/server/bs"
	"github.com/zerozwt/blivehl/server/engine"
	"github.com/zerozwt/blivehl/server/service"
)

func init() {
	engine.RegisterApi("/prepare", getPrepareInfo)
	engine.RegisterApi("/live_list", getLiveList)
}

func getPrepareInfo(req *bs.PrepareRequest) (*bs.PrepareResponse, error) {
	return service.GetLiveInfoService().GetPrepareInfo(req)
}

func getLiveList(req *bs.LiveListRequest) (*bs.LiveListResponse, error) {
	return service.GetLiveInfoService().GetLiveList(req)
}
