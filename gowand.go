package gowand

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type MagiWand struct {
	Threshold int
	IRDevice  string
	wandCast  chan WandCast
}

type WandCast struct {
	WandID   string
	MotionID string
	Timeout  int
	PulseLen int
}

func Wand() *MagiWand {
	return &MagiWand{
		Threshold: 410,
		IRDevice:  "/dev/lirc0",
		wandCast:  make(chan WandCast),
	}
}

// SetThreshold sets the threshold for determining 1s and 0s
func (w *MagiWand) SetThreshold(threshold int) {
	w.Threshold = threshold
}

// SetIRDevice sets the IR device to use
func (w *MagiWand) SetIRDevice(device string) {
	w.IRDevice = device
}

// Start starts the IR listener
func (w *MagiWand) Start() error {
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
		// First, split out the timeout
		pulses, timeout := parseLine(line)

		binary := ""

		if len(pulses) == 56 || len(pulses) == 112 {
			// We *might* have a valid signal
			// compares pulse width to determine binary values
			for _, pulse := range pulses {
				if pulse >= w.Threshold {
					binary += "1"
				} else {
					binary += "0"
				}

			}

			// we now have the 56 bit binary string
			// 0:8 is the zero
			// 8:32 is the wand
			// 32:56 is the motion

			wandId, _ := binaryToHex(binary[8:32])
			// motion ID is the REST of the binary string
			motionId, _ := binaryToHex(binary[32:])

			w.wandCast <- WandCast{
				WandID:   wandId,
				MotionID: motionId,
				Timeout:  timeout,
				PulseLen: len(pulses),
			}

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

func parseLine(line string) ([]int, int) {
	// Example line: +304 -853 +534 -596 +298 -859 +243 -886 +560 -597 +243 -889 +557 -595 +534 -597 +559 -597 +559 -571 +560 -596 +559 # timeout 23266
	s := strings.Split(line, "# timeout")
	timeout, err := strconv.Atoi(strings.Trim(s[1], " "))
	if err != nil {
		//fmt.Println("Error getting timeout:", err)
		timeout = 0
	}

	// Next, split out the pulses
	pulseandspace := strings.Split(s[0], " ")
	var pulses []int
	// seperate the pulses (+500) from the spaces (-500)
	for i, pulse := range pulseandspace {
		if i%2 == 0 {
			pulseS := strings.Split(pulse, "+")
			pulseInt, err := strconv.Atoi(pulseS[1])
			if err == nil {
				pulses = append(pulses, pulseInt)
			}
		}
	}

	return pulses, timeout
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
