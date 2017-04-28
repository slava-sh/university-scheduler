package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

const solveTimeLimit = 300 * time.Millisecond

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
			s := Solve(problem, solveTimeLimit)

			if s.Fatigue < 0 {
				t.Fatal("negative fatigue")
			}
			if s.Fatigue != s.computeFatigue() {
				t.Fatal("inconsistent fatigue")
			}

			for day := 1; day <= s.DaysPerWeek; day++ {
				for class := 1; class <= s.ClassesPerDay; class++ {
					numRooms := 0
					for group := 1; group <= s.NumGroups; group++ {
						prof := s.GroupSchedule[group][day][class]
						if prof != 0 {
							numRooms++
						}
					}
					if numRooms > s.NumRooms {
						t.Fatalf(
							"too many rooms occupied on (day=%d, class=%d)",
							day, class,
						)
					}
				}
			}

			for day := 1; day <= s.DaysPerWeek; day++ {
				for class := 1; class <= s.ClassesPerDay; class++ {
					for group := 1; group <= s.NumGroups; group++ {
						prof := s.GroupSchedule[group][day][class]
						if prof == 0 {
							continue
						}
						if s.ProfSchedule[prof][day][class] != group {
							t.Fatalf(
								"GroupSchedule is inconsistent with ProfSchedule "+
									"on (day=%d, class=%d)", day, class)
						}
					}
				}
			}

			type GP struct {
				Group int
				Prof  int
			}
			numClasses := make(map[GP]int)
			for day := 1; day <= s.DaysPerWeek; day++ {
				for class := 1; class <= s.ClassesPerDay; class++ {
					for group := 1; group <= s.NumGroups; group++ {
						prof := s.GroupSchedule[group][day][class]
						if prof == 0 {
							continue
						}
						numClasses[GP{group, prof}]++
					}
				}
			}
			for group := 1; group <= s.NumGroups; group++ {
				for prof := 1; prof <= s.NumProfs; prof++ {
					expected := s.NumClasses[group][prof]
					actual := numClasses[GP{group, prof}]
					if actual != expected {
						t.Fatalf(
							"expected %d classes for (group=%d, prof=%d), got %d",
							expected, group, prof, actual,
						)
					}
				}
			}
		})
	}
}

func BenchmarkSolve(b *testing.B) {
	in, err := os.Open("../input/04.txt")
	if err != nil {
		b.Fatal(err)
	}
	problem := ReadProblem(NewFastReader(in))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Solve(problem, solveTimeLimit)
	}
}
