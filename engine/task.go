// Copyright 2012 Tim Shannon. All rights reserved. 
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file. 

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
	taskList    Tasks
	taskQueue   Tasks
	tasksSorted bool
)

type taskFunc func(task *Task)

type Task struct {
	Name     string
	Func     taskFunc
	Data     interface{}
	start    float64
	frames   int
	state    uint
	delay    float64
	priority int
}

func (t *Task) Priority() int { return t.priority }
func (t *Task) SetPriority(priority int) {
	t.priority = priority
	//priority changed, resort
	tasksSorted = false
}

//Wait schedules the task to run the given # of seconds in the future
func (t *Task) Wait(seconds float64) {
	t.delay = Time() + seconds
	t.state = TaskWaiting
}

//Time is the number of seconds passed since this task first started
func (t *Task) Time() float64 {
	return Time() - t.start
}

//Frames is the number of frames/times this task has been called
func (t *Task) Frames() int { return t.frames }
func (t *Task) Stop()       { t.state = TaskStopped }
func (t *Task) Start()      { t.Wait(0) }
func (t *Task) State() uint { return t.state }
func (t *Task) Remove()     { t.state = TaskCompleted }

//sorting primitives
type Tasks []*Task

func (t Tasks) Len() int      { return len(t) }
func (t Tasks) Swap(i, j int) { t[i], t[j] = t[j], t[i] }

type ByPriority struct{ Tasks }

func (t ByPriority) Less(i, j int) bool { return t.Tasks[i].priority < t.Tasks[j].priority }

//AddTask creates a new task and adds it to the queue
func AddTask(name string, function taskFunc, data interface{}, priority int, delay float64) {
	task := &Task{Name: name,
		Func:     function,
		Data:     data,
		start:    0,
		frames:   0,
		priority: priority,
		state:    TaskWaiting,
		delay:    0}
	if delay != 0 {
		task.Wait(delay)
	}
	taskList = append(taskList, task)

	//new task added; resort
	tasksSorted = false
}

//runTasks sorts the taskList by priority, then adds all active,
// non-waiting task to the task queue and processes them
func runTasks() {

	if !tasksSorted {
		sort.Sort(ByPriority{taskList})
		tasksSorted = true
	}

	for i, task := range taskList {
		switch task.state {
		case TaskCompleted:
			//remove task from list
			taskList = append(taskList[:i], taskList[i+1:]...)
		case TaskWaiting:
			//check delay
			if task.delay <= Time() {
				task.state = TaskRunning
			}
		case TaskRunning:
			taskQueue = append(taskQueue, task)
		case TaskStopped:
			//do nothing
		}

	}

	//run through all queued tasks
	for _, task := range taskQueue {
		if task.start == 0 {
			task.start = Time()
		}
		task.frames++
		task.Func(task)

	}

	//empty queue
	taskQueue = taskQueue[0:0]
}
