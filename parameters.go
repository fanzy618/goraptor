package goraptor

import (
	"fmt"
	"math"
)

type parameters struct {
	k, x, s, h, l, lPrime int
	a, ai                 [][]uint8
}

func (p *parameters) String() string {
	return fmt.Sprintf("{K:%d, X:%d, S:%d, H:%d, L:%d, L_Prime:%d}",
		p.k, p.x, p.s, p.h, p.l, p.lPrime)
}

func (p *parameters) equal(s []int) bool {
	if p.k != s[0] || p.x != s[1] || p.s != s[2] || p.h != s[3] || p.l != s[4] || p.lPrime != s[5] {
		return false
	}
	return true
}

func isPrime(p int) bool {
	var r = int(math.Ceil(math.Sqrt(float64(p))))
	var i int = 2
	for ; i <= r; i++ {
		if p%i == 0 {
			return false
		}
	}
	return true
}

func nextPrime(p int) int {
	if p%2 == 0 {
		p++
	}
	for ; !isPrime(p); p += 2 {
	}
	return p
}

func findX(k int) (x int) {
	for x = 4; x*(x-1) < 2*k; x++ {
	}
	return x
}

//n!
func factorial(n int) int {
	switch n {
	case 0, 1:
		return 1
	}
	var p, i int = 1, 2
	for ; i <= n; i++ {
		p *= i
	}
	return p
}

func choose(n, m int) int {
	if n < m {
		return 0
	}
	r := factorial(n) / factorial(n-m) / factorial(m)
	return r
}

func findH(n int) int {
	var h int
	for h = 2; choose(h, int(math.Ceil(float64(h/2)))) < n; h++ {
	}
	return h
}

func getParameters(k int) (p *parameters) {
	p = &parameters{k: k}
	p.x = findX(k)
	p.s = nextPrime(int(math.Ceil(0.01*float64(p.k))) + p.x)
	p.h = findH(p.k + p.s)
	p.l = p.k + p.s + p.h
	p.lPrime = nextPrime(p.l)
	//For the p that has the same k, matrix a and ai are the same
	p.initA()

	return p
}
