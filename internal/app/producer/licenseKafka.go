package producer

import (
	"github.com/ozonmp/omp-demo-api/internal/app/repo"
	"github.com/ozonmp/omp-demo-api/internal/model/license"
	"sync"
	"time"

	"github.com/gammazero/workerpool"
	"github.com/ozonmp/omp-demo-api/internal/app/sender"
)

type LicenseProducer interface {
	Start()
	Close()
}

type licenseProducer struct {
	n       uint64
	timeout time.Duration

	sender sender.LicenseEventSender
	repo   repo.LicenseEventRepo
	events <-chan license.LicenseEvent

	workerPool *workerpool.WorkerPool

	wg   *sync.WaitGroup
	done chan bool
}

// [x] todo for students: add repo (ADDED)
func NewKafkaLicenseProducer(
	n uint64,
	sender sender.LicenseEventSender,
	repo repo.LicenseEventRepo,
	events <-chan license.LicenseEvent,
	workerPool *workerpool.WorkerPool,
) LicenseProducer {

	wg := &sync.WaitGroup{}
	done := make(chan bool)

	return &licenseProducer{
		n:          n,
		sender:     sender,
		repo:       repo,
		events:     events,
		workerPool: workerPool,
		wg:         wg,
		done:       done,
	}
}

func (p *licenseProducer) Start() {
	for i := uint64(0); i < p.n; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for {
				select {
				case event := <-p.events:
					if err := p.sender.Send(&event); err != nil {
						p.workerPool.Submit(func() {
							// ...
						})
					} else {
						p.workerPool.Submit(func() {
							// ...
						})
					}
				case <-p.done:
					return
				}
			}
		}()
	}
}

func (p *licenseProducer) Close() {
	close(p.done)
	p.wg.Wait()
}
