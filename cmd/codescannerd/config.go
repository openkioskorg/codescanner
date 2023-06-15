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
	"fmt"
	"log"
	"os"

	cs "gitlab.com/openkiosk/codescanner"
	"gopkg.in/yaml.v3"
)

type daemonConfig struct {
	// TODO: expose underlying serial library config
	Device *cs.CodeScannerConfig `yaml:"device"`
	Mqtt   brokerConfig          `yaml:"mqtt"`
}

func parseConfig() (conf daemonConfig) {
	file, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatal("Failed to read config: ", err)
	}
	if err := yaml.Unmarshal(file, &conf); err != nil {
		log.Fatal("Failed to unmarshal yaml: ", err)
	}
	fmt.Printf("%v", *conf.Device)
	return
}
