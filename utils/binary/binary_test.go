package binaryutil_test

import (
	"testing"

	binaryutil "github.com/cn-poe-community/poedat/utils/binary"
)

func TestFloat32(t *testing.T) {
	b := []byte{0xdb, 0x0f, 0x49, 0x40}
	var pi float32 = 3.141592653589793
	var diff float32 = 0.000001

	n, err := binaryutil.LittleEndian.Float32(b)
	if err != nil {
		t.Error(err)
	}
	if n <= pi && n+diff > pi ||
		n >= pi && n < pi+diff {
		//ok
	} else {
		t.Errorf("Float32 expected %v, got %v", pi, n)
	}
}
