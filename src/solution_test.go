package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSolve_sanity(t *testing.T) {
	filenames, err := filepath.Glob("../input/*.txt")
	if err != nil {
		t.Fatal(err)
	}
	for _, filename := range filenames {
		t.Run(filename[3:], func(t *testing.T) {
			in, err := os.Open(filename)
			if err != nil {
				t.Fatal(err)
			}
			problem := ReadProblem(NewFastReader(in))
			s := Solve(problem)
			if s.Fatigue < 0 {
				t.Fatal("negative fatigue")
			}
			for day := 1; day <= s.DaysPerWeek; day++ {
				for class := 1; class <= s.ClassesPerDay; class++ {
					profIsBusy := make(map[int]bool)
					numRooms := 0
					for group := 1; group <= s.NumGroups; group++ {
						prof := s.Schedule[group][day][class]
						if prof == 0 {
							continue
						}
						if profIsBusy[prof] {
							t.Fatalf(
								"prof=%d has multiple groups on (day=%d, class=%d)",
								prof, day, class,
							)
						}
						profIsBusy[prof] = true
						numRooms++
					}
					if numRooms > s.NumRooms {
						t.Fatalf(
							"too many rooms occupied on (day=%d, class=%d)",
							day, class,
						)
					}
				}
			}
		})
	}
}
