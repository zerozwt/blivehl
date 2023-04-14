package bs

type BasicInfoRequest struct {
	RoomID     int    `form:"room_id"`
	SaveRecent string `form:"save_recent"`
}

func (req BasicInfoRequest) NeedSaveRecent() bool {
	return req.SaveRecent == "true"
}

type BasicInfoResponse struct {
	Streamer struct {
		UID  int64  `json:"uid"`
		Name string `json:"name"`
		Icon string `json:"icon"`
	} `json:"streamer"`
}
