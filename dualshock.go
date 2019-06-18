package dualshock

import (
	"encoding/binary"
	"io"
)

// Controller describes the reference to the hardware device
type Controller struct {
	reader    io.Reader
	queue     chan []byte
	errors    chan error
	interrupt chan int
}

// DPad is the data structure describing the joysticks on the controller
type DPad struct {
	X, Y int
}

// TrackPad is the touch surface in the middle of the controller
type TrackPad struct {
	ID     int
	Active bool
	X, Y   int
}

// Motion contains information about acceleration in x, y and z axis
type Motion struct {
	X, Y, Z int16
}

// Orientation describes roll, yaw and pith of the controller
type Orientation struct {
	Roll, Yaw, Pitch int16
}

// State is the overall state of the controller, the controller will output 254
// states within a second.
type State struct {
	L1, L2, L3                      bool
	R1, R2, R3                      bool
	Up, Down, Left, Right           bool
	Cross, Circle, Square, Triangle bool
	Share, Options, PSHome          bool
	Timestamp, BatteryLevel         int
	LeftDPad                        DPad
	RightDPad                       DPad
	Motion                          Motion
	Orientation                     Orientation

	// TrackPad is true if its pressed
	TrackPad bool

	// TrackPad0 contains info relating to a touch event
	TrackPad0 TrackPad

	// TrackPad1 only register touches if TrackPad0 also does, this enables
	// multi touch functionality
	TrackPad1 TrackPad
	// Analog describes the analog position of buttons, on the PS4 controller its
	// only L2 and R2 which has analog output as well as digital.
	Analog struct{ L2, R2 int }
}

// transform reads a slice of bytes and turns them into a valid state for the
// controller
func transform(b []byte) State {
	return State{
		L1:       (b[6] & 0x01) != 0,
		L2:       (b[6] & 0x04) != 0,
		L3:       (b[6] & 0x40) != 0,
		R1:       (b[6] & 0x02) != 0,
		R2:       (b[6] & 0x08) != 0,
		R3:       (b[6] & 0x80) != 0,
		Up:       (b[5]&15) == 0 || (b[5]&15) == 1 || (b[5]&15) == 7,
		Down:     (b[5]&15) == 3 || (b[5]&15) == 4 || (b[5]&15) == 5,
		Left:     (b[5]&15) == 5 || (b[5]&15) == 6 || (b[5]&15) == 7,
		Right:    (b[5]&15) == 1 || (b[5]&15) == 2 || (b[5]&15) == 3,
		Cross:    (b[5] & 32) != 0,
		Circle:   (b[5] & 64) != 0,
		Square:   (b[5] & 16) != 0,
		Triangle: (b[5] & 128) != 0,
		Share:    (b[6] & 0x10) != 0,
		Options:  (b[6] & 0x20) != 0,
		PSHome:   (b[7] & 1) != 0,
		TrackPad: (b[7] & 2) != 0,
		TrackPad0: TrackPad{
			ID:     int(b[35] & 0x7f),
			Active: (b[35] >> 7) == 0,
			X:      int(((b[37] & 0x0f) << 4) | b[36]),
			Y:      int(b[38]<<4 | ((b[37] & 0xf0) >> 4)),
		},
		TrackPad1: TrackPad{
			ID:     int(b[39] & 0x7f),
			Active: (b[39] >> 7) == 0,
			X:      int(((b[41] & 0x0f) << 4) | b[40]),
			Y:      int(b[42]<<4 | ((b[41] & 0xf0) >> 4)),
		},
		LeftDPad: DPad{
			X: int(b[1]),
			Y: int(b[2]),
		},
		RightDPad: DPad{
			X: int(b[3]),
			Y: int(b[4]),
		},
		Motion: Motion{
			Y: int16(binary.LittleEndian.Uint16(b[13:])),
			X: -int16(binary.LittleEndian.Uint16(b[15:])),
			Z: -int16(binary.LittleEndian.Uint16(b[17:])),
		},
		Orientation: Orientation{
			Roll:  -int16(binary.LittleEndian.Uint16(b[19:])),
			Yaw:   int16(binary.LittleEndian.Uint16(b[21:])),
			Pitch: int16(binary.LittleEndian.Uint16(b[23:])),
		},
		Analog: struct{ L2, R2 int }{
			L2: int(b[8]),
			R2: int(b[9]),
		},
		Timestamp:    int(b[7] >> 2),
		BatteryLevel: int(b[12]),
	}
}

// New returns a new controller which transforms input from the device to a valid
// controller state
func New(reader io.Reader) *Controller {
	c := &Controller{
		reader,
		make(chan []byte),
		make(chan error),
		make(chan int),
	}
	go c.read()
	return c
}

// read transforms data from the io.Reader and pushes it to the queue of
// states
func (c *Controller) read() {
	for {
		select {
		case <-c.interrupt:
			close(c.errors)
			close(c.queue)
			return
		default:
			b := make([]byte, 64)
			n, err := c.reader.Read(b)

			if err != nil {
				c.errors <- err
				continue
			}

			c.queue <- b[:n]
		}
	}
}

// Listen for controller state changes
func (c *Controller) Listen(handle func(State)) {
	for {
		select {
		case <-c.interrupt:
			return
		default:
			handle(transform(<-c.queue))
		}
	}
}

// Errors returns a channel of reader errors
func (c *Controller) Errors() <-chan error {
	return c.errors
}

// Close the listener
func (c *Controller) Close() {
	close(c.interrupt)
}
