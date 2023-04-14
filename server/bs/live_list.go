package bs

type LiveListRequest struct {
	RoomID int   `form:"room_id"`
	Until  int64 `form:"until"`
}

type LiveListResponse struct {
	List  []BasicLiveInfo `json:"list"`
	Ended bool            `json:"ended"`
}

type AdminLiveListRequest struct {
	RoomID   int `form:"room_id"`
	Page     int `form:"page"`
	PageSize int `form:"size"`
}

type AdminLiveListResponse struct {
	List  []BasicLiveInfo `json:"list"`
	Total int             `json:"total"`
}
