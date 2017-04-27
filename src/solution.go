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

func Solve(p Problem, timeLimit time.Duration) Solution {
	best := solveNaive(p)
	firstFatigue := best.Fatigue
	start := time.Now()
	for i := 0; ; i++ {
		if i != 0 {
			elapsed := time.Since(start)
			timePerStep := time.Duration(int(elapsed) / i)
			if elapsed+timePerStep >= timeLimit {
				log.Println("steps:", i)
				log.Println("time per step:", timePerStep)
				break
			}
		}
		other := neighbor(best)
		if other.Fatigue > best.Fatigue {
			best = other
		}
	}
	log.Println("fatigue:", firstFatigue, "->", best.Fatigue)
	return best
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
		if g != 0 &&
			s.ProfSchedule[p][d2][c2] == 0 &&
			s.GroupSchedule[g][d2][c2] == 0 &&
			s.NumFreeRooms[d2][c2] != 0 {
			s.NumFreeRooms[d1][c1]++
			s.NumFreeRooms[d2][c2]--
			s.GroupSchedule[g][d1][c1] = 0
			s.GroupSchedule[g][d2][c2] = p
			s.ProfSchedule[p][d1][c1] = 0
			s.ProfSchedule[p][d2][c2] = g
			s.UpdateFatigue()
			break
		}
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

func makeInts2(size1, size2 int) [][]int {
	result := make([][]int, size1)
	for i := 0; i < size1; i++ {
		result[i] = make([]int, size2)
	}
	return result
}

func makeInts3(size1, size2, size3 int) [][][]int {
	result := make([][][]int, size1)
	for i := 0; i < size1; i++ {
		result[i] = makeInts2(size2, size3)
	}
	return result
}

func copyInts(a []int) []int {
	copy := make([]int, len(a))
	for i := 0; i < len(a); i++ {
		copy[i] = a[i]
	}
	return copy
}

func copyInts2(a [][]int) [][]int {
	copy := make([][]int, len(a))
	for i := 0; i < len(a); i++ {
		copy[i] = copyInts(a[i])
	}
	return copy
}

func copyInts3(a [][][]int) [][][]int {
	copy := make([][][]int, len(a))
	for i := 0; i < len(a); i++ {
		copy[i] = copyInts2(a[i])
	}
	return copy
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
