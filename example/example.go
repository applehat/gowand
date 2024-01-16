package main

import (
	"fmt"
	"time"

	"github.com/applehat/gowand"
)

func main() {
	// Get a new wand
	wand := gowand.Wand()

	wandIds := []string{}
	magnitudes := []string{}

	// You can override the default threshold and IR device if you want / need.
	wand.SetIRDevice("/dev/lirc0")
	//wand.Verbose()

	// Set a callback for when a wand is cast
	wand.OnCast(func(wc gowand.WandCast) {
		fmt.Printf("Wand Cast Detected: %v\n", wc)
		if !contains(wandIds, wc.WandID) {
			wandIds = append(wandIds, wc.WandID)
			fmt.Println("New Wand ID: " + wc.WandID)
		}
		if !contains(magnitudes, wc.Magnitude) {
			magnitudes = append(magnitudes, wc.Magnitude)
			fmt.Println("New Motion ID: " + wc.Magnitude)
		}
	})

	// Alternatively, you can get the wandCast channel and read from it directly
	// go func() {
	// 	for {
	// 		select {
	// 		case wc := <-wand.Chan():
	// 			println(wc.WandID)
	// 			println(wc.Magnitude)
	// 			println(wc.Timeout)
	// 			println(wc.PulseLen)
	// 		}
	// 	}
	// }()

	go func() {
		for {
			fmt.Print("\033[H\033[2J")
			fmt.Printf("Current Wand IDs: %v\n", wandIds)
			fmt.Printf("Current Motion IDs: %v\n", magnitudes)
			time.Sleep(5 * time.Second)
		}
	}()

	// Start the wand listener
	// This will block the main thread
	wand.Start()
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
