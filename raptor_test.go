package goraptor

import (
	"bytes"
	"testing"
)

func TestK4(t *testing.T) {
	p := getParameters(4)

	if !p.equal([]int{4, 4, 5, 5, 14, 17}) {
		t.Error(p.String())
	}

	//p.initA()
	a := [][]uint8{
		[]uint8{1, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		[]uint8{1, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0},
		[]uint8{1, 1, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0},
		[]uint8{0, 1, 1, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0},
		[]uint8{0, 0, 1, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0},
		[]uint8{1, 1, 0, 1, 1, 0, 0, 1, 0, 1, 0, 0, 0, 0},
		[]uint8{1, 0, 1, 1, 0, 1, 0, 0, 1, 0, 1, 0, 0, 0},
		[]uint8{1, 1, 1, 0, 0, 0, 1, 1, 1, 0, 0, 1, 0, 0},
		[]uint8{0, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 1, 0},
		[]uint8{0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 1},
		[]uint8{0, 1, 1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1, 0},
		[]uint8{0, 0, 0, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0},
		[]uint8{0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0},
		[]uint8{1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}

	for i := 0; i < len(a); i++ {
		for j := 0; j < len(a[0]); j++ {
			if p.a[i][j] != a[i][j] {
				t.Errorf("Un-match at %d:%d, where p.a is %d but test case is %d.\n", i, j, p.a[i][j], a[i][j])
			}
			if (a[i][j] == 1 && p.as.Has(i, j) == false) || (a[i][j] == 0 && p.as.Has(i, j)) {
				t.Errorf("Un-match at %d:%d, where a is %d but as is %v.\n", i, j, a[i][j], p.as.Has(i, j))
			}
		}
	}
}

func TestK(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	var testCases [][]int = [][]int{
		[]int{4, 4, 5, 5, 14, 17},
		[]int{5, 4, 5, 5, 15, 17},
		[]int{7, 5, 7, 6, 20, 23},
		[]int{340, 27, 31, 11, 382, 383},
		//		[]uint{1512, 56, 73, 13, 1598, 1601},
		//		[]uint{8162, 129, 211, 16, 8389, 8389},
		//		[]uint{8192, 129, 211, 16, 8419, 8419},
	}
	for _, c := range testCases {
		k := c[0]
		t.Logf("Testing K = %d", k)
		p := getParameters(k)
		if !p.equal(c) {
			t.Errorf("Test Failed when K=%d. Test case is %v but p is %v.", k, c, p)
			return
		}
	}
}

func TestMatrix(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	cases := []int{4, 8, 16, 32, 64, 128}
	for _, k := range cases {
		t.Logf("Testing matrix while k=%d", k)
		p := getParameters(k)
		p.initA()
		p.initAi()
		if p.a == nil || p.ai == nil {
			t.Errorf("Generate failed while k = %d", k)
			return
		}
	}
}

func TestRaptorEncode(t *testing.T) {
	src := [][]byte{[]byte{97}, []byte{98}, []byte{99}, []byte{100}}
	k := len(src)
	//Initialize
	raptor := New(k)
	for i, v := range src {
		raptor.Append(i, v)
	}
	raptor.Encode()

	//Check intermediate symbols
	exceptedC := []byte{97, 4, 102, 103, 6, 101, 3, 5, 1, 1, 4, 4, 101, 100}
	for idx, b := range raptor.c {
		if b[0] != exceptedC[idx] {
			t.Errorf("Expect %v at %v for L, but got %v.", exceptedC[idx], idx, b[0])
			return
		}
	}

	for i := 0; i < k; i++ {
		sym, _ := raptor.Symbol(i)
		if sym[0] != src[i][0] {
			t.Errorf("Expect %v at %v for symbol, but got %v.", src[i], i, sym)
			return
		}
	}
	repair := []byte{1, 2, 7, 5, 3, 99, 1, 4}
	for i := k; i < k+len(repair); i++ {
		sym, err := raptor.Symbol(i)
		if err != nil {
			t.Errorf("Get symbol %d failed, reason is %s.", i, err.Error())
			return
		}
		if sym[0] != repair[i-k] {
			t.Errorf("Expect %v at %v for symbol, but got %v.", repair[i], i, sym[0])
			return
		}
	}
}

func TestRaptorDecode(t *testing.T) {
	symbols := [][]byte{
		[]byte{97, 1}, []byte{98, 1}, []byte{99, 1}, []byte{100, 1},
		[]byte{1, 0}, []byte{2, 0}, []byte{7, 0}, []byte{5, 0},
		[]byte{3, 0}, []byte{99, 0}, []byte{1, 0}, []byte{4, 0},
	}
	k := 4
	raptor := New(k)
	for esi, data := range symbols {
		if data[1] == 1 {
			continue
		}
		raptor.Append(esi, data[:1])
	}
	if err := raptor.Decode(); err != nil {
		t.Errorf("Decode failed, error is %s.", err.Error())
		return
	}
	for esi := 0; esi < k; esi++ {
		sym, err := raptor.Symbol(esi)
		if err != nil {
			t.Errorf("Get symbol %d failed, error is %s.", esi, err.Error())
			return
		}
		if sym[0] != symbols[esi][0] {
			t.Errorf("Expect %v at %v for symbol, but got %v.", symbols[esi][0], esi, sym[0])
		}
	}
}

func TestRaptorEncodeAndDecode(t *testing.T) {
	symbols := make([][]byte, 128)
	for idx := range symbols {
		symbols[idx] = make([]byte, 1024)
		for offset := 0; offset < len(symbols[idx]); offset++ {
			symbols[idx][offset] = byte(idx + offset)
		}
	}

	k := len(symbols)
	encoder := New(k)
	for esi, data := range symbols {
		encoder.Append(esi, data)
	}
	encoder.Encode()

	decoder := New(k)
	for esi := 0; esi < 2*k; esi++ {
		if esi%3 == 0 {
			continue
		}
		sym, _ := encoder.Symbol(esi)
		decoder.Append(esi, sym)
	}
	decoder.Decode()

	for esi, data := range symbols {
		sym, _ := decoder.Symbol(esi)
		if !bytes.Equal(data, sym) {
			t.Errorf("Expect %v at %v for symbol, but got %v.", data, esi, sym)
			break
		}
	}
}
