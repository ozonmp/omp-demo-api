package main

import (
	"context"
	"github.com/ozonmp/lic-license-api/internal/app/retranslator"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	sigs := make(chan os.Signal, 1)

	cfg := retranslator.Config{
		ChannelSize:   512,
		ConsumerCount: 2,
		ConsumeSize:   10,
		ProducerCount: 28,
		WorkerCount:   2,
	}

	retranslator := retranslator.NewRetranslator(cfg)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	retranslator.Start(ctx)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
}
