# taskpool
golang task pool

# Usage

taskPool := NewTaskPool("My task pool", 10) // 创建10个协程的协程池
taskPool.AddTask(id, workFunc, cover, params)
