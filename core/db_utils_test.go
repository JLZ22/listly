package core

import (
	"testing"
)

func TestBoolToBytesAndBack(t *testing.T) {
	tests := []struct {
		input bool
	}{
		{true},
		{false},
	}

	for _, tt := range tests {
		b := boolToBytes(tt.input)
		got := bytesToBool(b)
		if got != tt.input {
			t.Errorf("boolToBytes/bytesToBool failed for input %v: got %v", tt.input, got)
		}
	}

	// Test bytesToBool with empty slice and other values
	if bytesToBool([]byte{}) != false {
		t.Error("bytesToBool([]byte{}) should return false")
	}
	if bytesToBool([]byte{0}) != false {
		t.Error("bytesToBool([]byte{0}) should return false")
	}
	if bytesToBool([]byte{2}) != false {
		t.Error("bytesToBool([]byte{2}) should return false")
	}
	if bytesToBool([]byte{1}) != true {
		t.Error("bytesToBool([]byte{1}) should return true")
	}
}

func TestItobAndBtoi(t *testing.T) {
	values := []int{0, 1, 255, 256, 1024, -1, -1000}

	for _, v := range values {
		b := itob(v)
		got := btoi(b)
		if got != v {
			t.Errorf("itob/btoi failed for %d: got %d", v, got)
		}
	}
}

func TestIntsToBytesAndBack(t *testing.T) {
	testCases := [][]int{
		{},
		{1},
		{1, 2, 3},
		{0, 255, 256, 1024},
		{-1, -2, -3},
	}

	for _, ints := range testCases {
		b := intsToBytes(ints)
		got := bytesToInts(b)
		if !equalIntSlices(got, ints) {
			t.Errorf("intsToBytes/bytesToInts failed for %v: got %v", ints, got)
		}
	}
}

func equalIntSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
