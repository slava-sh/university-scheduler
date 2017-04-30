package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
)

const (
	MaxClassesPerGroup = 24
	MaxClassesPerProf  = 24
	MaxRoomUtilization = 0.75
)

const (
	MaxRoom       = 60
	MaxGroup      = 60
	MaxProf       = 60
	DaysPerWeek   = 6
	ClassesPerDay = 7
)

func main() {
	log.SetPrefix("problem_generator: ")
	log.SetFlags(0)

	if len(os.Args) != 2 {
		usage()
		os.Exit(2)
	}

	seed, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	rand.Seed(int64(seed))
	p := GenerateProblem()
	p.Print(os.Stdout)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: ./problem_generator <seed>\n")
}

type Problem struct {
	NumRooms   int
	NumGroups  int
	NumProfs   int
	NumClasses [MaxGroup + 1][MaxProf + 1]int
}

func GenerateProblem() Problem {
	var p Problem
	for {
		p.NumGroups = 1 + rand.Intn(MaxGroup)
		p.NumProfs = 1 + rand.Intn(MaxProf)
		p.NumRooms = 1 + rand.Intn(MaxRoom)
		groupClasses := make([]int, p.NumGroups+1)
		profClasses := make([]int, p.NumProfs+1)
		totalClasses := 0
		for group := 1; group <= p.NumGroups; group++ {
			for prof := 1; prof <= p.NumProfs; prof++ {
				maxGroupClasses := MaxClassesPerGroup - groupClasses[group]
				maxProfClasses := MaxClassesPerProf - profClasses[prof]
				maxClasses := min(maxGroupClasses, maxProfClasses)
				numClasses := rand.Intn(maxClasses + 1)
				p.NumClasses[group][prof] = numClasses
				groupClasses[group] += numClasses
				profClasses[prof] += numClasses
				totalClasses += numClasses
			}
		}
		capacity := p.NumRooms * ClassesPerDay * DaysPerWeek
		roomUtilization := float64(totalClasses) / float64(capacity)
		if roomUtilization <= MaxRoomUtilization {
			break
		}
	}
	return p
}

func (p Problem) Print(out io.Writer) {
	fmt.Fprintln(out, p.NumGroups, p.NumProfs, p.NumRooms)
	for group := 1; group <= p.NumGroups; group++ {
		for prof := 1; prof <= p.NumProfs; prof++ {
			if prof != 1 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, p.NumClasses[group][prof])
		}
		fmt.Fprintln(out)
	}
}

func min(a, b int) int {
	if b < a {
		return b
	}
	return a
}
