package consumer

import (
	"fmt"
	"sync"
	"time"

	"github.com/ozonmp/omp-demo-api/internal/app/repo"
	"github.com/ozonmp/omp-demo-api/internal/model"
)

type Consumer interface {
	Start()
	Close()
}

type consumer struct {
	n      uint64
	events chan<- model.CardEvent

	repo repo.EventRepo

	batchSize uint64
	timeout   time.Duration

	done chan bool
	wg   *sync.WaitGroup
}

type Config struct {
	n         uint64
	events    chan<- model.CardEvent
	repo      repo.EventRepo
	batchSize uint64
	timeout   time.Duration
}

func NewDbConsumer(
	n uint64,
	batchSize uint64,
	consumeTimeout time.Duration,
	repo repo.EventRepo,
	events chan<- model.CardEvent) Consumer {

	wg := &sync.WaitGroup{}
	done := make(chan bool)

	return &consumer{
		n:         n,
		batchSize: batchSize,
		timeout:   consumeTimeout,
		repo:      repo,
		events:    events,
		wg:        wg,
		done:      done,
	}
}

func (c *consumer) Start() {
	fmt.Printf("Consumer started\n")
	c.wg.Add(int(c.n))
	for i := uint64(0); i < c.n; i++ {

		go func() {
			defer c.wg.Done()
			ticker := time.NewTicker(c.timeout)
			for {
				select {
				case <-ticker.C:
					fmt.Printf("Lock call with %d\n", c.batchSize)
					events, err := c.repo.Lock(c.batchSize)
					if err != nil {
						fmt.Printf("Error: %s\n", err)
						continue
					}
					for _, event := range events {
						fmt.Printf("Consumer pushed event into channel: %s\n", event.String())
						c.events <- event
					}
				case <-c.done:
					fmt.Printf("Done call\n")
					return
				}
			}
		}()
	}

	go func() {
		c.wg.Wait()
	}()
}

func (c *consumer) Close() {
	close(c.done)
	//c.wg.Wait()
}
