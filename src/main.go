package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"
)

const (
	timeLimit    = 10*time.Second - 60*time.Millisecond
	local        = true
	runsPerInput = 5
)

func main() {
	if local {
		runAll()
	} else {
		run(os.Stdin, os.Stdout)
	}
}

func run(stdin io.Reader, stdout io.Writer) int {
	in := NewFastReader(stdin)
	out := bufio.NewWriter(stdout)
	defer out.Flush()
	problem := ReadProblem(in)
	solution := Solve(problem, timeLimit)
	solution.Print(out)
	return solution.Fatigue
}

func runAll() {
	totalScore := 0.0
	for input := 1; input <= 10; input++ {
		filename := fmt.Sprintf("input/%02d.txt", input)
		fmt.Printf("%s ", filename)
		scoreSum := 0.0
		for i := 0; i < runsPerInput; i++ {
			stdin, err := os.Open(filename)
			if err != nil {
				panic(err)
			}
			score := run(stdin, ioutil.Discard)
			scoreSum += float64(score)
			fmt.Printf("%5d ", score)
		}
		averageScore := scoreSum / runsPerInput
		totalScore += averageScore
		fmt.Printf("%9.3f\n", averageScore)
	}
	fmt.Printf("score %9.3f\n", totalScore)
}
