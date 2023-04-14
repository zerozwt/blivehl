package service

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/zerozwt/blivehl/server/bs"
	"github.com/zerozwt/blivehl/server/db"
	"github.com/zerozwt/blivehl/server/engine"
	"github.com/zerozwt/blivehl/server/logger"
)

type HighlightService struct{}

var gHighlight *HighlightService = &HighlightService{}

func GetHightlightService() *HighlightService {
	return gHighlight
}

func (s *HighlightService) Commit(ctx *engine.Context, req *bs.CommitHighlightRequest) (*bs.CommitHighlightResponse, error) {
	return &bs.CommitHighlightResponse{}, db.SaveHightlight(ctx, req)
}

func (s *HighlightService) Query(ctx *engine.Context, req *bs.TimelineRequest) (*bs.TimelineResponse, error) {
	ret := bs.TimelineResponse{
		Timeline: []bs.TimelineEntry{},
	}

	list, err := db.QueryHightlight(ctx, req.RoomID, req.LiveID)
	if err != nil {
		return nil, err
	}

	ret.Timeline = list

	return &ret, nil
}

func (s *HighlightService) GenerateCSVLines(ctx *engine.Context, roomId int, liveId int64) ([][]string, string, error) {
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
	info, err := GetRoomInfoFetcher().Get(roomId)
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
	list, err := db.QueryHightlight(ctx, roomId, liveId)
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

func (s *HighlightService) QueryAll(req *bs.TimelineRequest) (*bs.AdminTimelineResponse, error) {
	list, err := db.QueryAllTimeline(req.RoomID, req.LiveID)
	if err != nil {
		return nil, err
	}

	sort.Slice(list, func(i, j int) bool { return list[i].Time < list[j].Time })

	authorsMap := map[string]bool{}
	for _, item := range list {
		authorsMap[item.Author] = true
	}
	authors := []string{}
	for key := range authorsMap {
		authors = append(authors, key)
	}

	users, err := db.BatchGetUser(authors)
	if err != nil {
		return nil, err
	}

	for idx := range list {
		tmp := list[idx]
		if user, ok := users[tmp.Author]; ok {
			tmp.Author = user.Name
			list[idx] = tmp
		}
	}

	return &bs.AdminTimelineResponse{Timeline: list}, nil
}

func (s *HighlightService) GenerateCSVLinesForAdmin(roomId int, liveId int64) ([][]string, string, error) {
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
	info, err := GetRoomInfoFetcher().Get(roomId)
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
	listRsp, err := s.QueryAll(&bs.TimelineRequest{RoomID: roomId, LiveID: liveId})
	if err != nil {
		return nil, "", err
	}
	list := listRsp.Timeline
	if len(list) > 0 {
		addLine("高能时间", "高能内容", "作者")
		for _, item := range list {
			addLine(s.renderTimeStamp(item.Time), item.Comment, item.Author)
		}
	}

	return ret, strings.Join(fileNameParts, "_"), nil
}
