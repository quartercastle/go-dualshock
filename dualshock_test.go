package dualshock

import (
	"testing"
)

type fakeDevice struct{}

func (f fakeDevice) Read(b []byte) (int, error) {
	copy(b, []byte{
		1, 134, 127, 128, 126, 8, 4, 88, 255, 0, 141, 219, 9, 188, 255, 4, 0,
		167, 255, 250, 6, 212, 31, 51, 254, 0, 0, 0, 0, 0, 27, 0, 0, 1, 252,
		129, 115, 70, 27, 130, 62, 97, 32, 0, 128, 0, 0, 0, 128, 0, 0, 0, 0,
		128, 0, 0, 0, 128, 0, 0, 0, 0, 128, 0,
	})
	return 64, nil
}

func TestDualshock(t *testing.T) {
	controller := New(fakeDevice{})

	result := make(chan State, 1)
	defer close(result)

	controller.Listen(func(state State) {
		controller.Close()
		result <- state
	})

	if r := <-result; !r.L2 {
		t.Errorf("Invalid state L2 should be true; got %v", r.L2)
	}
}

func BenchmarkDualshock(b *testing.B) {
	controller := New(fakeDevice{})

	result := make(chan State, 1)

	go controller.Listen(func(state State) {
		result <- state
	})

	for n := 0; n < b.N; n++ {
		<-result
	}
}
