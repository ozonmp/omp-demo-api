package producer

import (
	"fmt"
	"github.com/ozonmp/omp-demo-api/internal/app/repo"
	"sync"
	"time"

	"github.com/ozonmp/omp-demo-api/internal/app/sender"
	"github.com/ozonmp/omp-demo-api/internal/model"

	"github.com/gammazero/workerpool"
)

type Producer interface {
	Start()
	Close()
}

type producer struct {
	n       uint64
	timeout time.Duration

	sender sender.EventSender
	events <-chan model.CardEvent

	workerPool *workerpool.WorkerPool

	wg   *sync.WaitGroup
	done chan bool

	repo repo.EventRepo
}

// todo for students: add repo
func NewKafkaProducer(
	n uint64,
	sender sender.EventSender,
	events <-chan model.CardEvent,
	workerPool *workerpool.WorkerPool,
	repo repo.EventRepo,
) Producer {

	wg := &sync.WaitGroup{}
	done := make(chan bool)

	return &producer{
		n:          n,
		sender:     sender,
		events:     events,
		workerPool: workerPool,
		wg:         wg,
		done:       done,
		repo:       repo,
	}
}

func (p *producer) Start() {
	fmt.Printf("Producer started\n")
	p.wg.Add(int(p.n))
	for i := uint64(0); i < p.n; i++ {
		go func() {
			defer p.wg.Done()
			for {
				select {
				case event := <-p.events:
					fmt.Printf("Producer read event %s from event channel\n", event.String())
					if err := p.sender.Send(&event); err != nil {
						fmt.Printf("Failure send for %s\n", event.String())
						p.workerPool.Submit(func() {
							// ...
							p.repo.Unlock([]uint64{ event.ID })
						})
					} else {
						fmt.Printf("Success send for %s\n", event.String())
						p.workerPool.Submit(func() {
							// ...
							p.repo.Remove([]uint64{ event.ID })
						})
					}
				case <-p.done:
					return
				}
			}
		}()
	}
	go func() {
		p.wg.Wait()
	}()
}

func (p *producer) Close() {
	close(p.done)
}
