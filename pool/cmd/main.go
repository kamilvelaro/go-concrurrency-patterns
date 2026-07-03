package main

import (
	"log"
	"math/rand"
	"patterns/pool/db"
	"patterns/pool/pool"
	"sync"
	"time"
)

const (
	maxGoroutines   = 10
	pooledResources = 2
)

var p pool.Pool
var wg sync.WaitGroup

func main() {
	p, err := pool.New(db.CreateConnection, pooledResources)
	if err != nil {
		panic(err)
	}

	wg.Add(maxGoroutines)
	for query := 0; query < maxGoroutines; query++ {
		go func(q int) {
			performQuery(q, p)
			wg.Done()
		}(query)
	}
	wg.Wait()
}

func performQuery(query int, pool *pool.Pool) {
	conn, err := pool.Acquire()
	if err != nil {
		log.Println(err)
	} else {
		defer pool.Release(conn)
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		log.Printf("Query ID: %d, Connection ID: %d", query, conn.(*db.DbConnection).ID)
	}
}
