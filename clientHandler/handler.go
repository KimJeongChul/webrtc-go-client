package clientHandler

import (
	"net/url"
	"os"
	"strconv"
	"sync"

	"github.com/KimJeongChul/webrtc-go-client/common"
	"github.com/KimJeongChul/webrtc-go-client/message"
	"github.com/gofrs/uuid"
	"github.com/gorilla/websocket"
	"github.com/notedit/gst"
	"github.com/pion/webrtc/v2"
)

type ClientHandler struct {
	clientConfig   ClientConfigJson
	mediaRoomID    string
	userID         string
	handleID       string
	ws             *websocket.Conn
	msMsg          chan message.ResponseMessage
	webrtcAPI      *webrtc.API
	candidatesMux  *sync.Mutex
	interrupt      chan os.Signal
	elementEncSink *gst.Element
}

//Configuration Type
type ClientConfigJson struct {
	MediaServerAddr string `json:"mediaServerAddr"`
	MediaServerPort int    `json:"mediaServerPort"`
	VideoPath       string `json:"videoPath"`
}

func (c *ClientHandler) Initialize(configJson ClientConfigJson, mediaRoomID string, interrupt chan os.Signal) int {
	c.clientConfig = configJson
	c.mediaRoomID = mediaRoomID
	c.interrupt = interrupt

	uuidv4, err := uuid.NewV4()
	if err != nil {
		common.LogE("clientHandler", "Initialize", "Failed to generate UUID:"+err.Error())
	}

	c.userID = uuidv4.String()
	c.handleID = c.makeID(8)

	// Media Server URL
	mediaServerHost := c.clientConfig.MediaServerAddr + ":" + strconv.Itoa(c.clientConfig.MediaServerPort)
	mediaServerURL := url.URL{Scheme: "ws", Host: mediaServerHost, Path: "/ws"}

	c.ws, _, err = websocket.DefaultDialer.Dial(mediaServerURL.String(), nil)
	if err != nil {
		common.LogE("clientHandler", "Initialize", "Websocket connection error:"+err.Error())
		return -1
	}
	defer c.ws.Close()

	// Websocket Message
	c.msMsg = make(chan message.ResponseMessage, 1024) // 미디어서버 웹소켓 채널

	// Recv Handler
	go c.wsRecvHandler()

	// Message Handler
	go c.wsHandleMessage()

	// Create WebRTC API
	c.webrtcAPI = c.generateWebRTCAPI()
	c.candidatesMux = &sync.Mutex{}

	var encPipelineDefinition string
	if c.clientConfig.VideoPath == "" {
		encPipelineDefinition = "videotestsrc ! video/x-raw ! x264enc speed-preset=ultrafast tune=zerolatency key-int-max=20 ! video/x-h264,stream-format=byte-stream ! appsink name=encSink"
	} else {
		encPipelineDefinition = "filesrc location=" + c.clientConfig.VideoPath + " ! decodebin ! x264enc speed-preset=ultrafast tune=zerolatency key-int-max=20 ! video/x-h264,stream-format=byte-stream ! appsink name=encSink"
	}

	common.LogE("clientHandler", "Initalize", "encPipelineDefinition : "+encPipelineDefinition)

	// Launch GStreamer Pipeline
	encPipeline, err := gst.ParseLaunch(encPipelineDefinition)
	if err != nil {
		common.LogE("clientHandler", "Initalize", "Failed to create GStreamer Pipeline:"+err.Error())
	}

	c.elementEncSink = encPipeline.GetByName("encSink")

	encPipeline.SetState(gst.StatePlaying)

	return 0
}

func (c *ClientHandler) wsRecvHandler() {
	for {
		var res message.ResponseMessage
		err := c.ws.ReadJSON(&res)
		if err != nil {
			common.LogE("clientHandler", "wsRecvHandler", "Failed to read JSON Message:"+err.Error())
			break
		}
		c.msMsg <- res
	}
}

func (c *ClientHandler) wsHandleMessage() {
	for {
		msg := <-c.msMsg
		switch {
		case msg.Method == "SDP":
		}
	}
}

func (c *ClientHandler) Start() {
	c.createPublisherChannel()

	for {
		select {
		case <-c.interrupt:
			c.ws.Close()
			return
		}
	}
}
