package main

import (
	"log"

	cs "gitlab.com/openkiosk/codescanner"

)

func main() {
	conf := &cs.CodeScannerConfig{
		PortName: "/dev/ttyACM0",
		BuffLen: 100,
	}
	scanner, err := cs.Init(conf)
	if err != nil {
		log.Fatal(err)
	}
	res, err := scanner.Scan()
	if err != nil {
		log.Println(err)
	}
	log.Println(string(res))
}
