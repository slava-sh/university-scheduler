package main

import (
	"bufio"
	"encoding/csv"
	"log"
	"os"
	"time"
)

const timeLimit = 10*time.Second - 100*time.Millisecond

var sa *csv.Writer

func main() {
	file, err := os.Create("sa.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	sa = csv.NewWriter(file)
	defer sa.Flush()

	start := time.Now()
	in := NewFastReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	problem := ReadProblem(in)
	solution := Solve(problem, timeLimit)
	solution.Print(out)
	log.Println("elapsed:", time.Since(start))
}
