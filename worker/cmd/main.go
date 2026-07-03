package main

import (
	"log"
	"patterns/worker/worker"
	"sync"
	"time"
)

var names = []string{"steve", "bob", "mary"}

type namePrinter struct {
	name string
}

func (m *namePrinter) Task() {
	log.Println(m.name)
	time.Sleep(100 * time.Millisecond)
}

var wg sync.WaitGroup

func main() {

	p := worker.New(2)

	wg.Add(10 * len(names))

	for i := 0; i < 10; i++ {
		for _, name := range names {
			np := namePrinter{name: name}

			go func() {
				p.Run(&np)
				wg.Done()
			}()
		}

	}
	wg.Wait()
}
