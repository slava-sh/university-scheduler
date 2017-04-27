package main

import (
	"bufio"
	"log"
	"os"
	"time"
)

const timeLimit = 10*time.Second - 100*time.Millisecond

func main() {
	start := time.Now()
	in := NewFastReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	problem := ReadProblem(in)
	solution := Solve(problem, timeLimit)
	solution.Print(out)
	log.Println("elapsed:", time.Since(start))
}
