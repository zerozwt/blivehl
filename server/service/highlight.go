package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/zerozwt/blivehl/server/bs"
	"github.com/zerozwt/blivehl/server/db"
	"github.com/zerozwt/blivehl/server/logger"
)

type HighlightService struct{}

var gHighlight *HighlightService = &HighlightService{}

func GetHightlightService() *HighlightService {
	return gHighlight
}

func (s *HighlightService) Commit(req *bs.CommitHighlightRequest) (*bs.CommitHighlightResponse, error) {
	return &bs.CommitHighlightResponse{}, db.SaveHightlight(req)
}

func (s *HighlightService) Query(req *bs.TimelineRequest) (*bs.TimelineResponse, error) {
	ret := bs.TimelineResponse{
		Timeline: []bs.TimelineEntry{},
	}

	list, err := db.QueryHightlight(req.RoomID, req.LiveID)
	if err != nil {
		return nil, err
	}

	ret.Timeline = list

	return &ret, nil
}

func (s *HighlightService) GenerateCSVLines(roomId int, liveId int64) ([][]string, string, error) {
	ret := [][]string{}
	addLine := func(cells ...any) {
		line := []string{}
		for _, item := range cells {
			line = append(line, fmt.Sprint(item))
		}
		ret = append(ret, line)
	}

	meta := false
	fileNameParts := []string{}

	// get streamer info
	info, err := GetRoomInfoFetcher().GetRoomInfo(roomId)
	if err != nil {
		logger.WARN("get streamer info for room %d failed: %v", roomId, err)
		fileNameParts = append(fileNameParts, fmt.Sprint(roomId))
	} else {
		addLine("主播", info.Liver.Base.Name)
		addLine("UID", info.Base.Uid)
		addLine("直播间", roomId)
		fileNameParts = append(fileNameParts, s.filterFileName(info.Liver.Base.Name))
		meta = true
	}

	// get live info
	liveInfo, err := db.QueryLiveInfo(roomId, liveId, 1)
	if err != nil || len(liveInfo) == 0 {
		logger.WARN("get live info for room %d failed: %v", roomId, err)
	} else {
		addLine("直播标题", liveInfo[0].Title)
		addLine("开播时间", liveInfo[0].LiveStartTime)
		fileNameParts = append(fileNameParts, s.filterFileName(liveInfo[0].Title))
		meta = true
	}
	tm := time.Unix(liveId, 0).In(CSTZone)
	fileNameParts = append(fileNameParts, fmt.Sprintf("%04d%02d%02d_%02d%02d%02d",
		tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second()))

	if meta {
		addLine("")
	}

	// get timeline
	list, err := db.QueryHightlight(roomId, liveId)
	if err != nil {
		return nil, "", err
	}
	if len(list) > 0 {
		addLine("高能时间", "高能内容")
		for idx := len(list) - 1; idx >= 0; idx-- {
			addLine(s.renderTimeStamp(list[idx].Time), list[idx].Comment)
		}
	}

	return ret, strings.Join(fileNameParts, "_"), nil
}

func (s *HighlightService) renderTimeStamp(ts int64) string {
	ret := fmt.Sprintf("%02d分%02d秒", (ts%3600)/60, ts%60)
	if ts > 3600 {
		ret = fmt.Sprintf("%d小时", ts/3600) + ret
	}
	return ret
}

func (s *HighlightService) filterFileName(name string) string {
	ret := ""

	for _, ch := range name {
		switch ch {
		case '<', '>', ':', '"', '\'', '/', '\\', '|', '?', '*':
			ret += "_"
		default:
			ret += string(ch)
		}
	}

	return ret
}
