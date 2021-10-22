package retranslator

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/ozonmp/omp-demo-api/internal/mocks"
)

func TestStart(t *testing.T) {

	ctrl := gomock.NewController(t)
	repo := mocks.NewMockEventRepo(ctrl)
	sender := mocks.NewMockEventSender(ctrl)

	repo.EXPECT().Lock(gomock.Any()).AnyTimes()

	cfg := Config{
		ChannelSize:   512,
		ConsumerCount: 2,
		ConsumeSize:   10,
		ConsumeTimeout: 10 * time.Second,
		ProducerCount: 2,
		WorkerCount:   2,
		Repo:          repo,
		Sender:        sender,
	}

	retranslator := NewRetranslator(cfg)
	retranslator.Start()
	retranslator.Close()
}
