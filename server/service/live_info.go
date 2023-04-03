package service

import (
	"fmt"
	"time"

	"github.com/zerozwt/blivehl/server/bs"
	"github.com/zerozwt/blivehl/server/db"
	"github.com/zerozwt/blivehl/server/logger"
)

const (
	LIVE_LIST_PAGE_SIZE = 3
)

var CSTZone *time.Location = time.FixedZone("CST", 8*60*60)

func RenderTimeInCST(ts int64) string {
	CSTTime := time.Unix(ts, 0).In(CSTZone)
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d CST",
		CSTTime.Year(), CSTTime.Month(), CSTTime.Day(),
		CSTTime.Hour(), CSTTime.Minute(), CSTTime.Second())
}

type LiveInfoService struct{}

var gLiveInfo *LiveInfoService = &LiveInfoService{}

func GetLiveInfoService() *LiveInfoService {
	return gLiveInfo
}

func (s *LiveInfoService) GetPrepareInfo(req *bs.PrepareRequest) (*bs.PrepareResponse, error) {
	// get room info from fetcher
	info, err := GetRoomInfoFetcher().GetRoomInfo(req.RoomID)
	if err != nil {
		return nil, err
	}

	ret := &bs.PrepareResponse{
		LiveStatus: info.Base.LiveStatus,
		BasicLiveInfo: bs.BasicLiveInfo{
			Title:         info.Base.Title,
			LiveStartTime: RenderTimeInCST(info.Base.LiveStartTime),
			Cover:         "",
			LiveID:        info.Base.LiveStartTime,
		},
		LightTs: time.Now().Unix() - info.Base.LiveStartTime,
	}

	ret.Cover, err = SaveFileAndGetLocalUrl(info.Base.Uid, info.Base.Cover)
	if err != nil {
		logger.ERROR("get live stream cover for room %d failed: %v", req.RoomID, err)
	}

	// save live info
	if ret.LiveStatus != 0 {
		err = db.SaveLiveInfo(req.RoomID, &ret.BasicLiveInfo)
		if err != nil {
			logger.ERROR("save live stream info to db for room %d failed: %v", req.RoomID, err)
		}
	}

	return ret, nil
}

func (s *LiveInfoService) GetLiveList(req *bs.LiveListRequest) (*bs.LiveListResponse, error) {
	ret := bs.LiveListResponse{
		List: []bs.BasicLiveInfo{},
	}

	// fetch current live info if until==0
	var currentLive *bs.BasicLiveInfo
	if req.Until == 0 {
		info, err := GetRoomInfoFetcher().GetRoomInfo(req.RoomID)
		if err == nil && info.Base.LiveStatus != 0 {
			currentLive = &bs.BasicLiveInfo{
				Title:         info.Base.Title,
				LiveStartTime: RenderTimeInCST(info.Base.LiveStartTime),
				Cover:         "",
				LiveID:        info.Base.LiveStartTime,
			}
			currentLive.Cover, _ = SaveFileAndGetLocalUrl(info.Base.Uid, info.Base.Cover)
		} else if err != nil {
			logger.WARN("query current live status for room %d in live_list failed: %v", req.RoomID, err)
		}
	}

	// query db
	list, err := db.QueryLiveInfo(req.RoomID, req.Until, LIVE_LIST_PAGE_SIZE)
	if err != nil {
		return nil, err
	}
	ret.Ended = len(list) < LIVE_LIST_PAGE_SIZE

	// combine two of them
	if currentLive != nil {
		idxInList := -1
		for idx, item := range list {
			if item.LiveID == currentLive.LiveID {
				idxInList = idx
				break
			}
		}
		if idxInList < 0 {
			list = append([]*bs.BasicLiveInfo{currentLive}, list...)
		}
	}

	if len(list) > LIVE_LIST_PAGE_SIZE {
		list = list[:LIVE_LIST_PAGE_SIZE]
	}

	for _, item := range list {
		ret.List = append(ret.List, *item)
	}

	return &ret, nil
}
