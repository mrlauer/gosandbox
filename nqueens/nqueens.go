/*
nqueens computes the number of ways to place N mutually non-attacking nqueens on an
N x N chessboard.

	nqueens 10

takes N = 10, etc. 8 is the default.

	nqueens -p m N

restricts the computation to m processes.
*/
package main

import (
	"fmt"
	"flag"
	"strconv"
	"runtime"
	"time"
)

func nqueens(n int) int {
	pos := make([]int, n)
	var helper func(i int, pos []int) int
	helper = func(i int, pos []int) int {
		if i >= n {
			return 1
		}
		num := 0
		for j := 0; j < n; j++ {
			ok := true
			for k := 0; k < i; k++ {
				posk := pos[k]
				if j == posk || i+j == k+posk || i-j == k-posk {
					ok = false
					break
				}
			}
			if ok {
				pos[i] = j
				num += helper(i+1, pos)
			}
		}
		return num
	}
	return helper(0, pos)
}

func helper(i int, n int, pos []int) int {
	if i >= n {
		return 1
	}
	num := 0
	for j := 0; j < n; j++ {
		ok := true
		for k := 0; k < i; k++ {
			posk := pos[k]
			if j == posk || i+j == k+posk || i-j == k-posk {
				ok = false
				break
			}
		}
		if ok {
			pos[i] = j
			num += helper(i+1, n, pos)
		}
	}
	return num
}

func nqueens0(n int) int {
	pos := make([]int, n)
	return helper(0, n, pos)
}

func nqueens2(n int) int {
	var helper func(i int, pos []int) int
	helper = func(i int, pos []int) int {
		if i >= n {
			return 1
		}
		num := 0
		for j := 0; j < n; j++ {
			ok := true
			for k := 0; k < i; k++ {
				posk := pos[k]
				if j == posk || i+j == k+posk || i-j == k-posk {
					ok = false
					break
				}
			}
			if ok {
				pos[i] = j
				num += helper(i+1, pos)
			}
		}
		return num
	}
	num := 0
	ch := make(chan int, n)
	for i := 0; i < n; i++ {
		pos := make([]int, n)
		pos[0] = i
		go func() {
			ch <- helper(1, pos)
		}()
	}
	for i := 0; i < n; i++ {
		num += <-ch
	}
	return num
}

type intlist struct {
	data int
	next *intlist
}

func nqueensg(n int) int {
	var pos *intlist = nil
	var helper func(int, *intlist, chan int)
	ch := make(chan int)
	helper = func(i int, pos *intlist, out chan int) {
		if i >= n {
			out <- 1
			return
		} else {
			num := 0
			newch := make(chan int)
			nok := 0
			for j := 0; j < n; j++ {
				ok := true
				poshere := pos
				for k := i - 1; ok && k >= 0; k-- {
					posk := poshere.data
					if j == posk || i+j == k+posk || i-j == k-posk {
						ok = false
					}
					poshere = poshere.next
				}
				if ok {
					nok++
					newpos := &intlist{j, pos}
					go helper(i+1, newpos, newch)
				}
			}
			for j := 0; j < nok; j++ {
				num += <-newch
			}
			out <- num
		}
	}
	go helper(0, pos, ch)
	return <-ch
}

func timefn(f func(int) int, arg int, name string) {
	t1 := time.Nanoseconds()
	result := f(arg)
	t2 := time.Nanoseconds()
	delta := float64(t2-t1) / 1000000000.0
	fmt.Printf("%g seconds to computed %d ways to place %d queens (%s)\n", delta, result, arg, name)
}

func main() {
	pnproc := flag.Int("p", -1, "maximum number of procs allowed")
	flag.Parse()
	n, err := strconv.Atoi(flag.Arg(0))
	if err != nil {
		n = 8
	}
	np := *pnproc
	if np <= 0 {
		np = n+1
	}
	runtime.GOMAXPROCS(np)
	timefn(nqueens, n, "simple")
//	  timefn(nqueens0, n, "simple")
	timefn(nqueens2, n, "goroutines")
}
