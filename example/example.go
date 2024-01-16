package main

import (
	"github.com/applehat/gowand"
)

func main() {
	// Get a new wand
	wand := gowand.Wand()

	// You can override the default threshold and IR device if you want / need.
	//wand.SetThreshold(410)
	//wand.SetIRDevice("/dev/lirc0")

	// Start the wand listener
	wand.Start()

	// Now there are two ways to get the wand data.
	// You can either use the WandCast channel, or you can use the OnCast() function.
	wand.OnCast(func(wc gowand.WandCast) {
		println(wc.WandID)
		println(wc.MotionID)
		println(wc.Timeout)
		println(wc.PulseLen)
	})

	// Or if you don't want to use a go routine, you can get the channel and use it yourself.
	for {
		select {
		case wc := <-wand.Chan():
			println(wc.WandID)
			println(wc.MotionID)
			println(wc.Timeout)
			println(wc.PulseLen)
		}
	}
}
