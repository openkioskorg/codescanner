package codescanner

import (
	"errors"
	"time"

	"go.bug.st/serial"
)

var ErrEOF = errors.New("EOF")

type CodeScannerConfig struct {
	// Max length of scanned data/string
	BuffLen uint

	// Name of port in /dev/tty*
	PortName string

	Debounce time.Duration
}

// Represents a QR code/barcode scanner with a serial interface
type CodeScannerDevice struct {
	serial.Port
	buffLen  uint
	debounce time.Duration
}

func Init(conf *CodeScannerConfig) (*CodeScannerDevice, error) {
	port, err := serial.Open(conf.PortName, &serial.Mode{})
	if err != nil {
		return nil, err
	}
	if err := port.ResetInputBuffer(); err != nil {
		return nil, err
	}
	return &CodeScannerDevice{Port: port, buffLen: conf.BuffLen,
		debounce: conf.Debounce}, nil
}

// Blocks until something is scanned. Exits once input is received.
func (d *CodeScannerDevice) Scan() ([]byte, error) {
	buff := make([]byte, d.buffLen)
	n, err := d.Read(buff)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, ErrEOF
	}
	return buff, nil
}

type CodeScannerResult struct {
	BytesRead []byte
	Err       error
}

// Keeps scanning and shoves the read bytes down a channel
func (d *CodeScannerDevice) ScanWithHandler(ch chan<- *CodeScannerResult) {
	buff := make([]byte, d.buffLen)
	for {
		n, err := d.Read(buff)
		if err != nil {
			ch <- &CodeScannerResult{Err: err}
		}
		if n == 0 {
			ch <- &CodeScannerResult{Err: ErrEOF}
		}

		ch <- &CodeScannerResult{BytesRead: buff}

		time.Sleep(d.debounce)

		// Don't read things that were scanned during the debounce period
		if err := d.ResetInputBuffer(); err != nil {
			ch <- &CodeScannerResult{Err: err}
		}
	}
}
