package handler

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"

	"github.com/zerozwt/blivehl/server/bs"
	"github.com/zerozwt/blivehl/server/engine"
	"github.com/zerozwt/blivehl/server/logger"
	"github.com/zerozwt/blivehl/server/service"
	"github.com/zerozwt/blivehl/server/utils"
)

var notAdminRsp []byte = []byte(`{"code":1919810,"msg":"not admin","data":{}}`)

func init() {
	engine.RegisterApi("/admin/rooms", adminRoomList, loginChecker, adminChecker)
	engine.RegisterApi("/admin/lives", adminLiveList, loginChecker, adminChecker)
	engine.RegisterApi("/admin/timeline", adminHightlightTimeline, loginChecker, adminChecker)
	engine.RegisterRawApi("/admin/download", adminDownload, loginChecker, adminChecker)
}

func adminChecker(ctx *engine.Context) {
	if !utils.CheckCtxAdmin(ctx) {
		ctx.RawResponse.Header().Set("Content-Type", "application/json")
		ctx.RawResponse.Write(notAdminRsp)
		return
	}
	ctx.Next()
}

func adminRoomList(ctx *engine.Context, req *bs.RoomListRequest) (*bs.RoomListResponse, error) {
	return service.GetRoomInfoService().GetAllLiveRooms()
}

func adminLiveList(ctx *engine.Context, req *bs.AdminLiveListRequest) (*bs.AdminLiveListResponse, error) {
	list, total, err := service.GetLiveInfoService().GetLiveListByPage(req.RoomID, req.Page, req.PageSize)
	return &bs.AdminLiveListResponse{List: list, Total: total}, err
}

func adminHightlightTimeline(ctx *engine.Context, req *bs.TimelineRequest) (*bs.AdminTimelineResponse, error) {
	return service.GetHightlightService().QueryAll(req)
}

func adminDownload(ctx *engine.Context) {
	r := ctx.RawRequest
	w := ctx.RawResponse
	req := bs.DownloadRequest{}
	if err := engine.DecodeForm(r, &req); err != nil {
		logger.ERROR("decode download request failed: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	csvLines, fileName, err := service.GetHightlightService().GenerateCSVLinesForAdmin(req.RoomID, req.LiveID)

	if err != nil {
		logger.ERROR("download generate csv lines failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	csvData := &bytes.Buffer{}
	csvData.Write([]byte{0xEF, 0xBB, 0xBF}) // UTF8 BOM
	err = csv.NewWriter(csvData).WriteAll(csvLines)

	if err != nil {
		logger.ERROR("download generate csv data failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`inline; filename="%s.csv"`, fileName))
	w.Write(csvData.Bytes())
}
