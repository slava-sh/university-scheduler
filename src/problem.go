package main

type Problem struct {
	DaysPerWeek   int
	ClassesPerDay int
	NumRooms      int
	NumGroups     int
	NumProfs      int
	NumClasses    [][]int // [group][prof] -> numClasses
}

func ReadProblem(in *FastReader) Problem {
	var p Problem
	p.DaysPerWeek = 6
	p.ClassesPerDay = 7
	p.NumGroups = in.NextInt()
	p.NumProfs = in.NextInt()
	p.NumRooms = in.NextInt()
	p.NumClasses = make([][]int, p.NumGroups+1)
	for group := 1; group <= p.NumGroups; group++ {
		p.NumClasses[group] = make([]int, p.NumProfs+1)
		for prof := 1; prof <= p.NumProfs; prof++ {
			p.NumClasses[group][prof] = in.NextInt()
		}
	}
	return p
}
