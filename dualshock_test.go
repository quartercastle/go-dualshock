package dualshock_test

import (
	"testing"

	dualshock "github.com/kvartborg/go-dualshock"
)

type fakeDevice struct{}

func (f fakeDevice) Read(b []byte) (int, error) {
	copy(b, []byte{
		1, 134, 127, 128, 126, 8, 4, 88, 255, 0, 141, 219, 9, 188, 255, 4, 0,
		167, 255, 250, 6, 212, 31, 51, 254, 0, 0, 0, 0, 0, 27, 0, 0, 1, 252,
		129, 115, 70, 27, 130, 62, 97, 32, 0, 128, 0, 0, 0, 128, 0, 0, 0, 0,
		128, 0, 0, 0, 128, 0, 0, 0, 0, 128, 0,
	})
	return 0, nil
}

func TestDualshock(t *testing.T) {
	controller := dualshock.New(fakeDevice{})

	result := make(chan dualshock.State, 1)
	controller.Listen(func(state dualshock.State) {
		controller.Close()
		result <- state
	})

	if r := <-result; !r.L2 {
		t.Errorf("Invalid state L2 should be true; got %v", r.L2)
	}
}
