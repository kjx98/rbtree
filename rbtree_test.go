//
// Created by Jesse Kuang
//

package rbtree

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
	"testing"
)

const testVerbose = false

// Create a tree storing a set of integers
func testNewIntSet() *Tree[int, int] {
	return New[int, int](func(i1, i2 int) int {
		return i1 - i2
	})
}

func testAssert(t *testing.T, b bool, message string) {
	if !b {
		t.Fatal(message)
	}
}

func TestEmpty(t *testing.T) {
	tree := testNewIntSet()
	testAssert(t, tree.Len() == 0, "len!=0")
	testAssert(t, tree.Max().NegativeLimit(), "neglimit")
	testAssert(t, tree.Min().Limit(), "limit")
	testAssert(t, tree.FindGE(10).Limit(), "Not empty")
	testAssert(t, tree.FindLE(10).NegativeLimit(), "Not empty")
	testAssert(t, tree.Find(10) == nil, "Not empty")
	testAssert(t, tree.Limit().Equal(tree.Min()), "iter")
}

func TestFindGE(t *testing.T) {
	tree := testNewIntSet()
	testAssert(t, tree.Insert(10, 10), "Insert1")
	testAssert(t, !tree.Insert(10, 10), "Insert2")
	testAssert(t, tree.Len() == 1, "len==1")
	testAssert(t, tree.FindGE(10).Item() == 10, "FindGE 10")
	testAssert(t, tree.FindGE(11).Limit(), "FindGE 11")
	testAssert(t, tree.FindGE(9).Item() == 10, "FindGE 10")
}

func TestFindLE(t *testing.T) {
	tree := testNewIntSet()
	testAssert(t, tree.Insert(10, 10), "insert1")
	testAssert(t, tree.FindLE(10).Item() == 10, "FindLE 10")
	testAssert(t, tree.FindLE(11).Item() == 10, "FindLE 11")
	testAssert(t, tree.FindLE(9).NegativeLimit(), "FindLE 9")
}

func TestFind(t *testing.T) {
	tree := testNewIntSet()
	testAssert(t, tree.Insert(10, 10), "insert1")
	testAssert(t, *tree.Find(10) == 10, "Find 10")
	testAssert(t, tree.Find(9) == nil, "Find 9")
	testAssert(t, tree.Find(11) == nil, "Find 11")
}

func TestDelete(t *testing.T) {
	tree := testNewIntSet()
	testAssert(t, !tree.DeleteWithKey(10), "del")
	testAssert(t, tree.Len() == 0, "dellen")
	testAssert(t, tree.Insert(10, 10), "ins")
	testAssert(t, tree.DeleteWithKey(10), "del")
	testAssert(t, tree.Len() == 0, "dellen")

	// delete was deleting after the request if request not found
	// ensure this does not regress:
	testAssert(t, tree.Insert(10, 10), "ins")
	testAssert(t, !tree.DeleteWithKey(9), "del")
	testAssert(t, tree.Len() == 1, "dellen")

}

func iterToString(i *Iterator[int, int]) string {
	s := ""
	for ; !i.Limit(); i = i.Next() {
		if s != "" {
			s = s + ","
		}
		s = s + fmt.Sprintf("%d", i.Item())
	}
	return s
}

func reverseIterToString(i *Iterator[int, int]) string {
	s := ""
	for ; !i.NegativeLimit(); i = i.Prev() {
		if s != "" {
			s = s + ","
		}
		s = s + fmt.Sprintf("%d", i.Item())
	}
	return s
}

