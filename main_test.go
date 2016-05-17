package main

import "testing"

func TestSayHello(t *testing.T) {
	expected := "Hello filebroker"
	input := "filebroker"
	result := sayHello(input)
	if result != expected {
		t.Errorf("Got %s; want %s", result, expected)
	}

}
