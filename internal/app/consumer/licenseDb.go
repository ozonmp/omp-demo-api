package consumer

import (
	"context"
	"github.com/ozonmp/lic-license-api/internal/app/repo"
	"github.com/ozonmp/lic-license-api/internal/model/license"
	"sync"
	"time"
)

type LicenseConsumer interface {
	Start(ctx context.Context)
	Close()
}

type licenseConsumer struct {
	n      uint64
	events chan<- license.LicenseEvent

	repo repo.LicenseEventRepo

	batchSize uint64
	timeout   time.Duration

	ctx context.Context
	wg  *sync.WaitGroup
}

type LicenseConfig struct {
	n         uint64
	events    chan<- license.LicenseEvent
	repo      repo.LicenseEventRepo
	batchSize uint64
	timeout   time.Duration
}

func NewLicenseDbConsumer(
	n uint64,
	batchSize uint64,
	consumeTimeout time.Duration,
	repo repo.LicenseEventRepo,
	events chan<- license.LicenseEvent) LicenseConsumer {

	wg := &sync.WaitGroup{}

	return &licenseConsumer{
		n:         n,
		batchSize: batchSize,
		timeout:   consumeTimeout,
		repo:      repo,
		events:    events,
		wg:        wg,
	}
}

func (c *licenseConsumer) Start(ctx context.Context) {
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
				case <-ctx.Done():
					return
				}
			}
		}()
	}
}

func (c *licenseConsumer) Close() {
	c.wg.Wait()
}
