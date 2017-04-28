package main

import (
	"fmt"
	"io"
	"math/rand"
)

type Solution struct {
	*Problem
	Fatigue       int
	GroupSchedule [MaxGroup + 1][DaysPerWeek + 1][ClassesPerDay + 1]int // [group][day][class] -> prof
}

type State struct {
	Solution
	ProfSchedule [MaxProf + 1][DaysPerWeek + 1][ClassesPerDay + 1]int // [prof][day][class] -> group
	NumFreeRooms [DaysPerWeek + 1][ClassesPerDay + 1]int              // [day][class] -> numFreeRooms
	GroupFatigue [MaxGroup + 1][DaysPerWeek + 1]int                   // [group][day] -> fatigue
	ProfFatigue  [MaxProf + 1][DaysPerWeek + 1]int                    // [prof][day] -> fatigue
	Edges        *EdgeSet
}

func (s *Solution) Print(out io.Writer) {
	fmt.Fprintf(out, "%d\n", s.Fatigue)
	for group := 1; group <= s.NumGroups; group++ {
		fmt.Fprintf(out, "\n")
		for class := 1; class <= ClassesPerDay; class++ {
			for day := 1; day <= DaysPerWeek; day++ {
				if day != 1 {
					fmt.Fprintf(out, " ")
				}
				fmt.Fprintf(out, "%d", s.GroupSchedule[group][day][class])
			}
			fmt.Fprintf(out, "\n")
		}
	}
}

func (s *State) groupFatigue(group, day int) int {
	maxClass := 0
	for class := ClassesPerDay; class > 0; class-- {
		if s.GroupSchedule[group][day][class] != 0 {
			maxClass = class
			break
		}
	}
	if maxClass == 0 {
		return 0
	}
	minClass := 0
	for class := 1; class <= ClassesPerDay; class++ {
		if s.GroupSchedule[group][day][class] != 0 {
			minClass = class
			break
		}
	}
	return square(2 + maxClass - minClass + 1)
}

func (s *State) profFatigue(prof, day int) int {
	maxClass := 0
	for class := ClassesPerDay; class > 0; class-- {
		if s.ProfSchedule[prof][day][class] != 0 {
			maxClass = class
			break
		}
	}
	if maxClass == 0 {
		return 0
	}
	minClass := 0
	for class := 1; class <= ClassesPerDay; class++ {
		if s.ProfSchedule[prof][day][class] != 0 {
			minClass = class
			break
		}
	}
	return square(2 + maxClass - minClass + 1)
}

