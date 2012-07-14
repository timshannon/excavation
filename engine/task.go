package engine

import (
	"sort"
)

//task state
const (
	TaskStopped = iota
	TaskRunning
	TaskWaiting
	TaskCompleted
)

var (
	taskList  Tasks
	taskQueue Tasks
	sorted    bool
)

type taskFunc func(task *Task)

type Task struct {
	Name     string
	Func     taskFunc
	Time     float64
	Frames   int
	priority int
	Data     map[string]interface{}
	state    uint
	delay    float64
}

func (t *Task) Priority() int { return t.priority }
func (t *Task) SetPriority(priority int) {
	t.priority = priority
	//priority changed, resort
	sorted = false
}

func (t *Task) Wait(seconds float64) {
	t.delay = seconds
	t.state = TaskWaiting
}

type Tasks []*Task

func (t Tasks) Len() int      { return len(t) }
func (t Tasks) Swap(i, j int) { t[i], t[j] = t[j], t[i] }

type ByPriority struct{ Tasks }

func (t ByPriority) Less(i, j int) bool { return t.Tasks[i].priority < t.Tasks[j].priority }

func AddTask(name string, function taskFunc, priority int, delay float64) {
	task := &Task{name, function, 0, 0, priority, nil, TaskWaiting, delay}
	taskList = append(taskList, task)

	//new task added resort
	sorted = false
}

//runTasks sorts the taskList by priority, then adds all active,
// non-waiting task to the task queue and processes them
func runTasks() {

	if !sorted {
		sort.Sort(ByPriority{taskList})
		sorted = true
	}

	for i, task := range taskList {
		switch task.state {
		case TaskCompleted:
			//remove task from list
			taskList = append(taskList[:i], taskList[i+1:]...)
		case TaskWaiting:
			//check delay
			task.delay = task.delay - (Time() - task.Time)
			if task.delay <= 0 {
				task.delay = 0
				task.state = TaskRunning
			}
			fallthrough
		case TaskRunning:
			taskQueue = append(taskQueue, task)
		case TaskStopped:
			//do nothing
		}

	}

	//run through all queued tasks
	for _, task := range taskQueue {
		task.Time = Time()

	}
}
