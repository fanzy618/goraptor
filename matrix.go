package goraptor

import (
	"fmt"
	"math"
)

func (p *parameters) g_ldpc() {
	//G_LDPC
	var a, b int
	for i := 0; i < p.k; i++ {
		a = 1 + int(math.Floor(float64(i)/float64(p.s)))%(p.s-1)
		b = i % p.s
		p.a[b][i] = uint8(1)
		b = (b + a) % p.s
		p.a[b][i] = 1
		b = (b + a) % p.s
		p.a[b][i] = 1
	}
	//I_S
	for i := 0; i < p.s; i++ {
		p.a[i][i+p.k] = 1
	}
}

func (p *parameters) ldpc_as() {
	var a, b int
	for i := 0; i < p.k; i++ {
		a = 1 + int(math.Floor(float64(i)/float64(p.s)))%(p.s-1)
		b = i % p.s
		p.as.Set(b, i)

		b = (b + a) % p.s
		p.as.Set(b, i)
		b = (b + a) % p.s
		p.as.Set(b, i)
	}

	for i := 0; i < p.s; i++ {
		p.as.Set(i, i+p.k)
	}
}

func countOne(i int) int {
	var cnt int = 0
	for ; i > 0; i >>= 1 {
		if i&1 != 0 {
			cnt += 1
		}
	}
	return cnt
}

/*
   Let

      g[i] = i ^ (floor(i/2)) for all positive integers i

         Note: g[i] is the Gray sequence, in which each element differs
         from the previous one in a single bit position

      m[k] denote the subsequence of g[.] whose elements have exactly k non-zero bits in their binary representation.
      m[j,k] denote the jth element of the sequence m[k], where j=0, 1, 2, ...
*/
func genSeqM(j, h int) []int {
	var m []int = make([]int, 0, 256)
	for i := 1; len(m) < j; i++ {
		g := i ^ int(math.Floor(float64(i)/2.0))
		if countOne(g) == h {
			m = append(m, g)
		}
	}
	return m
}

func (p *parameters) g_half() {
	seqM := genSeqM(p.k+p.s, int(math.Ceil(float64(p.h)/2.0)))
	for i := 0; i < p.h; i++ {
		//G_Half
		for j := 0; j < p.k+p.s; j++ {
			if (1<<uint(i))&seqM[j] != 0 {
				p.a[p.s+i][j] = 1
			}
		}
		//I_H
		p.a[p.s+i][p.k+p.s+i] = 1
	}
}

func (p *parameters) half_as() {
	seqM := genSeqM(p.k+p.s, int(math.Ceil(float64(p.h)/2.0)))
	for i := 0; i < p.h; i++ {
		//G_Half
		for j := 0; j < p.k+p.s; j++ {
			if (1<<uint(i))&seqM[j] != 0 {
				p.as.Set(p.s+i, j)
			}
		}
		//I_H
		p.as.Set(p.s+i, p.k+p.s+i)
	}
}

func (p *parameters) ltRow(id int, buffer []uint8) {
	d, a, b := p.lt_triple(id)

	for ; b >= p.l; b = (b + a) % p.lPrime {
	}
	buffer[b] = 1

	if d > p.l {
		d = p.l
	}

	for j := 1; j < d; j++ {
		b = (b + a) % p.lPrime
		for ; b >= p.l; b = (b + a) % p.lPrime {
		}
		buffer[b] = 1
	}
}

func (p *parameters) g_lt() {
	for id := 0; id < p.k; id++ {
		rowIdx := id + p.s + p.h
		p.ltRow(id, p.a[rowIdx])
	}
}

func (p *parameters) lt_as() {
	for id := 0; id < p.k; id++ {
		rowIdx := id + p.s + p.h

		d, a, b := p.lt_triple(id)

		for ; b >= p.l; b = (b + a) % p.lPrime {
		}
		p.as.Set(rowIdx, b)

		if d > p.l {
			d = p.l
		}

		for j := 1; j < d; j++ {
			b = (b + a) % p.lPrime
			for ; b >= p.l; b = (b + a) % p.lPrime {
			}
			p.as.Set(rowIdx, b)
		}
	}
}

func (p *parameters) initA() {
	if p.a != nil {
		return
	}
	p.a = make([][]uint8, p.l)
	for i := 0; i < p.l; i++ {
		p.a[i] = make([]uint8, p.l)
	}
	p.g_ldpc()
	p.g_half()
	p.g_lt()
}

func (p *parameters) initAs() {

	p.ldpc_as()
	p.half_as()
	p.lt_as()
}

//For i = 0,.. n, r1[i] = r1[i] ^ r2[i]
func xorRow(r1, r2 []uint8) {
	if len(r1) != len(r2) {
		return
	}
	for idx, v := range r1 {
		r1[idx] = v ^ r2[idx]
	}
}

//Initialize the inverst matrix of A
func (p *parameters) initAi() {
	if p.ai != nil {
		return
	}
	var a, ai [][]uint8 = make([][]uint8, p.l), make([][]uint8, p.l)
	var l int = int(p.l)
	for i := 0; i < l; i++ {
		a[i], ai[i] = make([]uint8, p.l), make([]uint8, p.l)
	}

	for i := 0; i < l; i++ {
		ai[i][i] = 1
		for j := 0; j < l; j++ {
			a[i][j] = p.a[i][j]
		}
	}

	for k := 0; k < l; k++ {
		if a[k][k] == 0 {
			var i int = k
			for ; i < l; i++ {
				if a[i][k] == 1 {
					a[k], a[i] = a[i], a[k]
					ai[k], ai[i] = ai[i], ai[k]
					break
				}
			}
			if i == l {
				return
			}
		}
		for i := 0; i < l; i++ {
			if i == k {
				continue
			}
			if a[i][k] == 1 {
				xorRow(a[i], a[k])
				xorRow(ai[i], ai[k])
			}

		}
	}

	p.ai = ai
}

//PrintMatrix prints the 'a' and 'ai' for debug.
func (p *parameters) PrintMatrix() {
	fmt.Printf("Matrix A is:\n")
	for idx, row := range p.a {
		fmt.Printf("%v:\t", idx)
		for _, cell := range row {
			fmt.Printf("%v, ", cell)
		}
		fmt.Println()
	}
	fmt.Println()

	fmt.Printf("Matrix AI is:\n")
	for idx, row := range p.ai {
		fmt.Printf("%v:\t", idx)
		for _, cell := range row {
			fmt.Printf("%v, ", cell)
		}
		fmt.Println()
	}
	fmt.Println()
}