func Solve(p *Problem, shouldWork func() bool) *Solution {
	s := solveNaive(p)
	bestSolution := s.Solution
	s.createEdges()
	for i := 0; shouldWork(); i++ {
		for try := 0; try < 10; try++ {
			// Generate a swap.
			g, p, d1, c1 := s.Edges.Pop()
			d2 := 1 + rand.Intn(DaysPerWeek)
			c2 := 1 + rand.Intn(ClassesPerDay)
			if s.NumFreeRooms[d2][c2] == 0 ||
				s.ProfSchedule[p][d2][c2] != 0 ||
				s.GroupSchedule[g][d2][c2] != 0 {
				s.Edges.Push(g, p, d1, c1)
				continue
			}

			prevFatigue := s.Fatigue
			prevGroupFatigue1 := s.GroupFatigue[g][d1]
			prevGroupFatigue2 := s.GroupFatigue[g][d2]
			prevProfFatigue1 := s.ProfFatigue[p][d1]
			prevProfFatigue2 := s.ProfFatigue[p][d2]

			// Apply swap.
			s.Fatigue -= s.GroupFatigue[g][d1]
			s.Fatigue -= s.ProfFatigue[p][d1]
			if d2 != d1 {
				s.Fatigue -= s.GroupFatigue[g][d2]
				s.Fatigue -= s.ProfFatigue[p][d2]
			}
			s.NumFreeRooms[d1][c1]++
			s.NumFreeRooms[d2][c2]--
			s.GroupSchedule[g][d1][c1] = 0
			s.GroupSchedule[g][d2][c2] = p
			s.ProfSchedule[p][d1][c1] = 0
			s.ProfSchedule[p][d2][c2] = g
			s.GroupFatigue[g][d1] = s.groupFatigue(g, d1)
			s.ProfFatigue[p][d1] = s.profFatigue(p, d1)
			s.Fatigue += s.GroupFatigue[g][d1]
			s.Fatigue += s.ProfFatigue[p][d1]
			if d2 != d1 {
				s.GroupFatigue[g][d2] = s.groupFatigue(g, d2)
				s.ProfFatigue[p][d2] = s.profFatigue(p, d2)
				s.Fatigue += s.GroupFatigue[g][d2]
				s.Fatigue += s.ProfFatigue[p][d2]
			}

			if s.Fatigue <= prevFatigue {
				// Accept swap.
				if s.Fatigue < bestSolution.Fatigue {
					bestSolution = s.Solution
				}

				// Update edges.
				for _, delta := range []int{-1, 1} {
					newClass := c1 + delta
					if !(1 <= newClass && newClass <= ClassesPerDay) {
						continue
					}
					if s.isEdgeForGroup(g, d1, newClass) {
						newProf := s.GroupSchedule[g][d1][newClass]
						if newProf != 0 && s.isEdgeForProf(newProf, d1, newClass) {
							s.Edges.Push(g, newProf, d1, newClass)
						}
					}
					if s.isEdgeForProf(p, d1, newClass) {
						newGroup := s.ProfSchedule[p][d1][newClass]
						if newGroup != 0 && s.isEdgeForGroup(newGroup, d1, newClass) {
							s.Edges.Push(newGroup, p, d1, newClass)
						}
					}
				}
				if s.isEdgeForGroup(g, d2, c2) {
					if s.isEdgeForProf(p, d2, c2) {
						s.Edges.Push(g, p, d2, c2)
					}
					for _, delta := range []int{-1, 1} {
						newClass := c2 + delta
						if 1 <= newClass && newClass <= ClassesPerDay &&
							!s.isEdgeForGroup(g, d2, newClass) {
							newProf := s.GroupSchedule[g][d2][newClass]
							if newProf != 0 {
								s.Edges.Remove(g, newProf, d2, newClass)
							}
						}
					}
				}
				if s.isEdgeForProf(p, d2, c2) {
					for _, delta := range []int{-1, 1} {
						newClass := c2 + delta
						if 1 <= newClass && newClass <= ClassesPerDay &&
							!s.isEdgeForProf(p, d2, newClass) {
							newGroup := s.ProfSchedule[p][d2][newClass]
							if newGroup != 0 {
								s.Edges.Remove(p, newGroup, d2, newClass)
							}
						}
					}
				}
			} else {
				// Discard swap.
				s.NumFreeRooms[d1][c1]--
				s.NumFreeRooms[d2][c2]++
				s.GroupSchedule[g][d2][c2] = 0
				s.GroupSchedule[g][d1][c1] = p
				s.ProfSchedule[p][d2][c2] = 0
				s.ProfSchedule[p][d1][c1] = g
				s.Fatigue = prevFatigue
				s.GroupFatigue[g][d1] = prevGroupFatigue1
				s.ProfFatigue[p][d1] = prevProfFatigue1
				if d2 != d1 {
					s.GroupFatigue[g][d2] = prevGroupFatigue2
					s.ProfFatigue[p][d2] = prevProfFatigue2
				}
				s.Edges.Push(g, p, d1, c1)
			}

			break
		}
	}
	return &bestSolution
}

func solveNaive(p *Problem) *State {
	var s State
	s.Problem = p

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
	for day := 1; day <= DaysPerWeek; day++ {
		for class := 1; class <= ClassesPerDay; class++ {
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

	for day := 1; day <= DaysPerWeek; day++ {
		for group := 1; group <= s.NumGroups; group++ {
			s.GroupFatigue[group][day] = s.groupFatigue(group, day)
			s.Fatigue += s.GroupFatigue[group][day]
		}
		for prof := 1; prof <= s.NumProfs; prof++ {
			s.ProfFatigue[prof][day] = s.profFatigue(prof, day)
			s.Fatigue += s.ProfFatigue[prof][day]
		}
	}
	return &s
}

func (s *State) createEdges() {
	s.Edges = NewEdgeSet()
	for group := 1; group <= s.NumGroups; group++ {
		for day := 1; day <= DaysPerWeek; day++ {
			for class := 1; class <= ClassesPerDay; class++ {
				prof := s.GroupSchedule[group][day][class]
				if prof == 0 {
					continue
				}
				if s.isEdgeForGroup(group, day, class) &&
					s.isEdgeForProf(prof, day, class) {
					s.Edges.Push(group, prof, day, class)
				}
			}
		}
	}
}

func (s *State) isEdgeForGroup(group, day, class int) bool {
	if s.GroupSchedule[group][day][class] == 0 {
		return false
	}
	if class > 1 && s.GroupSchedule[group][day][class-1] == 0 {
		return true
	}
	if class < ClassesPerDay && s.GroupSchedule[group][day][class+1] == 0 {
		return true
	}
	return false
}

func (s *State) isEdgeForProf(prof, day, class int) bool {
	if s.ProfSchedule[prof][day][class] == 0 {
		return false
	}
	if class > 1 && s.ProfSchedule[prof][day][class-1] == 0 {
		return true
	}
	if class < ClassesPerDay && s.ProfSchedule[prof][day][class+1] == 0 {
		return true
	}
	return false
}
