package main

import (
	"runtime/debug"
	"testing"
)

func assertNil(t *testing.T, got interface{}) {
	if got != nil {
		t.Errorf("Expected variable to be: nil")
		t.Log(string(debug.Stack()))
	}
}

func assertTrue(t *testing.T, got bool) {
	if !got {
		t.Errorf("Expected variable to be: %t, got: %t", true, got)
		t.Log(string(debug.Stack()))
	}
}

func assertFalse(t *testing.T, got bool) {
	if got {
		t.Errorf("Expected variable to be: %t, got: %t", false, got)
		t.Log(string(debug.Stack()))
	}
}

func assertEqualsB(t *testing.T, expect byte, got byte) {
	if expect != got {
		t.Errorf("Expected variable to be: %x, got: %x", expect, got)
		t.Log(string(debug.Stack()))
	}
}

func assertEqualsW(t *testing.T, expect Word, got Word) {
	if expect != got {
		t.Errorf("Expected variable to be: %x, got: %x", expect, got)
		t.Log(string(debug.Stack()))
	}
}
