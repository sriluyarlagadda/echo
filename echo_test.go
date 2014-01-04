package echo

import "testing"

var onesComplimentTestData map[int32]uint16 = map[int32]uint16{
	56:    56,
	65536: 1,
	-6:    65529,
	1000:  1000,
}

func TestOnesComplimentFor16(t *testing.T) {

	for i, j := range onesComplimentTestData {
		if n := OnesCompliment(i); n != j {
			t.Errorf("OnesCompliment(%v) = %v but want %d", i, n, j)
		}
	}

}
