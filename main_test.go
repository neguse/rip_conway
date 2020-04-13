package main

import "testing"

func TestGameNext(t *testing.T) {
	// Test by blinker pattern.
	// https://www.conwaylife.com/wiki/Blinker
	g := NewGame(3, 3)
	*g.At(0, 0, 1) = true
	*g.At(0, 1, 1) = true
	*g.At(0, 2, 1) = true

	g.Print()

	g.UpdateCells()

	fn := func(x, y int, expect bool) {
		actual := *g.At(0, x, y)
		if actual != expect {
			t.Errorf("x:%d y:%d expect:%v actual:%v", x, y, expect, actual)
		}
	}
	fn(1, 0, true)
	fn(1, 1, true)
	fn(1, 2, true)
	fn(0, 0, false)
	fn(0, 1, false)
	fn(0, 2, false)
	fn(2, 0, false)
	fn(2, 1, false)
	fn(2, 2, false)
}
