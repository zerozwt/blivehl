package service

import (
	"fmt"
	"time"

	"github.com/zerozwt/blivehl/server/bs"
	"github.com/zerozwt/blivehl/server/db"
	"github.com/zerozwt/blivehl/server/logger"
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
	err = db.SaveLiveInfo(req.RoomID, &ret.BasicLiveInfo)
	if err != nil {
		logger.ERROR("save live stream info to db for room %d failed: %v", req.RoomID, err)
	}

	return ret, nil
}
