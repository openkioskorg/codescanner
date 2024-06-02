/* Daemon for using serial QR/barcode scanner devices.
   Copyright (C) 2023  Digilol OÜ

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as
   published by the Free Software Foundation, either version 3 of the
   License, or (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>. */

package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"log"
	"os"

	cs "gitlab.com/openkiosk/codescanner"
)

// By default ignore scans until start command
var isStarted = false

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: ./codescannerd config.yaml")
	}
	conf := parseConfig(os.Args[1])

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
			// Scanner disconnected
			log.Fatal(err)
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
