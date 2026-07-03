package worker

type Worker interface {
	Task()
}

type Pool struct {
	worker chan Worker
}

func New(wokerNum int) *Pool {
	p := Pool{worker: make(chan Worker)}

	for i := 0; i < wokerNum; i++ {
		go func() {
			for work := range p.worker {
				work.Task()
			}
		}()
	}
	return &p
}

func (p *Pool) Run(w Worker) {
	p.worker <- w
}
