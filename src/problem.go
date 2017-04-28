package main

const (
	MaxRoom       = 60
	MaxGroup      = 60
	MaxProf       = 60
	DaysPerWeek   = 6
	ClassesPerDay = 7
)

type Problem struct {
	NumRooms   int
	NumGroups  int
	NumProfs   int
	NumClasses [MaxGroup + 1][MaxProf + 1]int // [group][prof] -> numClasses
}

func ReadProblem(in *FastReader) *Problem {
	var p Problem
	p.NumGroups = in.NextInt()
	p.NumProfs = in.NextInt()
	p.NumRooms = in.NextInt()
	for group := 1; group <= p.NumGroups; group++ {
		for prof := 1; prof <= p.NumProfs; prof++ {
			p.NumClasses[group][prof] = in.NextInt()
		}
	}
	return &p
}
