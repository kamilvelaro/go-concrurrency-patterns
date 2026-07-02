package runner

import (
	"log"
	"os"
	"os/signal"
	"patterns/runner/task"
	"sync"
	"time"
)

type Runner struct {
	interrupt  chan os.Signal
	timeout    <-chan time.Time
	workersNum uint
	tasks      chan task.Task
	completed  chan task.Completed
}

func New(workersNum uint, tasksNum int) *Runner {
	return &Runner{
		interrupt:  make(chan os.Signal),
		timeout:    make(<-chan time.Time),
		workersNum: workersNum,
		tasks:      make(chan task.Task, tasksNum),
		completed:  make(chan task.Completed),
	}
}

func (r *Runner) AddTask(tasks ...task.Task) {
	for _, task := range tasks {
		r.tasks <- task
	}
}

func (r *Runner) addRoutine() {
	log.Println("Adding new routine")
	go func() {
		var task task.Task
		task, ok := <-r.tasks
		log.Println("===")
		if !ok {
			log.Println("Failed to run task, because task pipe is closed")
			return
		}
		r.completed <- task.Run()
		log.Println("Finishing routine")
	}()
}

func (r *Runner) Run(wg *sync.WaitGroup) error {
	signal.Notify(r.interrupt, os.Interrupt)

	for i := 0; i < int(r.workersNum); i++ {
		r.addRoutine()
	}

	go func() {
	runnerLoop:
		for {
			select {
			case <-r.interrupt:
				log.Println("Stopping runner")
				wg.Done()
				break runnerLoop
			case completed := <-r.completed:
				if completed.Err != nil {
					log.Printf("Task finished with error: %v", completed.Err.Error())
				} else {
					log.Println("Task finished with success")
				}
				r.addRoutine()
			}
		}
	}()

	return nil
}
