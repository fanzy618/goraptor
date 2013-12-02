package goraptor

import (
	"errors"
	"fmt"
)

const MaxDupFactor = 3

var (
	EINVAL       = errors.New("Invalid argument")
	ENeedMore    = errors.New("Need More data to decode")
	EBufferSmall = errors.New("Buffer is too small")
)

type Raptor struct {
	k        int
	p        *parameters
	blockLen int
	data     [][]byte
	dataCnt  int
	c        [][]byte
	// Underlying buffer for l
	bufferC []byte
	// d = p.a * l
	d [][]byte
}

func New(k int) *Raptor {
	if k < 4 || k > 8192 {
		return nil
	}
	var raptor *Raptor = new(Raptor)
	raptor.k = k
	raptor.p = getParameters(k)
	raptor.data = make([][]byte, MaxDupFactor*k)
	return raptor
}

func (raptor *Raptor) Append(esi int, data []byte) error {
	if esi < 0 || esi >= MaxDupFactor*raptor.k || data == nil {
		return EINVAL
	}
	if raptor.blockLen == 0 {
		raptor.blockLen = len(data)
	} else {
		if raptor.blockLen != len(data) {
			return EINVAL
		}
	}

	if raptor.data[esi] == nil {
		raptor.dataCnt++
	}
	raptor.data[esi] = data
	return nil
}

func (raptor *Raptor) Symbol(esi int) ([]byte, error) {
	if esi < 0 || esi >= MaxDupFactor*raptor.k {
		return nil, EINVAL
	}

	if esi < raptor.k && raptor.data[esi] != nil {
		return raptor.data[esi], nil
	}
	//	if raptor.data[esi] == nil {
	p := raptor.p
	d, a, b := p.lt_triple(esi)

	for ; b >= p.l; b = (b + a) % raptor.p.lPrime {
	}

	if d > p.l {
		d = p.l
	}

	var data []byte
	if raptor.data[esi] == nil {
		raptor.data[esi] = make([]byte, raptor.blockLen)
	}
	data = raptor.data[esi]
	if len(data) != len(raptor.c[b]) {
		str := fmt.Sprintf("Length unmatch. data is %v, c[%v] is %v.", len(data), b, len(raptor.c[b]))
		return nil, errors.New(str)
	}
	copy(data, raptor.c[b])

	for i := 1; i < d; i++ {
		b = (b + a) % p.lPrime
		for ; b >= p.l; b = (b + a) % p.lPrime {
		}
		for j := range data {
			data[j] ^= raptor.c[b][j]
		}
	}

	//	}

	return raptor.data[esi], nil
}

func (raptor *Raptor) Encode() error {
	if raptor.dataCnt < raptor.k {
		return ENeedMore
	}
	for i := 0; i < raptor.k; i++ {
		//The source symbols must have been appended.
		if raptor.data[i] == nil {
			return ENeedMore
		}
	}

	raptor.p.initAi()
	if raptor.bufferC == nil {
		raptor.bufferC = make([]byte, raptor.p.l*raptor.blockLen)
		raptor.c = make([][]byte, raptor.p.l)
		for i := 0; i < int(raptor.p.l); i++ {
			raptor.c[i] = raptor.bufferC[i*raptor.blockLen : (i+1)*raptor.blockLen]
		}
	}

	sh := int(raptor.p.s + raptor.p.h)
	for row, rowData := range raptor.p.ai {
		for col := sh; col < len(rowData); col++ {
			if rowData[col] == 1 {
				for i := 0; i < raptor.blockLen; i++ {
					raptor.c[row][i] ^= raptor.data[col-sh][i]
				}
			}
		}
	}

	return nil
}

func (raptor *Raptor) Decode() error {
	if raptor.dataCnt < raptor.k {
		return ENeedMore
	}
	sh := int(raptor.p.s + raptor.p.h)
	// Matrix A for decode
	ad := make([][]uint8, raptor.dataCnt+sh)
	temp := make([][]byte, raptor.dataCnt+sh)
	for i := 0; i < raptor.dataCnt+sh; i++ {
		temp[i] = make([]byte, raptor.blockLen)
		ad[i] = make([]uint8, raptor.p.l)
		if i < sh {
			copy(ad[i], raptor.p.a[i])
		}
	}

	offset := 0
	for i := sh; i < len(ad); i++ {
		for ; raptor.data[i+offset-sh] == nil; offset++ {

		}
		raptor.p.ltRow(i+offset-sh, ad[i])
		copy(temp[i], raptor.data[i+offset-sh])
	}

	var l int = int(raptor.p.l)
	for k := 0; k < l; k++ {
		if ad[k][k] == 0 {
			var i int = k
			for ; i < len(ad); i++ {
				if ad[i][k] == 1 {
					ad[k], ad[i] = ad[i], ad[k]
					temp[k], temp[i] = temp[i], temp[k]
					break
				}
			}
			if i == len(ad) {
				return ENeedMore
			}
		}
		for i := 0; i < len(ad); i++ {
			if i == k {
				continue
			}
			if ad[i][k] == 1 {
				xorRow(ad[i], ad[k])
				xorRow(temp[i], temp[k])
				//				for j := 0; j < raptor.blockLen; j++ {
				//					temp[i][j] ^= temp[k][j]
				//				}
			}
		}
	}
	raptor.c = make([][]byte, raptor.p.l)
	for i := 0; i < raptor.p.l; i++ {
		raptor.c[i] = temp[i]
	}

	return nil
}
