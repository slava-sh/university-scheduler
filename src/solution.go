package main

import (
	"fmt"
	"io"
	"math/rand"
	"time"
)

type Solution struct {
	*Problem
	Fatigue       int
	GroupSchedule Treap   // (group, day, class) -> prof
	ProfSchedule  Treap   // (prof, day, class) -> group
	NumFreeRooms  [][]int // [day][class] -> numFreeRooms
}

func (s *Solution) Copy() *Solution {
	copy := new(Solution)
	*copy = *s
	copy.NumFreeRooms = copyInts2(s.NumFreeRooms)
	return copy
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
				fmt.Fprintf(out, "%d", s.GroupSchedule.Get(group, day, class))
			}
			fmt.Fprintf(out, "\n")
		}
	}
}

func (s *Solution) computeFatigue() int {
	fatigue := 0
	for day := 1; day <= s.DaysPerWeek; day++ {
		for group := 1; group <= s.NumGroups; group++ {
			fatigue += s.groupFatigue(group, day)
		}
		for prof := 1; prof <= s.NumProfs; prof++ {
			fatigue += s.profFatigue(prof, day)
		}
	}
	return fatigue
}

func (s *Solution) groupFatigue(group, day int) int {
	maxClass := 0
	for class := s.ClassesPerDay; class > 0; class-- {
		if s.GroupSchedule.Get(group, day, class) != 0 {
			maxClass = class
			break
		}
	}
	if maxClass == 0 {
		return 0
	}
	minClass := 0
	for class := 1; class <= s.ClassesPerDay; class++ {
		if s.GroupSchedule.Get(group, day, class) != 0 {
			minClass = class
			break
		}
	}
	return square(2 + maxClass - minClass + 1)
}

func (s *Solution) profFatigue(prof, day int) int {
	maxClass := 0
	for class := s.ClassesPerDay; class > 0; class-- {
		if s.ProfSchedule.Get(prof, day, class) != 0 {
			maxClass = class
			break
		}
	}
	if maxClass == 0 {
		return 0
	}
	minClass := 0
	for class := 1; class <= s.ClassesPerDay; class++ {
		if s.ProfSchedule.Get(prof, day, class) != 0 {
			minClass = class
			break
		}
	}
	return square(2 + maxClass - minClass + 1)
}

const PopulationSize = 5
const NumPops = 2

func Solve(p Problem, timeLimit time.Duration) *Solution {
	start := time.Now()
	firstSolution := solveNaive(p)
	bestSolution := firstSolution
	population := MakeRandomSet()
	for i := 0; i < PopulationSize; i++ {
		population.Push(firstSolution)
	}
	loopStart := time.Now()
	for i := 0; ; i++ {
		if i != 0 {
			timePerStep := time.Duration(int(time.Since(loopStart)) / i)
			timeLeft := timeLimit - time.Since(start)
			if timeLeft <= timePerStep {
				break
			}
		}

		solution := population.Pop()
		for j := 1; j < NumPops; j++ {
			other := population.Pop()
			if other.Fatigue < solution.Fatigue {
				solution = other
			}
		}

		population.Push(solution)
		for j := 1; j < NumPops; j++ {
			solution = randomNeighbor(solution)
			if solution.Fatigue < bestSolution.Fatigue {
				bestSolution = solution
			}
			population.Push(solution)
		}
	}
	return bestSolution
}

func randomNeighbor(s *Solution) *Solution {
	s = s.Copy()
	for try := 0; try < 100; try++ {
		d1 := 1 + rand.Intn(s.DaysPerWeek)
		d2 := 1 + rand.Intn(s.DaysPerWeek)
		c1 := 1 + rand.Intn(s.ClassesPerDay)
		c2 := 1 + rand.Intn(s.ClassesPerDay)
		p := 1 + rand.Intn(s.NumProfs)
		g := s.ProfSchedule.Get(p, d1, c1)
		if g == 0 ||
			s.NumFreeRooms[d2][c2] == 0 ||
			s.ProfSchedule.Get(p, d2, c2) != 0 ||
			s.GroupSchedule.Get(g, d2, c2) != 0 {
			continue
		}
		if 0 < c1 && c1 < s.ClassesPerDay {
			groupWillHaveEmptySlot :=
				s.GroupSchedule.Get(g, d1, c1-1) != 0 &&
					s.GroupSchedule.Get(g, d1, c1+1) != 0
			if groupWillHaveEmptySlot {
				continue
			}
			profWillHaveEmptySlot :=
				s.ProfSchedule.Get(p, d1, c1-1) != 0 &&
					s.ProfSchedule.Get(p, d1, c1+1) != 0
			if profWillHaveEmptySlot {
				continue
			}
		}
		s.Fatigue -= s.groupFatigue(g, d1)
		s.Fatigue -= s.profFatigue(p, d1)
		if d2 != d1 {
			s.Fatigue -= s.groupFatigue(g, d2)
			s.Fatigue -= s.profFatigue(p, d2)
		}
		s.NumFreeRooms[d1][c1]++
		s.NumFreeRooms[d2][c2]--
		s.GroupSchedule = s.GroupSchedule.Remove(g, d1, c1).Set(g, d2, c2, p)
		s.ProfSchedule = s.ProfSchedule.Remove(p, d1, c1).Set(p, d2, c2, g)
		s.Fatigue += s.groupFatigue(g, d1)
		s.Fatigue += s.profFatigue(p, d1)
		if d2 != d1 {
			s.Fatigue += s.groupFatigue(g, d2)
			s.Fatigue += s.profFatigue(p, d2)
		}
		break
	}
	return s
}

func solveNaive(p Problem) *Solution {
	var s Solution
	s.Problem = &p
	s.NumFreeRooms = makeInts2(p.DaysPerWeek+1, p.ClassesPerDay+1)

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
			s.NumFreeRooms[day][class] = s.NumRooms
			groupIsBusy := make(map[int]bool)
			profIsBusy := make(map[int]bool)
			for groupAndProf := range classesToSchedule {
				group := groupAndProf.Group
				prof := groupAndProf.Prof
				if profIsBusy[prof] || groupIsBusy[group] {
					continue
				}
				if s.NumFreeRooms[day][class] == 0 {
					break
				}
				s.NumFreeRooms[day][class]--
				classesToSchedule[groupAndProf]--
				if classesToSchedule[groupAndProf] == 0 {
					delete(classesToSchedule, groupAndProf)
				}
				s.GroupSchedule = s.GroupSchedule.Set(group, day, class, prof)
				s.ProfSchedule = s.ProfSchedule.Set(prof, day, class, group)
				groupIsBusy[group] = true
				profIsBusy[prof] = true
			}
		}
	}

	s.Fatigue = s.computeFatigue()
	return &s
}
