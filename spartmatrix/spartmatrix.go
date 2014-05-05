package spartmatrix

import (
	"bytes"
	"container/list"
)

//Spart matrix in GF(2)
type SpartMatrix struct {
	Row, Colomn int
	rows        []*list.List
}

func (sm *SpartMatrix) Set(row, col int) {
	l := sm.rows[row]
	for e := l.Front(); e != nil; e = e.Next() {
		val := e.Value.(int)
		if val > col {
			l.InsertBefore(col, e)
			return
		}
		if val == col {
			return
		}
	}
	l.PushBack(col)
}

func (sm *SpartMatrix) Has(row, col int) bool {
	if row >= sm.Row || col >= sm.Colomn {
		return false
	}

	l := sm.rows[row]
	for e := l.Front(); e != nil; e = e.Next() {
		val := e.Value.(int)
		if val == col {
			return true
		}
		if val > col {
			return false
		}
	}
	return false
}

//SwapRow swap row r1 and row r2
func (sm *SpartMatrix) SwapRow(r1, r2 int) {
	if r1 == r2 {
		return
	}
	sm.rows[r1], sm.rows[r2] = sm.rows[r2], sm.rows[r1]
}

//XorRow let row[r1][i] = row[r1][i] ^ row[r2][i] for i in range(m.Colomn)
func (sm *SpartMatrix) XorRow(r1, r2 int) {
	row1 := sm.rows[r1]
	row2 := sm.rows[r2]
	e1 := row1.Front()
	e2 := row2.Front()
	for e1 != nil && e2 != nil {
		v1 := e1.Value.(int)
		v2 := e2.Value.(int)
		switch {
		case v1 == v2:
			// Remove() will set the Next pointer to nil, it must be saved first.
			tmp := e1.Next()
			row1.Remove(e1)
			e1 = tmp
			e2 = e2.Next()
		case v1 < v2:
			e1 = e1.Next()
		case v1 > v2:
			row1.InsertBefore(v2, e1)
			e2 = e2.Next()
		}
	}

	for ; e2 != nil; e2 = e2.Next() {
		row1.PushBack(e2.Value)
	}

}

func (sm *SpartMatrix) SwapColomn(c1, c2 int) {
	if c1 == c2 {
		return
	}
	for _, row := range sm.rows {
		sm.swapAColomn(row, c1, c2)
	}
}

func (sm *SpartMatrix) swapAColomn(row *list.List, c1, c2 int) {
	if c1 == c2 {
		return
	}

	if c1 > c2 {
		// let c1 < c2
		c1, c2 = c2, c1
	}

	var e1, e2 *list.Element
	for e := row.Front(); e != nil; e = e.Next() {
		val := e.Value.(int)
		if val >= c1 {
			e1 = e
		}

		if val >= c2 {
			e2 = e
			break
		}
	}

	foundC1 := (e1 != nil) && (c1 == e1.Value.(int))
	foundC2 := (e2 != nil) && (c2 == e2.Value.(int))
	switch {
	case !foundC1 && !foundC2:
		return
	case foundC1 && !foundC2:
		e1.Value = c2
		if e2 != nil {
			row.MoveBefore(e1, e2)
		} else {
			row.MoveToBack(e1)
		}
	case !foundC1 && foundC2:
		e2.Value = c1
		if e1 != nil {
			row.MoveBefore(e2, e1)
		} else {
			panic("Is it possible?")
		}
	case foundC1 && foundC2:
		return
	}
}

func (sm *SpartMatrix) Copy() *SpartMatrix {
	nsm := New(sm.Row, sm.Colomn)
	for idx, row := range sm.rows {
		for e := row.Front(); e != nil; e = e.Next() {
			nsm.rows[idx].PushBack(e.Value)
		}
	}
	return nsm
}

func (self *SpartMatrix) CountBetween(rowNum, start, end int) int {
	row := self.rows[rowNum]
	cnt := 0
	for e := row.Front(); e != nil; e = e.Next() {
		val := e.Value.(int)
		if val < start {
			continue
		} else if val >= start && val < end {
			cnt++
		} else {
			break
		}
	}
	return cnt
}

func (self *SpartMatrix) String() string {
	buffer := new(bytes.Buffer)
	for row := 0; row < self.Row; row++ {
		for col := 0; col < self.Colomn; col++ {
			if self.Has(row, col) {
				buffer.WriteString("1 ")
			} else {
				buffer.WriteString("0 ")
			}
		}
		buffer.WriteString("\n")
	}
	return buffer.String()
}

func (self *SpartMatrix) IterateRow(rowNum int, desc bool) *RowIterator {
	iter := new(RowIterator)
	iter.row = self.rows[rowNum]
	iter.desc = desc
	if !desc {
		iter.pos = iter.row.Front()
	} else {
		iter.pos = iter.row.Back()
	}

	return iter
}

type RowIterator struct {
	row  *list.List
	pos  *list.Element
	desc bool
}

func (self *RowIterator) HasNext() bool {
	return self.pos != nil
}

func (self *RowIterator) Next() {
	if !self.desc {
		self.pos = self.pos.Next()
	} else {
		self.pos = self.pos.Prev()
	}
}

func (self *RowIterator) Value() int {
	return self.pos.Value.(int)
}

func New(row, colomn int) *SpartMatrix {
	sm := new(SpartMatrix)
	sm.Row = row
	sm.Colomn = colomn
	sm.rows = make([]*list.List, row)
	for i := range sm.rows {
		sm.rows[i] = list.New()
		sm.rows[i].Init()
	}
	return sm
}
