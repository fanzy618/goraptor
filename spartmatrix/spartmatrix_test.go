package spartmatrix

import (
	"testing"
)

func getIdentityMatrix(i int) *SpartMatrix {
	m := New(i, i)
	for x := 0; x < i; x++ {
		m.Set(x, x)
	}
	return m
}

func TestSetAndHas(t *testing.T) {
	sm := New(3, 3)
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			sm.Set(i, j)
		}
	}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if !sm.Has(i, j) {
				t.Errorf("Matrix should have %d, %d but not.\n", i, j)
			}
		}
	}

	if sm.Has(3, 3) {
		t.Errorf("Matrix should not have 3,3!\n")
	}
}

func TestSwapRow(t *testing.T) {
	m := getIdentityMatrix(2)
	m.SwapRow(0, 1)
	if !m.Has(0, 1) || !m.Has(1, 0) {
		t.Errorf("Test failed!")
	}
	if m.Has(0, 0) || m.Has(1, 1) {
		t.Errorf("Test failed!")
	}
}

func TestXorRow(t *testing.T) {
	m := getIdentityMatrix(2)
	m.XorRow(0, 1)
	if !(m.Has(0, 0) && m.Has(0, 1)) {
		t.Errorf("Test failed in row 0! Matrix is:\n%s", m.String())
		return
	}
	m.XorRow(1, 0)
	if !(m.Has(1, 0) && !m.Has(1, 1)) {
		t.Errorf("Test failed in row 0! Matrix is:\n%s", m.String())
	}

}

func TestSwapColomn(t *testing.T) {
	m := getIdentityMatrix(2)
	m.SwapColomn(0, 1)
	if !(!m.Has(0, 0) && m.Has(0, 1) && m.Has(1, 0) && !m.Has(1, 1)) {
		t.Errorf("Test failed! Matrix is:\n%s", m.String())
	}

}
