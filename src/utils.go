package main

func square(x int) int {
	return x * x
}

func min(a, b int) int {
	if b < a {
		return b
	}
	return a
}

func max(a, b int) int {
	if b > a {
		return b
	}
	return a
}

func makeInts2(size1, size2 int) [][]int {
	result := make([][]int, size1)
	for i := 0; i < size1; i++ {
		result[i] = make([]int, size2)
	}
	return result
}

func makeInts3(size1, size2, size3 int) [][][]int {
	result := make([][][]int, size1)
	for i := 0; i < size1; i++ {
		result[i] = makeInts2(size2, size3)
	}
	return result
}

func copyInts(a []int) []int {
	copy := make([]int, len(a))
	for i := 0; i < len(a); i++ {
		copy[i] = a[i]
	}
	return copy
}

func copyInts2(a [][]int) [][]int {
	copy := make([][]int, len(a))
	for i := 0; i < len(a); i++ {
		copy[i] = copyInts(a[i])
	}
	return copy
}

func copyInts3(a [][][]int) [][][]int {
	copy := make([][][]int, len(a))
	for i := 0; i < len(a); i++ {
		copy[i] = copyInts2(a[i])
	}
	return copy
}
