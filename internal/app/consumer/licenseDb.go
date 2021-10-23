package consumer

import (
	"github.com/ozonmp/omp-demo-api/internal/model/license"
	"sync"
	"time"

	"github.com/ozonmp/omp-demo-api/internal/app/repo"
	"github.com/ozonmp/omp-demo-api/internal/model"
)

type LicenseConsumer interface {
	Start()
	Close()
}

type licenseConsumer struct {
	n      uint64
	events chan<- model.SubdomainEvent

	repo repo.EventRepo

	batchSize uint64
	timeout   time.Duration

	done chan bool
	wg   *sync.WaitGroup
}

type LicenseConfig struct {
	n         uint64
	events    chan<- license.LicenseEvent
	repo      repo.EventRepo
	batchSize uint64
	timeout   time.Duration
}

func NewLicenseDbConsumer(
	n uint64,
	batchSize uint64,
	consumeTimeout time.Duration,
	repo repo.EventRepo,
	events chan<- model.SubdomainEvent) LicenseConsumer {

	wg := &sync.WaitGroup{}
	done := make(chan bool)

	return &licenseConsumer{
		n:         n,
		batchSize: batchSize,
		timeout:   consumeTimeout,
		repo:      repo,
		events:    events,
		wg:        wg,
		done:      done,
	}
}

func (c *licenseConsumer) Start() {
	for i := uint64(0); i < c.n; i++ {
		c.wg.Add(1)

		go func() {
			defer c.wg.Done()
			ticker := time.NewTicker(c.timeout)
			for {
				select {
				case <-ticker.C:
					events, err := c.repo.Lock(c.batchSize)
					if err != nil {
						continue
					}
					for _, event := range events {
						c.events <- event
					}
				case <-c.done:
					return
				}
			}
		}()
	}
}

func (c *licenseConsumer) Close() {
	close(c.done)
	c.wg.Wait()
}
