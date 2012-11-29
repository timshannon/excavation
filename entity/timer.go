package entity

import (
	"errors"
	"excavation/engine"
	"strconv"
	"strings"
)

//Triggers a list of entities passed in as the following format
// entityName,delay,entityName,delay,entityName,delay
type Timer struct {
	node     *engine.Node
	triggers map[float64]Entity
}

func (t *Timer) Add(node *engine.Node, args EntityArgs) {
	t.node = node
	t.triggers = make(map[float64]Entity)
	triggerList := strings.Split(args.String("triggers"), ",")

	for i := 0; i < len(triggerList); i += 2 {
		f, err := strconv.ParseFloat(triggerList[i+1], 64)
		if err != nil {
			engine.RaiseError(errors.New("Invalid delay for timer: " + t.node.Name() +
				" and trigger: " + triggerList[i]))
		}
		if trigger, ok := EntityFromName(triggerList[i]); ok {
			t.triggers[f] = trigger

		} else {
			engine.RaiseError(errors.New("Entity Name: " + triggerList[i] + " not found for timer."))
		}
	}
}

func (t *Timer) Trigger(value float32) {
	if value > 0 {
		for k, v := range t.triggers {
			engine.AddTask(t.node.Name()+"_TimerItem", triggerTask, v, 0, k)
		}

	}
}

func triggerTask(t *engine.Task) {
	t.Data.(Entity).Trigger(1)
	t.Remove()
}
