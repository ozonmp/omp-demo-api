package retranslator

import (
	"context"
	"github.com/ozonmp/lic-license-api/internal/app/consumer"
	"github.com/ozonmp/lic-license-api/internal/app/producer"
	"github.com/ozonmp/lic-license-api/internal/app/repo"
	"github.com/ozonmp/lic-license-api/internal/app/sender"
	"github.com/ozonmp/lic-license-api/internal/model/license"
	"time"

	"github.com/gammazero/workerpool"
)

type Retranslator interface {
	Start(ctx context.Context)
	Close()
}

type Config struct {
	ChannelSize uint64

	ConsumerCount  uint64
	ConsumeSize    uint64
	ConsumeTimeout time.Duration

	ProducerCount uint64
	WorkerCount   int

	Repo   repo.LicenseEventRepo
	Sender sender.LicenseEventSender
}

type retranslator struct {
	events   chan license.LicenseEvent
	consumer consumer.LicenseConsumer
	producer producer.LicenseProducer
	cancel   context.CancelFunc
}

func NewRetranslator(cfg Config) Retranslator {
	events := make(chan license.LicenseEvent, cfg.ChannelSize)

	consumer := consumer.NewLicenseDbConsumer(
		cfg.ConsumerCount,
		cfg.ConsumeSize,
		cfg.ConsumeTimeout,
		cfg.Repo,
		events)
	producer := producer.NewKafkaLicenseProducer(
		cfg.ProducerCount,
		cfg.Sender,
		events,
		workerpool.NewWorkerPool(cfg.WorkerCount, cfg.Repo))

	return &retranslator{
		events:   events,
		consumer: consumer,
		producer: producer,
	}
}

func (r *retranslator) Start(ctx context.Context) {
	r.producer.Start(ctx)
	r.consumer.Start(ctx)
}

func (r *retranslator) Close() {
	r.consumer.Close()
	r.producer.Close()
}
