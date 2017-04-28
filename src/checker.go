package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
)

func main() {
	in := NewFastReader(os.Stdin)
	out := log.New(os.Stdout, "", 0)

	p := ReadProblem(in)
	s := ReadSolution(in, p)

	fatigue := s.computeFatigue()
	if s.Fatigue != fatigue {
		out.Printf("expected fatigue %d, got %d", fatigue, s.Fatigue)
	}

	for day := 1; day <= DaysPerWeek; day++ {
		for class := 1; class <= ClassesPerDay; class++ {
			numRooms := 0
			for group := 1; group <= s.NumGroups; group++ {
				prof := s.GroupSchedule[group][day][class]
				if prof != 0 {
					numRooms++
				}
			}
			if numRooms > s.NumRooms {
				out.Printf(
					"too many rooms occupied on (day=%d, class=%d)",
					day, class,
				)
			}
		}
	}

	type GP struct {
		Group int
		Prof  int
	}
	numClasses := make(map[GP]int)
	for day := 1; day <= DaysPerWeek; day++ {
		for class := 1; class <= ClassesPerDay; class++ {
			for group := 1; group <= s.NumGroups; group++ {
				prof := s.GroupSchedule[group][day][class]
				if prof == 0 {
					continue
				}
				numClasses[GP{group, prof}]++
			}
		}
	}
	for group := 1; group <= s.NumGroups; group++ {
		for prof := 1; prof <= s.NumProfs; prof++ {
			expected := s.NumClasses[group][prof]
			actual := numClasses[GP{group, prof}]
			if actual != expected {
				out.Printf(
					"expected %d classes for (group=%d, prof=%d), got %d",
					expected, group, prof, actual,
				)
			}
		}
	}
}

func ReadSolution(in *FastReader, p *Problem) *Solution {
	s := new(Solution)
	s.Problem = p
	s.Fatigue = in.NextInt()
	for group := 1; group <= p.NumGroups; group++ {
		for class := 1; class <= ClassesPerDay; class++ {
			for day := 1; day <= DaysPerWeek; day++ {
				prof := in.NextInt()
				s.GroupSchedule[group][day][class] = prof
				s.ProfSchedule[prof][day][class] = group
			}
		}
	}
	return s
}

const (
	MaxGroup      = 60
	MaxProf       = 60
	DaysPerWeek   = 6
	ClassesPerDay = 7
)

type Problem struct {
	NumRooms   int
	NumGroups  int
	NumProfs   int
	NumClasses [MaxGroup + 1][MaxProf + 1]int
}

type Solution struct {
	*Problem
	Fatigue       int
	GroupSchedule [MaxGroup + 1][DaysPerWeek + 1][ClassesPerDay + 1]int
	ProfSchedule  [MaxProf + 1][DaysPerWeek + 1][ClassesPerDay + 1]int
}

func ReadProblem(in *FastReader) *Problem {
	var p Problem
	p.NumGroups = in.NextInt()
	p.NumProfs = in.NextInt()
	p.NumRooms = in.NextInt()
	for group := 1; group <= p.NumGroups; group++ {
		for prof := 1; prof <= p.NumProfs; prof++ {
			p.NumClasses[group][prof] = in.NextInt()
		}
	}
	return &p
}

func (s *Solution) computeFatigue() int {
	fatigue := 0
	for day := 1; day <= DaysPerWeek; day++ {
		for group := 1; group <= s.NumGroups; group++ {
			fatigue += s.groupFatigue(group, day)
		}
		for prof := 1; prof <= s.NumProfs; prof++ {
			fatigue += s.profFatigue(prof, day)
		}
	}
	return fatigue
}

func (s *Solution) groupFatigue(group, day int) int {
	maxClass := 0
	for class := ClassesPerDay; class > 0; class-- {
		if s.GroupSchedule[group][day][class] != 0 {
			maxClass = class
			break
		}
	}
	if maxClass == 0 {
		return 0
	}
	minClass := 0
	for class := 1; class <= ClassesPerDay; class++ {
		if s.GroupSchedule[group][day][class] != 0 {
			minClass = class
			break
		}
	}
	return square(2 + maxClass - minClass + 1)
}

func (s *Solution) profFatigue(prof, day int) int {
	maxClass := 0
	for class := ClassesPerDay; class > 0; class-- {
		if s.ProfSchedule[prof][day][class] != 0 {
			maxClass = class
			break
		}
	}
	if maxClass == 0 {
		return 0
	}
	minClass := 0
	for class := 1; class <= ClassesPerDay; class++ {
		if s.ProfSchedule[prof][day][class] != 0 {
			minClass = class
			break
		}
	}
	return square(2 + maxClass - minClass + 1)
}

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

type FastReader struct {
	r   *bufio.Reader
	buf string
}

func NewFastReader(r io.Reader) *FastReader {
	return &FastReader{
		r: bufio.NewReader(r),
	}
}

func (r *FastReader) advance() {
	for len(r.buf) != 0 && r.buf[0] == ' ' {
		r.buf = r.buf[1:]
	}
	var err error
	for len(r.buf) == 0 && err != io.EOF {
		buf := new(bytes.Buffer)
		for {
			var chunk []byte
			var more bool
			chunk, more, err = r.r.ReadLine()
			buf.Write(chunk)
			if !more {
				break
			}
		}
		r.buf = buf.String()
	}
}

func (r *FastReader) NextLine() string {
	r.advance()
	line := r.buf
	r.buf = ""
	return line
}

func (r *FastReader) NextWord() string {
	r.advance()
	var word string
	wordStart := 0
	for i := 0; i < len(r.buf); i++ {
		b := r.buf[i]
		if b == ' ' {
			if i != wordStart {
				word = r.buf[wordStart:i]
				r.buf = r.buf[i:]
				break
			}
			wordStart = i + 1
		}
	}
	if len(word) == 0 {
		word = r.buf
		r.buf = ""
	}
	return word
}

func (r *FastReader) NextInt() int {
	return parseInt(r.NextWord())
}

func parseInt(word string) int {
	sign := 1
	if word[0] == '-' {
		sign = -1
		word = word[1:]
	}
	result := 0
	for i := 0; i < len(word); i++ {
		result = result*10 + int(word[i]) - '0'
	}
	result *= sign
	return result
}
