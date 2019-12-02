package taskpool

import (
	"container/list"
	"sync"
)

type task struct {
	id     string
	fn     interface{}
	params []interface{}
}

type queue struct {
	name  string
	store map[string]*list.Element
	ll    *list.List
	cond  *sync.Cond
}

// cover指明当元素已存在时是否覆盖
func (q *queue) push(id string, fn interface{}, cover bool, params ...interface{}) (int, bool) {
	q.cond.L.Lock()

	defer func() {
		q.cond.Broadcast()
		q.cond.L.Unlock()
	}()

	ok := false
	var ele *list.Element
	if ele, ok = q.store[id]; ok {
		if !cover {
			return len(q.store), false
		}
	}

	task := task{
		id:     id,
		fn:     fn,
		params: params,
	}

	if ok {
		q.store[id] = q.ll.InsertAfter(&task, ele)
		q.ll.Remove(ele)
	} else {
		q.store[id] = q.ll.PushBack(&task)
	}
	return len(q.store), true
}

func (q *queue) pop() *task {
	q.cond.L.Lock()

	for q.ll.Len() == 0 {
		q.cond.Wait()
	}

	defer func() {
		q.cond.L.Unlock()
	}()

	ele := q.ll.Front()
	if ele == nil {
		return nil
	}

	task := ele.Value.(*task)
	q.ll.Remove(ele)
	delete(q.store, task.id)
	return task
}
