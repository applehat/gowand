package main

import (
	"github.com/applehat/gowand"
)

func main() {
	// Get a new wand
	wand := gowand.Wand()

	// You can override the default threshold and IR device if you want / need.
	wand.SetIRDevice("/dev/lirc0")
	wand.Verbose()

	// Set a callback for when a wand is cast
	wand.OnCast(func(wc gowand.WandCast) {
		fmt.Println(wc.WandID)
		fmt.Println(wc.MotionID)
		//fmt.Println(wc.Timeout)
		//fmt.Println(wc.PulseLen)

	})

	// Alternatively, you can get the wandCast channel and read from it directly
	// go func() {
	// 	for {
	// 		select {
	// 		case wc := <-wand.Chan():
	// 			println(wc.WandID)
	// 			println(wc.MotionID)
	// 			println(wc.Timeout)
	// 			println(wc.PulseLen)
	// 		}
	// 	}
	// }()

	fmt.Println("Starting Server")

	// Start the wand listener
	wand.Start()

}
