package main

import (
	"bufio"
	"os"
	"time"
)

const timeLimit = 10*time.Second - 60*time.Millisecond

func main() {
	start := time.Now()
	in := NewFastReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	problem := ReadProblem(in)
	solution := Solve(problem, func() bool {
		return time.Since(start) <= timeLimit
	})
	solution.Print(out)
}
