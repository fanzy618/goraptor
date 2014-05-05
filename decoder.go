package goraptor

import (
	"goraptor/spartmatrix"
)

type decodeScheduler struct {
	dmatrix *spartmatrix.SpartMatrix
	l, m    int
	c, d    []int
}

func (self *decodeScheduler) xor(from, into int) {
	self.dmatrix.XorRow(self.d[into], self.d[from])
}

func (self *decodeScheduler) swapRow(i, j int) {
	self.d[i], self.d[j] = self.d[j], self.d[i]
}

func (self *decodeScheduler) swapCol(i, j int) {
	self.c[i], self.c[j] = self.c[j], self.c[i]
}

func findMinR(start, end int, matrix *spartmatrix.SpartMatrix) (rowNum, r int) {
	//FIXME
	for row := 0; row < matrix.Row; row++ {
		cnt := matrix.CountBetween(row, start, end)
		if cnt > r {
			r = cnt
			rowNum = row
		}
	}
	return
}

func (self *decodeScheduler) phaseOne(matrix *spartmatrix.SpartMatrix) {
	i, u := 0, 0
	for i+u < self.l {
		rowNum, r := findMinR(i, self.l-u, matrix)
		if r == 0 {
			//Decode failed.
			return
		}
		matrix.SwapRow(i, rowNum)
		self.SwapRow(i, rowNum)
	}
}

func (self *decodeScheduler) Decode(matrix *spartmatrix.SpartMatrix) *spartmatrix.SpartMatrix {
	return nil
}

func newScheduler(l, m int, D [][]byte) *decodeScheduler {
	scheduler := new(decodeScheduler)
	scheduler.m = m
	scheduler.dmatrix = spartmatrix.New(m, m)
	for i := 0; i < m; i++ {
		scheduler.dmatrix.Set(m, m)
	}

	scheduler.l = l
	scheduler.c = make([]int, l)
	for i := range scheduler.c {
		scheduler.c[i] = i
	}

	scheduler.d = make([]int, m)
	for i := range scheduler.d {
		scheduler.d[i] = i
	}

	return scheduler
}
