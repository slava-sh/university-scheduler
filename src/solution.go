package main

import (
	"fmt"
	"io"
)

type Solution struct {
	Problem
	Fatigue       int
	GroupSchedule [][][]int // [group][day][class] -> prof
	ProfSchedule  [][][]int // [group][day][class] -> group
}

func (s *Solution) Print(out io.Writer) {
	fmt.Fprintf(out, "%d\n", s.Fatigue)
	for group := 1; group <= s.NumGroups; group++ {
		fmt.Fprintf(out, "\n")
		for class := 1; class <= s.ClassesPerDay; class++ {
			for day := 1; day <= s.DaysPerWeek; day++ {
				if day != 1 {
					fmt.Fprintf(out, " ")
				}
				fmt.Fprintf(out, "%d", s.GroupSchedule[group][day][class])
			}
			fmt.Fprintf(out, "\n")
		}
	}
}

func (s *Solution) UpdateFatigue() {
	fatigue := 0
	for day := 1; day <= s.DaysPerWeek; day++ {
		for group := 1; group <= s.NumGroups; group++ {
			maxClass := 0
			minClass := s.ClassesPerDay
			for class := 1; class <= s.ClassesPerDay; class++ {
				if s.GroupSchedule[group][day][class] == 0 {
					continue
				}
				minClass = min(minClass, class)
				maxClass = max(maxClass, class)
			}
			if maxClass == 0 {
				continue
			}
			fatigue += square(2 + maxClass - minClass + 1)
		}
		for prof := 1; prof <= s.NumProfs; prof++ {
			maxClass := 0
			minClass := s.ClassesPerDay
			for class := 1; class <= s.ClassesPerDay; class++ {
				if s.ProfSchedule[prof][day][class] == 0 {
					continue
				}
				minClass = min(minClass, class)
				maxClass = max(maxClass, class)
			}
			if maxClass == 0 {
				continue
			}
			fatigue += square(2 + maxClass - minClass + 1)
		}
	}
	s.Fatigue = fatigue
}

func Solve(p Problem) Solution {
	return solveNaive(p)
}

func solveNaive(p Problem) Solution {
	var s Solution
	s.Problem = p
	s.GroupSchedule = makeSchedule(p, p.NumGroups)
	s.ProfSchedule = makeSchedule(p, p.NumProfs)

	type GroupAndProf struct {
		Group int
		Prof  int
	}
	classesToSchedule := make(map[GroupAndProf]int)
	for group := 1; group <= s.NumGroups; group++ {
		for prof := 1; prof <= s.NumProfs; prof++ {
			if s.NumClasses[group][prof] == 0 {
				continue
			}
			groupAndProf := GroupAndProf{group, prof}
			classesToSchedule[groupAndProf] = s.NumClasses[group][prof]
		}
	}
	for day := 1; day <= s.DaysPerWeek; day++ {
		for class := 1; class <= s.ClassesPerDay; class++ {
			room := 0
			groupIsBusy := make(map[int]bool)
			profIsBusy := make(map[int]bool)
			for groupAndProf := range classesToSchedule {
				group := groupAndProf.Group
				prof := groupAndProf.Prof
				if profIsBusy[prof] || groupIsBusy[group] {
					continue
				}
				if room == s.NumRooms {
					break
				}
				room++
				classesToSchedule[groupAndProf]--
				if classesToSchedule[groupAndProf] == 0 {
					delete(classesToSchedule, groupAndProf)
				}
				s.GroupSchedule[group][day][class] = prof
				s.ProfSchedule[prof][day][class] = group
				groupIsBusy[group] = true
				profIsBusy[prof] = true
			}
		}
	}

	s.UpdateFatigue()
	return s
}

func makeSchedule(p Problem, size int) [][][]int {
	schedule := make([][][]int, size+1)
	for i := 1; i <= size; i++ {
		schedule[i] = make([][]int, p.DaysPerWeek+1)
		for day := 1; day <= p.DaysPerWeek; day++ {
			schedule[i][day] = make([]int, p.ClassesPerDay+1)
		}
	}
	return schedule
}

func square(x int) int {
	return x * x
}

func min(a, b int) int {
	if b < a {
		return b
	}
	return a
}

func max(a, b int) int {
	if b > a {
		return b
	}
	return a
}
