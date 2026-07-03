package db

import (
	"io"
	"log"
	"sync/atomic"
)

var idCounter int32

type DbConnection struct {
	ID int32
}

func (dbConn *DbConnection) Close() error {
	log.Printf("Closing connection %d", dbConn.ID)
	return nil
}

func CreateConnection() (io.Closer, error) {
	id := atomic.AddInt32(&idCounter, 1)
	log.Printf("New connection %d has been created", id)
	return &DbConnection{ID: id}, nil
}
