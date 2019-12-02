package taskpool

import (
	"container/list"
	"fmt"
	"reflect"
	"sync"
)

type TaskPool struct {
	Name  string
	queue queue
	sem   chan struct{}
}

func NewTaskPool(name string, count int) *TaskPool {
	pool := TaskPool{Name: name}
	pool.queue = queue{
		name:  name,
		store: make(map[string]*list.Element),
		ll:    list.New(),
		cond:  sync.NewCond(new(sync.Mutex)),
	}
	pool.sem = make(chan struct{}, count)

	go pool.run()
	return &pool
}

func (pool *TaskPool) AddTask(id string, fn interface{}, cover bool, params ...interface{}) (int, bool) {
	return pool.queue.push(id, fn, cover, params)
}

func (pool *TaskPool) run() {
	for {
		t := pool.queue.pop()
		if t != nil {
			pool.sem <- struct{}{}
			go pool.work(t)
		}
	}
}

func (pool *TaskPool) work(t *task) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recover:", r)
		}
	}()

	f := reflect.ValueOf(t.fn)
	args := make([]reflect.Value, 0, len(t.params))

	for _, p := range t.params {
		args = append(args, reflect.ValueOf(p))
	}
	f.Call(args)
	<-pool.sem
}
