package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Node[K comparable, V any] struct {
	key         K
	value       V
	left, right *Node[K, V]
}

type OrderedMap[K comparable, V any] struct {
	root *Node[K, V]
	size int
	less func(K, K) bool
}

func NewOrderedMap[K comparable, V any](less func(K, K) bool) OrderedMap[K, V] {
	return OrderedMap[K, V]{less: less}
}

func (m *OrderedMap[K, V]) insert(node **Node[K, V], key K, value V) {
	if (*node) == nil {
		*node = &Node[K, V]{key: key, value: value}
		m.size++
		return
	}

	isLess := m.less(key, (*node).key)
	isGreater := (*node).key != key && !isLess
	switch {
	case isLess:
		m.insert(&(*node).left, key, value)
	case isGreater:
		m.insert(&(*node).right, key, value)
	default:
		(*node).value = value
	}
}

func (m *OrderedMap[K, V]) Insert(key K, value V) {
	m.insert(&m.root, key, value)
}

func (m *OrderedMap[K, V]) erase(node **Node[K, V], key K) {
	if (*node) == nil {
		return
	}

	isLess := m.less(key, (*node).key)
	isGreater := (*node).key != key && !isLess
	switch {
	case isLess:
		m.erase(&(*node).left, key)
	case isGreater:
		m.erase(&(*node).right, key)
	default:
		m.size--
		left := (*node).left
		right := (*node).right

		if right != nil {
			// find min in right
			for right.left != nil {
				right = right.left
			}

			(*node).right = right.right
			(*node).key = right.key
			(*node).value = right.value

			return
		}

		if left != nil {
			// replace current node with left node
			(*node).key = left.key
			(*node).value = left.value
			(*node).left = left.left
			(*node).right = left.right
			return
		}

		*node = nil
	}
}

func (m *OrderedMap[K, V]) Erase(key K) {
	m.erase(&m.root, key)
}

func (m *OrderedMap[K, V]) findNode(node *Node[K, V], key K) *Node[K, V] {
	if node == nil {
		return nil
	}

	isLess := m.less(key, (*node).key)
	isGreater := (*node).key != key && !isLess

	switch {
	case isLess:
		return m.findNode(node.left, key)
	case isGreater:
		return m.findNode(node.right, key)
	default:
		return node
	}
}

func (m *OrderedMap[K, V]) Contains(key K) bool {
	return m.findNode(m.root, key) != nil
}

func (m *OrderedMap[K, V]) Size() int {
	return m.size
}

func forEach[K comparable, V any](node *Node[K, V], action func(K, V)) {
	if node == nil {
		return
	}

	forEach(node.left, action)
	action(node.key, node.value)
	forEach(node.right, action)
}

func (m *OrderedMap[K, V]) ForEach(action func(K, V)) {
	forEach(m.root, action)
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap[int, int](func(a, b int) bool { return a < b })
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}

func TestInsertDuplicateKey(t *testing.T) {
	m := NewOrderedMap[int, int](func(a, b int) bool { return a < b })
	m.Insert(1, 10)
	m.Insert(1, 20) // duplicate key, should update value, not size
	assert.Equal(t, 1, m.Size())
	var val int
	m.ForEach(func(_, v int) { val = v })
	assert.Equal(t, 20, val)
}

func TestEraseNonExistent(t *testing.T) {
	m := NewOrderedMap[int, int](func(a, b int) bool { return a < b })
	m.Insert(1, 1)
	m.Insert(2, 2)
	m.Erase(3) // erase non-existent key
	assert.Equal(t, 2, m.Size())
}

func TestEraseFromEmpty(t *testing.T) {
	m := NewOrderedMap[int, int](func(a, b int) bool { return a < b })
	m.Erase(1) // should not panic
	assert.Zero(t, m.Size())
}

func TestInsertEraseSingle(t *testing.T) {
	m := NewOrderedMap[int, int](func(a, b int) bool { return a < b })
	m.Insert(42, 100)
	assert.Equal(t, 1, m.Size())
	m.Erase(42)
	assert.Zero(t, m.Size())
	assert.False(t, m.Contains(42))
}

func TestForEachEmpty(t *testing.T) {
	m := NewOrderedMap[int, int](func(a, b int) bool { return a < b })
	called := false
	m.ForEach(func(_, _ int) { called = true })
	assert.False(t, called)
}