func TestIterator(t *testing.T) {
	tree := testNewIntSet()
	for i := 0; i < 10; i = i + 2 {
		tree.Insert(i, i)
	}
	if iterToString(tree.FindGE(3)) != "4,6,8" {
		t.Error("iter")
	}
	if iterToString(tree.FindGE(4)) != "4,6,8" {
		t.Error("iter")
	}
	if iterToString(tree.FindGE(8)) != "8" {
		t.Error("iter")
	}
	if iterToString(tree.FindGE(9)) != "" {
		t.Error("iter")
	}

	if reverseIterToString(tree.FindLE(3)) != "2,0" {
		t.Error("iter", reverseIterToString(tree.FindLE(3)))
	}
	if reverseIterToString(tree.FindLE(2)) != "2,0" {
		t.Error("iter")
	}
	if reverseIterToString(tree.FindLE(0)) != "0" {
		t.Error("iter")
	}
	if reverseIterToString(tree.FindLE(-1)) != "" {
		t.Error("iter")
	}
}

//
// Randomized tests
//

// oracle stores provides an interface similar to rbtree, but stores
// data in an sorted array
type oracle struct {
	data []int
}

func newOracle() *oracle {
	return &oracle{data: make([]int, 0)}
}

func (o *oracle) Len() int {
	return len(o.data)
}

// interface needed for sorting
func (o *oracle) Less(i, j int) bool {
	return o.data[i] < o.data[j]
}

func (o *oracle) Swap(i, j int) {
	e := o.data[j]
	o.data[j] = o.data[i]
	o.data[i] = e
}

func (o *oracle) Insert(key int) bool {
	for _, e := range o.data {
		if e == key {
			return false
		}
	}

	n := len(o.data) + 1
	newData := make([]int, n)
	copy(newData, o.data)
	newData[n-1] = key
	o.data = newData
	sort.Sort(o)
	return true
}

func (o *oracle) RandomExistingKey(rand *rand.Rand) int {
	index := rand.Intn(len(o.data))
	return o.data[index]
}

func (o *oracle) FindGE(t *testing.T, key int) oracleIterator {
	prev := int(-1)
	for i, e := range o.data {
		if e <= prev {
			t.Fatal("Nonsorted oracle ", e, prev)
		}
		if e >= key {
			return oracleIterator{o: o, index: i}
		}
	}
	return oracleIterator{o: o, index: len(o.data)}
}

func (o *oracle) FindLE(t *testing.T, key int) oracleIterator {
	iter := o.FindGE(t, key)
	if !iter.Limit() && o.data[iter.index] == key {
		return iter
	}
	return oracleIterator{o, iter.index - 1}
}

func (o *oracle) Delete(key int) bool {
	for i, e := range o.data {
		if e == key {
			newData := make([]int, len(o.data)-1)
			copy(newData, o.data[0:i])
			copy(newData[i:], o.data[i+1:])
			o.data = newData
			return true
		}
	}
	return false
}

//
// Test iterator
//
type oracleIterator struct {
	o     *oracle
	index int
}

func (oiter oracleIterator) Limit() bool {
	return oiter.index >= len(oiter.o.data)
}

func (oiter oracleIterator) Min() bool {
	return oiter.index == 0
}

func (oiter oracleIterator) NegativeLimit() bool {
	return oiter.index < 0
}

func (oiter oracleIterator) Max() bool {
	return oiter.index == len(oiter.o.data)-1
}

func (oiter oracleIterator) Item() int {
	return oiter.o.data[oiter.index]
}

func (oiter oracleIterator) Next() oracleIterator {
	return oracleIterator{oiter.o, oiter.index + 1}
}

func (oiter oracleIterator) Prev() oracleIterator {
	return oracleIterator{oiter.o, oiter.index - 1}
}

