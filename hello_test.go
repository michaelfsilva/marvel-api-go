package main

import (
	"testing"
)

// test hello function
func TestHelloWithoutParam(t *testing.T) {
	// test for empty argument
	emptyResult := hello("") // should return "Hello Dude!"

	if emptyResult != "Hello Dude!" {
		t.Errorf("hello(\"\") failed, expected %v, got %v", "Hello Dude!", emptyResult)
	} else {
		t.Logf("hello(\"\") success, expected %v, got %v", "Hello Dude!", emptyResult)
	}
}

func TestHelloWithParam(t *testing.T) {
	// test for valid argument
	result := hello("Mike") // should return "Hello Mike!"

	if result != "Hello Mike!" {
		t.Errorf("hello(\"Mike\") failed, expected %v, got %v", "Hello Mike!", result)
	} else {
		t.Logf("hello(\"Mike\") success, expected %v, got %v", "Hello Mike!", result)
	}
}
