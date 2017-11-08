package models

import (
	"reflect"
	"sort"
)

// LinesByID implements sort.Interface for []*TransactionLine based on
// the AccountID field.
type LinesByID []*TransactionLine

func (lines LinesByID) Len() int           { return len(lines) }
func (lines LinesByID) Swap(i, j int)      { lines[i], lines[j] = lines[j], lines[i] }
func (lines LinesByID) Less(i, j int) bool { return lines[i].AccountID < lines[j].AccountID }

func containsSameElements(l1 []*TransactionLine, l2 []*TransactionLine) bool {
	lc1 := make([]*TransactionLine, len(l1))
	copy(lc1, l1)
	lc2 := make([]*TransactionLine, len(l2))
	copy(lc2, l2)
	sort.Sort(LinesByID(lc1))
	sort.Sort(LinesByID(lc2))
	return reflect.DeepEqual(lc1, lc2)
}
