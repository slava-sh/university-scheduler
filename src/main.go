package main

import (
	"bufio"
	"os"
	"time"
)

const timeLimit = 10*time.Second - 60*time.Millisecond

func main() {
	in := NewFastReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	problem := ReadProblem(in)
	solution := Solve(problem, timeLimit)
	solution.Print(out)
}
