package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/KimJeongChul/webrtc-go-client/clientHandler"
	"github.com/KimJeongChul/webrtc-go-client/common"
	"github.com/notedit/gst"
)

// Read Configuration File
func LoadConfigJson(inFile *string) (clientHandler.ClientConfigJson, error) {
	var parsedResult clientHandler.ClientConfigJson
	file, err := os.Open(*inFile)
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&parsedResult)
	if err != nil {
		log.Fatal(err)
	}
	return parsedResult, err
}

func main() {
	// Start GStreamer
	go gst.MainLoopNew().Run()

	configFilePath := flag.String("c", "./clientConfig.json", "set server config file")
	mediaRoomID := flag.String("r", "", "media room id")
	flag.Parse()

	//Load Configuration
	config, err := LoadConfigJson(configFilePath)
	if err != nil {
		common.LogE("main", "main", "Config File:"+*configFilePath+" Load Error.")
		os.Exit(-1)
	}

	//Media Room ID
	if *mediaRoomID == "" {
		common.LogE("main", "main", "Media Room ID is Nil Error.")
		os.Exit(-1)
	}

	// Register interrupt
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	clientHandler := clientHandler.ClientHandler{}
	result := clientHandler.Initialize(config, *mediaRoomID, interrupt)
	if result == -1 {
		common.LogE("main", "main", "Client Handler failed to initialize")
		os.Exit(-1)
	}

	clientHandler.Start()
}
