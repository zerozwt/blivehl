package bs

type TimelineRequest struct {
	RoomID int   `form:"room_id"`
	LiveID int64 `form:"live_id"`
}

type TimelineEntry struct {
	Time    int64  `json:"time"`
	Comment string `json:"comment"`
}

type TimelineResponse struct {
	Timeline []TimelineEntry `json:"timeline"`
}

type AdminTimelineEntry struct {
	TimelineEntry
	Author string `json:"author"`
}

type AdminTimelineResponse struct {
	Timeline []AdminTimelineEntry `json:"timeline"`
}
