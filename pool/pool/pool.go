package pool

import (
	"errors"
	"fmt"
	"io"
	"log"
	"sync"
)

type Pool struct {
	sync.Mutex
	resources       chan io.Closer
	factory         func() (io.Closer, error)
	poolCapacity    uint
	connectionsMade uint
	closed          bool
}

func New(fn func() (io.Closer, error), size uint) (*Pool, error) {
	return &Pool{resources: make(chan io.Closer, size), factory: fn, poolCapacity: size}, nil
}

func (p *Pool) Acquire() (io.Closer, error) {
	select {
	case r, ok := <-p.resources:
		if !ok {
			log.Println("Pool is closed, no resources available")
		}
		log.Println("Giving resource from pool")
		return r, nil
	default:
		if p.connectionsMade < p.poolCapacity {
			log.Println("Acquired new resource")
			p.Lock()
			p.connectionsMade++
			p.Unlock()
			return p.factory()
		} else {
			return nil, errors.New("All resources in use")
		}

	}
}

func (p *Pool) Release(r io.Closer) {
	if p.closed {
		fmt.Println("Pool is closed, release is impossible")
		r.Close()
		return
	}

	p.Lock()
	p.connectionsMade--
	p.resources <- r
	p.Unlock()
	log.Println("Released")

}

func (p *Pool) Close() {
	if p.closed {
		return
	}
	p.closed = true

	close(p.resources)

	// for r := range p.resources {
	// 	close(r)
	// }
}
