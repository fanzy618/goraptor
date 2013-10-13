package goraptor

import (
//	"fmt"
)

const MaxDupFactor = 3

type Error struct {
	errstr string
}

func (err *Error) Error() string {
	return err.errstr
}

var (
	EINVAL       = &Error{"Invalid argument"}
	ENeedMore    = &Error{"Need More data to decode"}
	EBufferSmall = &Error{"Buffer is too small"}
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
	raptor.data = make([][]byte, MaxDupFactor * k)
	return raptor
}

func (raptor *Raptor) append(esi int, data []byte) error {
	if esi < 0 || esi >= MaxDupFactor * raptor.k || data == nil {
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
	if esi < 0 || esi >= MaxDupFactor * raptor.k {
		return nil, EINVAL
	}

	if raptor.data[esi] == nil {
		p := raptor.p
		d, a, b := p.lt_triple(esi)

		for ; b >= p.l; b = (b + a) % raptor.p.lPrime {
		}

		if d > p.l {
			d = p.l
		}

		data := make([]byte, raptor.blockLen)
		copy(data, raptor.c[b])

		for i := 1; i < d; i++ {
			b = (b + a) % p.lPrime
			for ; b >= p.l; b = (b + a) % p.lPrime {
			}
			for j := range data {
				data[j] ^= raptor.c[b][j]
			}
		}
		raptor.data[esi] = data
	}

	return raptor.data[esi], nil
}

func (raptor *Raptor) Encode() error {
	if raptor.dataCnt != raptor.k {
		return EINVAL
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
	for i := range ad {
		temp[i] = make([]byte, raptor.dataCnt)
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
				for j := 0; j < raptor.blockLen; j++ {
					temp[i][j] ^= temp[k][j]
				}
			}
		}
	}
	raptor.c = make([][]byte, raptor.p.l)
	for i := range raptor.c {
		raptor.c[i] = temp[i]
	}

	return nil
}