func compareContents(t *testing.T, oiter oracleIterator, titer *Iterator[int, int]) {
	oi := oiter
	ti := titer

	// Test forward iteration
	testAssert(t, oi.NegativeLimit() == ti.NegativeLimit(), "rend")
	if oi.NegativeLimit() {
		oi = oi.Next()
		ti = ti.Next()
	}

	for !oi.Limit() && !ti.Limit() {
		// log.Print("Item: ", oi.Item(), ti.Item())
		if ti.Item() != oi.Item() {
			t.Fatal("Wrong item", ti.Item(), oi.Item())
		}
		oi = oi.Next()
		ti = ti.Next()
	}
	if !ti.Limit() {
		t.Fatal("!ti.done", ti.Item())
	}
	if !oi.Limit() {
		t.Fatal("!oi.done", oi.Item())
	}

	// Test reverse iteration
	oi = oiter
	ti = titer
	testAssert(t, oi.Limit() == ti.Limit(), "end")
	if oi.Limit() {
		oi = oi.Prev()
		ti = ti.Prev()
	}

	for !oi.NegativeLimit() && !ti.NegativeLimit() {
		if ti.Item() != oi.Item() {
			t.Fatal("Wrong item", ti.Item(), oi.Item())
		}
		oi = oi.Prev()
		ti = ti.Prev()
	}
	if !ti.NegativeLimit() {
		t.Fatal("!ti.done", ti.Item())
	}
	if !oi.NegativeLimit() {
		t.Fatal("!oi.done", oi.Item())
	}
}

func compareContentsFull(t *testing.T, o *oracle, tree *Tree[int, int]) {
	compareContents(t, o.FindGE(t, int(-1)), tree.FindGE(-1))
}

func TestRandomized(t *testing.T) {
	const numKeys = 1000

	o := newOracle()
	tree := testNewIntSet()
	r := rand.New(rand.NewSource(0))
	for i := 0; i < 10000; i++ {
		op := r.Intn(100)
		if op < 50 {
			key := r.Intn(numKeys)
			if testVerbose {
				log.Print("Insert ", key)
			}
			o.Insert(key)
			tree.Insert(key, key)
			compareContentsFull(t, o, tree)
		} else if op < 90 && o.Len() > 0 {
			key := o.RandomExistingKey(r)
			if testVerbose {
				log.Print("DeleteExisting ", key)
			}
			o.Delete(key)
			if !tree.DeleteWithKey(key) {
				t.Fatal("DeleteExisting", key)
			}
			compareContentsFull(t, o, tree)
		} else if op < 95 {
			key := int(r.Intn(numKeys))
			if testVerbose {
				log.Print("FindGE ", key)
			}
			compareContents(t, o.FindGE(t, key), tree.FindGE(key))
		} else {
			key := int(r.Intn(numKeys))
			if testVerbose {
				log.Print("FindLE ", key)
			}
			compareContents(t, o.FindLE(t, key), tree.FindLE(key))
		}
	}
}

//
// Examples
//

func TestIntString(t *testing.T) {
	type MyItem struct {
		key   int
		value string
	}

	tree := New[string, int](func(a, b int) int { return a - b })
	tree.Insert(10, "value10")
	tree.Insert(12, "value12")

	log.Println("Find(10) ->", tree.Find(10))
	log.Println("Find(11) ->", tree.Find(11))

	// Find an element >= 11
	iter := tree.FindGE(11)
	fmt.Println("FindGE(11) ->", iter.Item())

	// Find an element >= 13
	iter = tree.FindGE(13)
	if !iter.Limit() {
		t.Fatal("There should be no element >= 13")
	}

	// Output:
	// Find(10) -> {10 value10}
	// Find(11) -> <nil>
	// FindGE(11) -> {12 value12}
}

func BenchmarkRBInsert(b *testing.B) {
	b.StopTimer()
	tree := New[int, int](func(a, b int) int {
		return a - b
	})
	for i := 0; i < 1e6; i++ {
		tree.Insert(i, i)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		v := (rand.Int() % 1e6) + 2e6
		tree.Insert(v, v)
	}
}

func BenchmarkFind(b *testing.B) {
	b.StopTimer()
	tree := New[int, int](func(a, b int) int {
		return a - b
	})
	for i := 0; i < 1e6; i++ {
		tree.Insert(i, i)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		v := (rand.Int() % 1e6)
		tree.Find(v)
	}
}
