package message

// ResponseMessage ...
type ResponseMessage struct {
	Method string `json:"method"`
	RoomID string `json:"roomID"`
	Status int    `json:"status"`
}
