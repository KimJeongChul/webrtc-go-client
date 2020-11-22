package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/KimJeongChul/webrtc-go-client/clientHandler"
	"github.com/KimJeongChul/webrtc-go-client/common"
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
	configFilePath := flag.String("c", "./clientConfig.json", "set server config file")
	mediaRoomID := flag.String("r", "", "media room id")
	flag.Parse()

	//Load Configuration
	config, err := LoadConfigJson(configFilePath)
	if err != nil {
		common.LogE("main", "Config File:"+*configFilePath+" Load Error.")
		os.Exit(-1)
	}

	//Media Room ID
	if *mediaRoomID == "" {
		common.LogE("main", "Media Room ID is Nil Error.")
		os.Exit(-1)
	}

	clientHandler := clientHandler.ClientHandler{}
	clientHandler.Initialize(config, *mediaRoomID)
}
