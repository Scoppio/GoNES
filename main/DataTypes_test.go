package main

import (
	"testing"
)

func TestFnEquals(t *testing.T) {
	fn := IMM

	if got := FnEquals(fn, IMM); got != true {
		t.Errorf("Expected: %t, got: %t", true, got)
	}
}
