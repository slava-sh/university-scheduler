package main

import (
	"bufio"
	"os"
)

func main() {
	in := NewFastReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	problem := ReadProblem(in)
	solution := Solve(problem)
	solution.Print(out)
}
