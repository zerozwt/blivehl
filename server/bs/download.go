package bs

type DownloadRequest struct {
	RoomID int   `form:"room_id"`
	LiveID int64 `form:"live_id"`
}
