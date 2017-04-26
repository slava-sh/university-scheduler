package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	in  = NewFastReader(os.Stdin)
	out = bufio.NewWriter(os.Stdout)
)

func main() {
	defer out.Flush()
	n := in.NextInt()
	m := in.NextInt()
	a := in.NextInt()
	printf("%d %d %d\n", n, m, a)
}

func printf(format string, a ...interface{}) {
	fmt.Fprintf(out, format, a...)
}
