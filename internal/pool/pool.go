package pool

import (
	"sync"
)

type Pool[R any] struct {
	workerCount int

	workList  []work
	inputChan chan work

	wg sync.WaitGroup
}

type work func()

func NewPool[R any](
	workerCount int,
) *Pool[R] {
	return &Pool[R]{
		workerCount: workerCount,
		wg:          sync.WaitGroup{},
	}
}

func (p *Pool[R]) AddWork(w work) {
	p.workList = append(p.workList, w)
}

func (p *Pool[R]) start() {
	for _, v := range p.workList {
		p.inputChan <- v
	}
	close(p.inputChan)
}

func (p *Pool[R]) Start() {
	p.wg.Add(len(p.workList))
	ic := make(chan work, len(p.workList))
	p.inputChan = ic
	for range p.workerCount {
		go func() {
			for v := range p.inputChan {
				v()
				p.wg.Done()
			}
		}()
	}
	p.start()
	p.wg.Wait()
	p.workList = []work{}
}
