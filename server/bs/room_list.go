package bs

type RoomListRequest struct{}

type RecentLiveRoom struct {
	ID   int    `json:"room_id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type RoomListResponse struct {
	List []RecentLiveRoom `json:"list"`
}
