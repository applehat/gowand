# Gowand

## Description

`Gowand` is a Go package that provides functionality for handling and processing "Casts" from a [MagiQuest](https://en.wikipedia.org/wiki/MagiQuest) wand. 

This package is based partially on work done by `rveach` in his [pywand](https://gitlab.com/rveach/wand) package.

`Gowand` is not associated with MagiQuest or Creative Kingdoms in any way, and is not intended to be used for any commercial purposes. MagiQuest is a registered trademark of Great Wolf Resorts, Inc.

## Pre-requisites

`Gowand` requires the following:
- LIRC (Linux Infrared Remote Control) installed and configured 
  - This is installed by default on many popular Linux distributions
- An IR receiver connected to your computer and configured in LIRC

## Installation 

To install the `gowand` package, you can use `go get`:

```bash
go get github.com/applehat/gowand
```

## Usage

The most basic usage of `gowand` is to create a new wand and start listening for wand casts. The following example shows how to do this:

```go
package main

import (
	"github.com/applehat/gowand"
)

func main() {
	// Get a new wand
	wand := gowand.Wand()
	wand.Start()

	wand.OnCast(func(wc gowand.WandCast) {
		// Dump the wand cast data to the console
        fmt.Printf("%+v\n", wc)
	})
}

```

You can alternatively handle the channel that wand casts are sent to yourself. The following example shows how to do this:

```go
package main

import (
    "github.com/applehat/gowand"
)

func main() {
    // Get a new wand
    wand := gowand.Wand()
    wand.Start()

    // Get the channel that wand casts are sent to
    wandCastChan := wand.Chan()

    // Loop forever, handling wand casts
    for {
        select {
        case wc := <-wandCastChan:
            // Dump the wand cast data to the console
            fmt.Printf("%+v\n", wc)
        }
    }
}
```

## Configuration

By default, `gowand` will assume your IR device is at `/dev/lirc0`. If this is not the case, you can specify the device before starting:

```go
wand := gowand.Wand()
wand.SetDevice("/dev/lirc1")
wand.Start()
```

## Hardware

I uased a Raspberry Pi running Raspbian for testing, with an 1838 IR receiver connected to GPIO pin 18. 

Add the following to your `/boot/config.txt` file to enable the IR receiver (requires a reboot):

```
dtoverlay=gpio-ir,gpio_pin=18
```

Then on the Pi, the wiring is:

IR Receiver      | Raspberry Pi
---------------- | ------------
GND              | GND
VCC              | 3.3V
OUT              | GPIO 18

There is a lot of information on the internet suggested that a 36KHz IR Receiver will work better then the more common 38KHz IR Receiver, tho I have been able to somewhat reliably get the 38KHz receiver to work.
