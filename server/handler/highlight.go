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
)

func init() {
	engine.RegisterApi("/commit", commitHighlight)
	engine.RegisterApi("/timeline", queryTimeline)
	engine.RegisterRawApi("/download", download)
}

func commitHighlight(ctx *engine.Context, req *bs.CommitHighlightRequest) (*bs.CommitHighlightResponse, error) {
	return service.GetHightlightService().Commit(req)
}

func queryTimeline(ctx *engine.Context, req *bs.TimelineRequest) (*bs.TimelineResponse, error) {
	return service.GetHightlightService().Query(req)
}

func download(ctx *engine.Context) {
	r := ctx.RawRequest
	w := ctx.RawResponse
	req := bs.DownloadRequest{}
	if err := engine.DecodeForm(r, &req); err != nil {
		logger.ERROR("decode download request failed: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	csvLines, fileName, err := service.GetHightlightService().GenerateCSVLines(req.RoomID, req.LiveID)

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
