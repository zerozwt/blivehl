package bs

type PrepareRequest struct {
	RoomID int `form:"room_id"`
}

type BasicLiveInfo struct {
	Title         string `json:"title"`
	LiveStartTime string `json:"live_start_time"`
	Cover         string `json:"cover"`
	LiveID        int64  `json:"live_id"`
}

type PrepareResponse struct {
	BasicLiveInfo
	LiveStatus int   `json:"live_status"`
	LightTs    int64 `json:"light_ts"`
}
