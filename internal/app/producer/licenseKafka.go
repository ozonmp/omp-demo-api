package producer

import (
	"context"
	"github.com/gammazero/workerpool"
	"github.com/ozonmp/lic-license-api/internal/app/repo"
	"github.com/ozonmp/lic-license-api/internal/app/sender"
	"github.com/ozonmp/lic-license-api/internal/model/license"
	"sync"
	"time"
)

type LicenseProducer interface {
	Start(ctx context.Context)
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
	events <-chan license.LicenseEvent,
	workerPool *workerpool.WorkerPool,
) LicenseProducer {

	wg := &sync.WaitGroup{}
	done := make(chan bool)

	return &licenseProducer{
		n:          n,
		sender:     sender,
		events:     events,
		workerPool: workerPool,
		wg:         wg,
		done:       done,
	}
}

func (p *licenseProducer) Start(ctx context.Context) {
	for i := uint64(0); i < p.n; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for {
				select {
				case event := <-p.events:
					if event.Type == license.Created {
						if err := p.sender.Send(&event); err != nil {
							p.workerPool.Update(event)
						} else {
							p.workerPool.Clean(event)
						}
					}
				case <-ctx.Done():
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
