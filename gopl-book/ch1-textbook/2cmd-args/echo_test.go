package main

import (
	"testing"
)

var args = []string{
	"Alpha",
	"Bravo",
	"Charlie",
	"Delta",
	"Echo",
	"Foxtrot",
	"Golf",
	"Hotel",
	"India",
	"Juliett",
	"Kilo",
	"Lima",
	"Mike",
	"November",
	"Oscar",
	"Papa",
	"Quebec",
	"Romeo",
	"Sierra",
	"Tango",
	"Uniform",
	"Victor",
	"Whiskey",
	"X-ray",
	"Yankee",
	"Zulu",
}

func BenchmarkEcho1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		echo1(args)
	}
}
func BenchmarkEcho3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		echo3(args)
	}
}
