package clientHandler

import (
	"context"
	"math/rand"

	"github.com/KimJeongChul/webrtc-go-client/common"
	"github.com/pion/webrtc/v2"
	"github.com/pion/webrtc/v2/pkg/media"
)

var config = webrtc.Configuration{
	ICEServers: []webrtc.ICEServer{
		{
			URLs: []string{"stun:stun.l.google.com:19302"},
		},
	},
}

// WebRTC Configuration
func (c *ClientHandler) generateWebRTCAPI() *webrtc.API {
	// Add to Setting Engine
	se := webrtc.SettingEngine{}
	se.SetTrickle(true)

	// Add to MediaEngine
	me := webrtc.MediaEngine{}
	me.RegisterCodec(webrtc.NewRTPOpusCodec(webrtc.DefaultPayloadTypeOpus, 48000)) // Opus 코덱 등록
	me.RegisterCodec(webrtc.NewRTPH264Codec(webrtc.DefaultPayloadTypeH264, 90000)) // H264 코덱 등록
	api := webrtc.NewAPI(webrtc.WithMediaEngine(me), webrtc.WithSettingEngine(se))

	return api
}

func (c *ClientHandler) createPublisherChannel() {
	pc, err := c.webrtcAPI.NewPeerConnection(config)
	if err != nil {
		common.LogE("clientHandler", "createPublisherChannel", "Failed to create Peer Connection:"+err.Error())
	}

	iceConnectedCtx, iceConnectedCtxCancel := context.WithCancel(context.Background())
	pendingCandidates := make([]*webrtc.ICECandidate, 0)

	// Create Video Track
	videoTrack, addTrackErr := pc.NewTrack(webrtc.DefaultPayloadTypeH264, rand.Uint32(), "video", "pion")
	if addTrackErr != nil {
		common.LogE("clientHandler", "createPublisherChannel", "Failed to create video track:"+addTrackErr.Error())
	}

	// Add Video Track
	if _, addTrackErr = pc.AddTrack(videoTrack); err != nil {
		common.LogE("clientHandler", "createPublisherChannel", "Failed to add video track:"+addTrackErr.Error())
	}

	go func() {
		<-iceConnectedCtx.Done()

		for {
			sample, err := c.elementEncSink.PullSample()
			if err != nil {
				common.LogE("clientHandler", "createPublisherChannel", "Failed to pull Gstreamer element enc sink:"+err.Error())
			}

			samples := uint32(90000 * (float32(sample.Duration) / 1000000000))

			// Write Sample H264 Data
			if h264Err := videoTrack.WriteSample(media.Sample{Data: sample.Data, Samples: samples}); h264Err != nil {
				panic(h264Err)
			}
		}
	}()

	// ICE Connection Stat
	pc.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		if connectionState.String() == "connected" {
			//TODO: Subscribe

			iceConnectedCtxCancel()
		}
	})

	pc.OnICECandidate(func(ice *webrtc.ICECandidate) {
		if ice != nil {
			c.candidatesMux.Lock()
			defer c.candidatesMux.Unlock()

			desc := pc.RemoteDescription()
			if desc == nil {
				pendingCandidates = append(pendingCandidates, ice)
			} else {
				candidate := ice.ToJSON()
				common.LogD("clientHandler", "createPublisherChannel", "candidate:"+candidate.Candidate)
			}
		}
	})

	// Create Offer
	offer, err := pc.CreateOffer(nil)
	if err != nil {
		common.LogE("clientHandler", "createPublisherChannel", "Failed to create offer:"+err.Error())
	}

	// Local Description 설정
	err = pc.SetLocalDescription(offer)
	if err != nil {
		common.LogE("clientHandler", "createPublisherChannel", "Failed to create local description:"+err.Error())
	}

	// TODO: Send ICE Candidate
	for _, ice := range pendingCandidates {
		candidate := ice.ToJSON()
		common.LogD("clientHandler", "createPublisherChannel", "candidate:"+candidate.Candidate)
	}

	// TODO: Send SDP
}
