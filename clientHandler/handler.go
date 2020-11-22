package clientHandler

import (
	"github.com/KimJeongChul/webrtc-go-client/common"
	"github.com/gofrs/uuid"
)

type ClientHandler struct {
	clientConfig ClientConfigJson
	mediaRoomID  string
	userID       string
	handleID     string
}

//Configuration Type
type ClientConfigJson struct {
	MediaServerAddr string `json:"mediaServerAddr"`
	MediaServerPort int    `json:"mediaServerPort"`
}

func (c *ClientHandler) Initialize(configJson ClientConfigJson, mediaRoomID string) int {
	c.clientConfig = configJson
	c.mediaRoomID = mediaRoomID

	uuidv4, err := uuid.NewV4()
	if err != nil {
		common.LogE("Initialize", "Failed to generate UUID"+err.Error())
	}

	c.userID = uuidv4.String()
	c.handleID = c.makeID(8)
	return 0
}
