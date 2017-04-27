package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"time"
)

type Solution struct {
	Problem
	Fatigue       int
	GroupSchedule [][][]int // [group][day][class] -> prof
	ProfSchedule  [][][]int // [prof][day][class] -> group
	NumFreeRooms  [][]int   // [day][class] -> numFreeRooms
}

func (s *Solution) Copy() Solution {
	var copy Solution
	copy = *s
	copy.GroupSchedule = copyInts3(s.GroupSchedule)
	copy.ProfSchedule = copyInts3(s.ProfSchedule)
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
			for class := s.ClassesPerDay; class > 0; class-- {
				if s.GroupSchedule[group][day][class] != 0 {
					maxClass = class
					break
				}
			}
			if maxClass == 0 {
				continue
			}
			minClass := 0
			for class := 1; class <= s.ClassesPerDay; class++ {
				if s.GroupSchedule[group][day][class] != 0 {
					minClass = class
					break
				}
			}
			fatigue += square(2 + maxClass - minClass + 1)
		}
		for prof := 1; prof <= s.NumProfs; prof++ {
			maxClass := 0
			for class := s.ClassesPerDay; class > 0; class-- {
				if s.ProfSchedule[prof][day][class] != 0 {
					maxClass = class
					break
				}
			}
			if maxClass == 0 {
				continue
			}
			minClass := 0
			for class := 1; class <= s.ClassesPerDay; class++ {
				if s.ProfSchedule[prof][day][class] != 0 {
					minClass = class
					break
				}
			}
			fatigue += square(2 + maxClass - minClass + 1)
		}
	}
	s.Fatigue = fatigue
}

func Solve(p Problem, timeLimit time.Duration) Solution {
	start := time.Now()
	solution := solveNaive(p)
	bestSolution := solution
	firstFatigue := solution.Fatigue
	stepStart := time.Now()
	for i := 0; ; i++ {
		if i != 0 {
			timePerStep := time.Duration(int(time.Since(stepStart)) / i)
			timeLeft := timeLimit - time.Since(start)
			if timeLeft <= timePerStep {
				log.Println("steps:", i)
				log.Println("time per step:", timePerStep)
				log.Println("fatigue:", firstFatigue, "->", bestSolution.Fatigue)
				break
			}
		}
		newSolution := neighbor(solution)
		delta := newSolution.Fatigue - solution.Fatigue
		if shouldAccept(delta) {
			solution = newSolution
			if solution.Fatigue < bestSolution.Fatigue {
				bestSolution = solution
			}
		}
	}
	return bestSolution
}

func shouldAccept(delta int) bool {
	return delta <= 0
}

func neighbor(s Solution) Solution {
	s = s.Copy()
	for try := 0; try < 100; try++ {
		d1 := 1 + rand.Intn(s.DaysPerWeek)
		d2 := 1 + rand.Intn(s.DaysPerWeek)
		c1 := 1 + rand.Intn(s.ClassesPerDay)
		c2 := 1 + rand.Intn(s.ClassesPerDay)
		p := 1 + rand.Intn(s.NumProfs)
		g := s.ProfSchedule[p][d1][c1]
		if g == 0 ||
			s.NumFreeRooms[d2][c2] == 0 ||
			s.ProfSchedule[p][d2][c2] != 0 ||
			s.GroupSchedule[g][d2][c2] != 0 {
			continue
		}
		s.NumFreeRooms[d1][c1]++
		s.NumFreeRooms[d2][c2]--
		s.GroupSchedule[g][d1][c1] = 0
		s.GroupSchedule[g][d2][c2] = p
		s.ProfSchedule[p][d1][c1] = 0
		s.ProfSchedule[p][d2][c2] = g
		s.UpdateFatigue()
		break
	}
	return s
}

func solveNaive(p Problem) Solution {
	var s Solution
	s.Problem = p
	s.GroupSchedule = makeInts3(p.NumGroups+1, p.DaysPerWeek+1, p.ClassesPerDay+1)
	s.ProfSchedule = makeInts3(p.NumProfs+1, p.DaysPerWeek+1, p.ClassesPerDay+1)
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
