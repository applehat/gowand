package gowand

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type MagiWand struct {
	IRDevice string
	wandCast chan WandCast
	verbose  bool
}

type WandCast struct {
	WandID    string // Hex Wand ID
	Magnitude string // Hex Motion ID
	Checksum  string // Hex Checksum
	Timeout   int    // Timeout from ir-ctl (might be useful for something)
}

func Wand() *MagiWand {
	return &MagiWand{
		IRDevice: "/dev/lirc0",
		wandCast: make(chan WandCast),
		verbose:  false,
	}
}

func (w *MagiWand) log(str ...string) {
	if w.verbose {
		for _, s := range str {
			fmt.Printf("[GoWand] %s\n", s)
		}
	}
}

func (w *MagiWand) Verbose() {
	w.verbose = true
	w.log("Verbose Logging Enabled")
}

// SetIRDevice sets the IR device to use
func (w *MagiWand) SetIRDevice(device string) {
	w.IRDevice = device
}

// Start starts the IR listener
func (w *MagiWand) Start() error {
	w.log("Starting IR listener")
	out := make(chan string)
	command := exec.Command("ir-ctl", "-r", "-d", w.IRDevice)
	stdout, err := command.StdoutPipe()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(stdout)
	go func() {
		for scanner.Scan() {
			out <- scanner.Text()
		}
		close(out)
	}()

	if err := command.Start(); err != nil {
		return err
	}

	for line := range out {
		w.log("Reading Line: " + line)
		// First, split out the timeout
		// Example line: +304 -853 +534 -596 +298 -859 +243 -886 +560 -597 +243 -889 +557 -595 +534 -597 +559 -597 +559 -571 +560 -596 +559 # timeout 23266
		s := strings.Split(line, "# timeout")
		timeout, err := strconv.Atoi(strings.Trim(s[1], " "))
		if err != nil {
			//fmt.Println("Error getting timeout:", err)
			timeout = 0
		}

		// Next, split out the pulses (So we have + and -)
		pulseandspace := strings.Split(s[0], " ")

		binary := ""
		plus := 0
		minus := 0

		for i, pulse := range pulseandspace {
			if strings.Contains(pulse, "+") {
				pulseS := strings.Split(pulse, "+")
				pulseInt, _ := strconv.Atoi(pulseS[1])
				plus = pulseInt
			}
			if strings.Contains(pulse, "-") {
				pulseS := strings.Split(pulse, "-")
				pulseInt, _ := strconv.Atoi(pulseS[1])
				minus = pulseInt
			}
			if i%2 == 1 {
				// We have a pulse and a space
				// Adding these together we can get the duty cycle
				total := plus + minus
				// if plus >= 1/3 of total, it's a 1
				if plus >= (total / 3) {
					binary += "1"
				} else {
					binary += "0"
				}
			}
		}

		// Note -- we are technically potentially missing the last pulse, but it's not needed for the wand cast because its all
		// checksum stuff that we don't care about.

		if len(binary) == 56 || len(binary) == 112 {
			// we now have the 56 bit binary string (or 112 that we will ignore the second half of)
			// 0:8 is always zero -- 8 bits
			// 8:39 is the wand -- 31 Bits
			// 39:48 is the magnitude -- 9 bits

			if binary[0:8] != "00000000" {
				w.log("Invalid binary string:" + binary)
				// This is not a wand cast. It's something else.
				continue
			}

			wandId, _ := binaryToHex(binary[8:39])
			magnitude, _ := binaryToHex(binary[39:48])
			checksum, _ := binaryToHex(binary[48:56])

			// verify checksum

			cast := WandCast{
				WandID:    "0x" + wandId,
				Magnitude: "0x" + magnitude,
				Checksum:  "0x" + checksum,
				Timeout:   timeout,
			}

			w.log("Wand Cast Found! Wand: " + wandId + " Motion:" + magnitude)

			w.wandCast <- cast
		} else {
			w.log("Invalid binary length:" + strconv.Itoa(len(binary)) + " -- " + binary)
		}
	}
	return nil
}

// OnCast defines a function to be called when a wand is cast
// This creates a goroutine to listen for wand casts
func (w *MagiWand) OnCast(function func(WandCast)) {
	go func() {
		for {
			select {
			case wc := <-w.wandCast:
				function(wc)
			}
		}
	}()
}

// Chan returns the channel that wand casts are sent to
// This allows you to create your own goroutine to listen for wand casts
func (w *MagiWand) Chan() chan WandCast {
	return w.wandCast
}

func binaryToHex(binary string) (string, error) {
	// Convert binary to integer
	intValue, err := strconv.ParseInt(binary, 2, 64)
	if err != nil {
		return "", err
	}

	// Convert integer to hexadecimal
	hexValue := fmt.Sprintf("%x", intValue)

	return hexValue, nil
}
