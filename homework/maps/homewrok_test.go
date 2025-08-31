package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Node struct {
	key, value  int
	left, right *Node
}

type OrderedMap struct {
	root *Node
	size int
}

func NewOrderedMap() OrderedMap {
	return OrderedMap{}
}

func (m *OrderedMap) insert(node **Node, key, value int) {
	switch {
	case node == nil || (*node) == nil:
		*node = &Node{key: key, value: value}
		m.size++
	case key < (*node).key:
		m.insert(&(*node).left, key, value)
	case key > (*node).key:
		m.insert(&(*node).right, key, value)
	default:
		(*node).value = value
	}
}

func (m *OrderedMap) Insert(key, value int) {
	m.insert(&m.root, key, value)
}

func (m *OrderedMap) erase(node **Node, key int) {
	switch {
	case node == nil || (*node) == nil:
		return
	case key < (*node).key:
		m.erase(&(*node).left, key)
	case key > (*node).key:
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
			(*node).key = right.key
			(*node).value = right.value
			(*node).right = right.right
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

func (m *OrderedMap) Erase(key int) {
	m.erase(&m.root, key)
}

func findNode(node *Node, key int) *Node {
	switch {
	case node == nil:
		return nil
	case key < node.key:
		return findNode(node.left, key)
	case key > node.key:
		return findNode(node.right, key)
	default:
		return node
	}
}

func (m *OrderedMap) Contains(key int) bool {
	return findNode(m.root, key) != nil
}

func (m *OrderedMap) Size() int {
	return m.size
}

func forEach(node *Node, action func(int, int)) {
	if node == nil {
		return
	}

	forEach(node.left, action)
	action(node.key, node.value)
	forEach(node.right, action)
}

func (m *OrderedMap) ForEach(action func(int, int)) {
	forEach(m.root, action)
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap()
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
