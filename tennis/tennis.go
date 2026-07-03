package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var wg sync.WaitGroup

func init() {
	source := rand.NewSource(time.Now().UnixNano())
	rand.New(source)
}

func main() {
	court := make(chan int)

	wg.Add(2)

	go player("Puu1", court)
	go player("Puu2", court)

	court <- 1

	wg.Wait()
}

func player(name string, court chan int) {
	defer wg.Done()

	for {
		ball, ok := <-court
		if !ok {
			fmt.Printf("Player %s won\n", name)
			return
		}

		n := rand.Intn(100)
		if n%13 == 0 {
			fmt.Printf("Player %s missed the ball\n", name)
			close(court)
			return
		}

		fmt.Printf("Player %s hit %d\n", name, ball)
		ball++
		court <- ball
	}
}
