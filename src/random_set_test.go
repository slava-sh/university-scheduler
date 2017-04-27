package main

import (
	"math/rand"
	"testing"
)

func TestRandomSet(t *testing.T) {
	const actions = "UUUOUOUUOUOOUOOO"
	rand.Seed(0)
	rs := MakeRandomSet()
	solutions := make(map[*Solution]bool)
	for _, action := range actions {
		if action == 'U' {
			s := new(Solution)
			solutions[s] = true
			rs.Push(s)
		} else {
			s := rs.Pop()
			if !solutions[s] {
				t.Fatal("unknown s")
			}
			delete(solutions, s)
		}
		length := rs.Len()
		if length != len(solutions) {
			t.Fatalf("wrong length: expected %d, got %d", len(solutions), length)
		}
	}
}
