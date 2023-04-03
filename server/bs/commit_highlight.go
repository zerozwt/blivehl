package bs

type CommitHighlightRequest struct {
	RoomID  int    `json:"room_id"`
	LiveID  int64  `json:"live_id"`
	Time    int64  `json:"ts"`
	Comment string `json:"comment"`
}

type CommitHighlightResponse struct {
}
