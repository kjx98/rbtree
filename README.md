### rbtree - RedBlack tree
#### Jesse Kuang <jkuang@21cn.com>

[![GoDoc](https://godoc.org/github.com/kjx98/rbtree.git?status.svg)](https://godoc.org/github.com/kjx98/rbtree.git)

A red-black tree with an API similar to C++ STL's.

### INSTALLATION
    go get github.com/kjx98/rbtree

### EXAMPLE

    More examples can be found in rbtree_test.go

```
    import "github.com/kjx98/rbtree"

	type MyItem struct {
		key   int
		value string
	}

	tree := rbtree.New(func(a, b interface {}) int { return a.(MyItem).key - b.(MyItem).key })
	tree.Insert(MyItem{10, "value10"})
	tree.Insert(MyItem{12, "value12"})

	fmt.Println("Find(10) ->", tree.Find(MyItem{10, ""}))
	fmt.Println("Find(11) ->", tree.Find(MyItem{11, ""}))

	// Find an element >= 11
	iter := tree.FindGE(MyItem{11, ""})
	fmt.Println("FindGE(11) ->", iter.Item())

	// Find an element >= 13
	iter = tree.FindGE(MyItem{13, ""})
	if !iter.End() { panic("There should be no element >= 13") }

	// Output:
	// Find(10) -> {10 value10}
	// Find(11) -> <nil>
	// FindGE(11) -> {12 value12}
```

### TYPES

```
type CompareFunc func(a, b interface {}) int
    CompareFunc returns 0 if a==b, <0 if a<b, >0 if a>b.

type Iterator struct {
    // contains filtered or unexported fields
}
    Iterator allows scanning tree elements in sort order.

    Iterator invalidation rule is the same as C++ std::map<>'s. That is, if
    you delete the element that an iterator points to, the iterator becomes
    invalid. For other operation types, the iterator remains valid.

func (iter *Iterator) Equal(iter2 *Iterator) bool

func (iter *Iterator) Item() interface{}
    Return the current element.

    REQUIRES: !iter.Limit() && !iter.NegativeLimit()

func (iter *Iterator) Limit() bool
    Check if the iterator points beyond the max element in the tree

func (iter *Iterator) Max() bool
    Check if the iterator points to the maximum element in the tree

func (iter *Iterator) Min() bool
    Check if the iterator points to the minimum element in the tree

func (iter *Iterator) NegativeLimit() bool
    Check if the iterator points before the minumum element in the tree

func (iter *Iterator) Next() *Iterator
    Create a new iterator that points to the successor of the current
    element.

    REQUIRES: !iter.Limit()

func (iter *Iterator) Prev() *Iterator
    Create a new iterator that points to the predecessor of the current
    node.

    REQUIRES: !iter.NegativeLimit()

type Tree struct {
    // contains filtered or unexported fields
}

func NewTree(compare CompareFunc) *Tree
    Create a new empty tree.

func (root *Tree) DeleteWithIterator(iter *Iterator)
    Delete the current item.

    REQUIRES: !iter.Limit() && !iter.NegativeLimit()

func (root *Tree) DeleteWithKey(key interface {}) bool
    Delete an item with the given key. Return true iff the item was found.

func (root *Tree) FindGE(key interface {}) *Iterator
    Find the smallest element N such that N >= key, and return the iterator
    pointing to the element. If no such element is found, return
    root.Limit().

func (root *Tree) FindLE(key interface {}) *Iterator
    Find the largest element N such that N <= key, and return the iterator
    pointing to the element. If no such element is found, return
    iter.NegativeLimit().

func (root *Tree) Find(key interface {}) interface {}
    A convenience function for finding an element equal to key. Return nil
    if not found.

func (root *Tree) Insert(item interface {}) bool
    Insert an item. If the item is already in the tree, do nothing and
    return false. Else return true.

func (root *Tree) Len() int
    Return the number of elements in the tree.

func (root *Tree) Limit() *Iterator
    Create an iterator that points beyond the maximum item in the tree

func (root *Tree) Max() *Iterator
    Create an iterator that points at the maximum item in the tree

    If the tree is empty, return NegativeLimit()

func (root *Tree) Min() *Iterator
    Create an iterator that points to the minimum item in the tree If the
    tree is empty, return Limit()

func (root *Tree) NegativeLimit() *Iterator
    Create an iterator that points before the minimum item in the tree
```
