package engine

//task state
const (
	TaskStopped = iota
	TaskRunning
	TaskWaiting
	TaskCompleted
)

type Task func(event Event)

type Event struct {
	Name     string
	Time     float64
	Frames   int
	State    uint
	Priority uint
	Data     map[string]interface{}
}