func TestOrderAfterMixedOps(t *testing.T) {
	m := NewOrderedMap[int, int](func(a, b int) bool { return a < b })
	m.Insert(1, 50)
	m.Insert(2, 10)
	m.Insert(3, 30)
	m.Insert(4, 20)
	m.Insert(5, 40)

	m.Erase(3)

	m.Insert(3, 300)

	m.Erase(1)
	m.Insert(6, 60)

	var keys []int
	m.ForEach(func(k, _ int) { keys = append(keys, k) })
	expected := []int{2, 3, 4, 5, 6}
	assert.True(t, reflect.DeepEqual(expected, keys))
}

func TestNegativeAndLargeKeys(t *testing.T) {
	m := NewOrderedMap[int, int](func(a, b int) bool { return a < b })
	m.Insert(-10, 1)
	m.Insert(0, 2)
	m.Insert(1000000, 3)
	m.Insert(-100, 4)
	var keys []int
	m.ForEach(func(k, _ int) { keys = append(keys, k) })
	expected := []int{-100, -10, 0, 1000000}
	assert.True(t, reflect.DeepEqual(expected, keys))
}

func TestMinMaxIntKeys(t *testing.T) {
	m := NewOrderedMap[int, int](func(a, b int) bool { return a < b })
	min := -1 << 63
	max := 1<<63 - 1
	m.Insert(min, 111)
	m.Insert(max, 222)
	m.Insert(0, 333)
	var keys []int
	m.ForEach(func(k, _ int) { keys = append(keys, k) })
	expected := []int{min, 0, max}
	assert.True(t, reflect.DeepEqual(expected, keys))
	assert.True(t, m.Contains(min))
	assert.True(t, m.Contains(max))
	assert.True(t, m.Contains(0))
}

func TestStringKeysAndValues(t *testing.T) {
	m := NewOrderedMap[string, string](func(a, b string) bool { return a < b })
	m.Insert("b", "bee")
	m.Insert("a", "alpha")
	m.Insert("c", "cat")
	var keys []string
	var values []string
	m.ForEach(func(k, v string) {
		keys = append(keys, k)
		values = append(values, v)
	})
	assert.True(t, reflect.DeepEqual([]string{"a", "b", "c"}, keys))
	assert.True(t, reflect.DeepEqual([]string{"alpha", "bee", "cat"}, values))
}

func TestFloat64Keys(t *testing.T) {
	m := NewOrderedMap[float64, int](func(a, b float64) bool { return a < b })
	m.Insert(2.2, 22)
	m.Insert(1.1, 11)
	m.Insert(3.3, 33)
	var keys []float64
	m.ForEach(func(k float64, _ int) { keys = append(keys, k) })
	assert.True(t, reflect.DeepEqual([]float64{1.1, 2.2, 3.3}, keys))
}

func TestRuneKeys(t *testing.T) {
	m := NewOrderedMap[rune, string](func(a, b rune) bool { return a < b })
	m.Insert('b', "bee")
	m.Insert('a', "alpha")
	m.Insert('c', "cat")
	var keys []rune
	m.ForEach(func(k rune, _ string) { keys = append(keys, k) })
	assert.True(t, reflect.DeepEqual([]rune{'a', 'b', 'c'}, keys))
}

type person struct {
	Name string
	Age  int
}

func TestStructValues(t *testing.T) {
	m := NewOrderedMap[int, person](func(a, b int) bool { return a < b })
	m.Insert(2, person{"Bob", 30})
	m.Insert(1, person{"Alice", 25})
	m.Insert(3, person{"Carol", 40})
	var names []string
	m.ForEach(func(_ int, v person) { names = append(names, v.Name) })
	assert.True(t, reflect.DeepEqual([]string{"Alice", "Bob", "Carol"}, names))
}

type keyByField struct {
	ID  int
	Tag string
}

func TestStructKeysByField(t *testing.T) {
	m := NewOrderedMap[keyByField, string](func(a, b keyByField) bool { return a.ID < b.ID })
	m.Insert(keyByField{2, "b"}, "bee")
	m.Insert(keyByField{1, "a"}, "alpha")
	m.Insert(keyByField{3, "c"}, "cat")
	var ids []int
	m.ForEach(func(k keyByField, _ string) { ids = append(ids, k.ID) })
	assert.True(t, reflect.DeepEqual([]int{1, 2, 3}, ids))
}
