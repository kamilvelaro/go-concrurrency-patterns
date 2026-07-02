package dummy

import (
	"log"
	"patterns/runner/task"
)

type Dummy struct {
}

func (d Dummy) Run(args ...interface{}) task.Completed {
	log.Println("Dummy task")
	return task.Completed{nil}
}
