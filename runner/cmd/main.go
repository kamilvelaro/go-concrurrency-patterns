package main

import (
	"log"
	"math/rand"
	"patterns/runner/runner"
	"patterns/runner/task/dummy"
	"sync"
	"time"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile | log.Lmsgprefix)
}

func main() {

	var wg sync.WaitGroup
	wg.Add(1)

	dummyTask := dummy.Dummy{}

	log.Println("Starting runner ... ")
	r := runner.New(2, 10)
	r.AddTask(dummyTask)
	r.Run(&wg)

	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	r.AddTask(dummyTask)

	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	r.AddTask(dummyTask)

	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	r.AddTask(dummyTask)

	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	r.AddTask(dummyTask)

	wg.Wait()

}
