package main

import (
	"fmt"
	"go.bug.st/serial"
	"log"
)

var ErrEOF = errors.New("EOF")

type CodeScannerConfig struct {
	// Max length of scanned data/string
	BuffLen uint

	// Name of port in /dev/tty*
	PortName string

	// TODO: field for delimiter
}

// Represents a QR code/barcode scanner with a serial interface
type CodeScannerDevice struct {
	serial.Port
	buffLen uint
	// TODO: there should be a time.Duration member to be used for
	// debouncing.
}

// Blocks until something is scanned. Exits once input is received.
func (d *CodeScannerDevice) Scan() ([]byte, error) {
	buff := make([]byte, d.buffLen)
	n, err := port.Read(buff)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, ErrEOF
	}
	return n, nil
}

// Keeps scanning and shoves the read bytes down a channel
func (d *CodeScannerDevice) ScanWithHandler(ch chan<- []byte) {
}
