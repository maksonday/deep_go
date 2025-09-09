package main

import (
	"container/heap"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Task struct {
	Identifier int
	Priority   int
}

type TaskQueue struct {
	queue []Task
	mu    sync.Mutex
}

func (t *TaskQueue) Len() int           { return len(t.queue) }
func (t *TaskQueue) Less(i, j int) bool { return t.queue[i].Priority >= t.queue[j].Priority }
func (t *TaskQueue) Swap(i, j int)      { t.queue[i], t.queue[j] = t.queue[j], t.queue[i] }

func (t *TaskQueue) Push(x any) {

	t.queue = append(t.queue, x.(Task))
}

func (t *TaskQueue) Pop() any {
	old := t.queue
	n := len(old)
	x := old[n-1]
	t.queue = old[:n-1]

	return x
}

type Scheduler struct {
	queue *TaskQueue
}

func NewScheduler() Scheduler {
	queue := TaskQueue{queue: []Task{}}
	heap.Init(&queue)

	return Scheduler{
		queue: &queue,
	}
}

func (s *Scheduler) AddTask(task Task) {
	s.queue.mu.Lock()
	defer s.queue.mu.Unlock()
	heap.Push(s.queue, task)
}

func (s *Scheduler) ChangeTaskPriority(taskID int, newPriority int) {
	for i := range s.queue.Len() {
		if s.queue.queue[i].Identifier == taskID {
			s.queue.queue[i].Priority = newPriority
			heap.Fix(s.queue, i)
			break
		}
	}
}

func (s *Scheduler) GetTask() Task {
	s.queue.mu.Lock()
	defer s.queue.mu.Unlock()
	return heap.Pop(s.queue).(Task)
}

func TestTrace(t *testing.T) {
	task1 := Task{Identifier: 1, Priority: 10}
	task2 := Task{Identifier: 2, Priority: 20}
	task3 := Task{Identifier: 3, Priority: 30}
	task4 := Task{Identifier: 4, Priority: 40}
	task5 := Task{Identifier: 5, Priority: 50}

	scheduler := NewScheduler()
	scheduler.AddTask(task1)
	scheduler.AddTask(task2)
	scheduler.AddTask(task3)
	scheduler.AddTask(task4)
	scheduler.AddTask(task5)

	task := scheduler.GetTask()
	assert.Equal(t, task5, task)

	task = scheduler.GetTask()
	assert.Equal(t, task4, task)

	scheduler.ChangeTaskPriority(1, 100)

	task1.Priority = 100

	task = scheduler.GetTask()
	assert.Equal(t, task1, task)

	task = scheduler.GetTask()
	assert.Equal(t, task3, task)
}
