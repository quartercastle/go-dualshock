# go-dualshock

[![Version](https://img.shields.io/github/release/kvartborg/go-dualshock.svg)](https://github.com/kvartborg/go-dualshock/releases)
[![Build Status](https://travis-ci.org/kvartborg/go-dualshock.svg?branch=master)](https://travis-ci.org/kvartborg/go-dualshock)
[![GoDoc](https://godoc.org/github.com/kvartborg/go-dualshock?status.svg)](https://godoc.org/github.com/kvartborg/go-dualshock)
[![Go Report Card](https://goreportcard.com/badge/github.com/kvartborg/go-dualshock)](https://goreportcard.com/report/github.com/kvartborg/go-dualshock)


Connect a PS4 DualShock controller with your go program.

### Install
```sh
go get github.com/kvartborg/go-dualshock
```

### Example
```go
package main

import (
    "fmt"
    "log"
    "github.com/karalabe/hid"
    dualshock "github.com/kvartborg/go-dualshock"
)

func main() {
    vendorID, productID := uint16(1356), uint16(1476)
    devices := hid.Enumerate(vendorID, productID)

    if len(devices) == 0 {
        log.Fatal("no dualshock controller where found")
    }

    device, err := devices[0].Open()

    if err != nil {
        log.Fatal(err)
    }

    controller := dualshock.New(device)

    controller.Listen(func(state dualshock.State) {
        fmt.Println(state.Analog.L2)
    })
}
```

### License
This project is licensed under the [MIT License](https://github.com/kvartborg/go-dualshock/blob/master/LICENSE).
