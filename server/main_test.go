package main

import (
	"testing"
)

func TestRefBit(t *testing.T) {
	result := refBit(0, 130)
	if result != 0 {
		t.Error(result)
	}
}

func TestByteToBinaryDigit(t *testing.T) {
	result := ByteToBinaryDigit(129)
	test := "10000001"
	if result != test {
		t.Error(result)
	}
}

func TestByteToBinaryDigit2(t *testing.T) {
	result := ByteToBinaryDigit(130)
	test := "10000010"
	if result != test {
		t.Error(result)
	}
}
