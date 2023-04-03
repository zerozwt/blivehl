package bs

type LiveListRequest struct {
	RoomID int   `form:"room_id"`
	Until  int64 `form:"until"`
}

type LiveListResponse struct {
	List  []BasicLiveInfo `json:"list"`
	Ended bool            `json:"ended"`
}
