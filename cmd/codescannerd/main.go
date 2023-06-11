package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"log"

	cs "gitlab.com/openkiosk/codescanner"
)

// By default ignore scans until start command
var isStarted = false

func main() {
	conf := parseConfig()
	scanner, err := cs.Init(conf.Device)
	if err != nil {
		log.Fatal(err)
	}

	broker, err := newBroker(conf.Mqtt)
	if err != nil {
		log.Fatal("Failed to connect to MQTT broker: ", err)
	}

	for {
		res, err := scanner.Scan()
		if err != nil {
			log.Println(err)
		}
		if isStarted {
			if err := broker.publishScan(context.Background(), b64(res)); err != nil {
				log.Println("Failed to publish event: ", err)
			}
		}
	}
}

// Trim null bytes and base64 encode
func b64(b []byte) string {
	b = bytes.Trim(b, "\x00")
	return base64.StdEncoding.EncodeToString(b)
}
