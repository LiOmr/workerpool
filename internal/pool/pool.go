package pool

import (
	"context"
	"fmt"
	"log"
	"sync"
)

type Pool struct {
	ctx    context.Context
	cancel context.CancelFunc

	jobs chan string

	addWorker    chan struct{}
	removeWorker chan struct{}

	mu      sync.Mutex
	workers map[int]chan struct{}
	nextID  int

	wg sync.WaitGroup
}

func New(buffer int) *Pool {
	ctx, cancel := context.WithCancel(context.Background())
	p := &Pool{
		ctx:          ctx,
		cancel:       cancel,
		jobs:         make(chan string, buffer),
		addWorker:    make(chan struct{}),
		removeWorker: make(chan struct{}),
		workers:      make(map[int]chan struct{}),
	}
	go p.dispatch()
	return p
}

func (p *Pool) dispatch() {
	for {
		select {
		case <-p.ctx.Done():
			p.stopAll()
			return
		case <-p.addWorker:
			p.startWorker()
		case <-p.removeWorker:
			p.stopOne()
		}
	}
}

func (p *Pool) Submit(s string) {
	select {
	case p.jobs <- s:
	case <-p.ctx.Done():
		log.Println("pool closed, task dropped")
	}
}

func (p *Pool) AddWorker() { p.addWorker <- struct{}{} }

func (p *Pool) RemoveWorker() { p.removeWorker <- struct{}{} }

func (p *Pool) Shutdown() {
	p.cancel()
	p.wg.Wait()
}

func (p *Pool) startWorker() {
	p.mu.Lock()
	p.nextID++
	id := p.nextID
	stop := make(chan struct{})
	p.workers[id] = stop
	p.mu.Unlock()

	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		for {
			select {
			case <-p.ctx.Done():
				return
			case <-stop:
				return
			case job := <-p.jobs:
				fmt.Printf("Worker %d: %s\n", id, job)
			}
		}
	}()
}

func (p *Pool) stopOne() {
	p.mu.Lock()
	for id, stop := range p.workers {
		close(stop)
		delete(p.workers, id)
		break
	}
	p.mu.Unlock()
}

func (p *Pool) stopAll() {
	p.mu.Lock()
	for id, stop := range p.workers {
		close(stop)
		delete(p.workers, id)
	}
	p.mu.Unlock()
}
