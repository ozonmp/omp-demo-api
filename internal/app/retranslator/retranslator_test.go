package retranslator

import (
	"errors"
	"fmt"
	"github.com/ozonmp/omp-demo-api/internal/model"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/ozonmp/omp-demo-api/internal/mocks"
)

func TestStart(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := mocks.NewMockEventRepo(ctrl)
	sender := mocks.NewMockEventSender(ctrl)

	consumeSize := uint64(2)

	card := model.Card{
		OwnerId: 6, Number: "1234567887654321", Cvv: "097", ExpirationDate: "1.01.2022", CardType: model.DEBIT,
	}

	expEvents := []model.CardEvent{
		{ID: 3, Type: model.Created, Status: model.Processed, Entity: &card},
		{ID:4, Type: model.Updated, Status: model.Deferred, Entity: &card},
	}

	repo.EXPECT().Lock(consumeSize).Return(expEvents, nil).AnyTimes()
	sender.EXPECT().Send(&expEvents[0]).Return(nil).AnyTimes()
	sender.EXPECT().Send(&expEvents[1]).Return(errors.New("Failed event\n")).AnyTimes()
	repo.EXPECT().Remove([]uint64{ expEvents[0].ID }).Return(nil).AnyTimes()
	repo.EXPECT().Unlock([]uint64{ expEvents[1].ID }).Return(nil).AnyTimes()

	cfg := Config{
		ChannelSize:   5,
		ConsumerCount: 2,
		ConsumeSize:  consumeSize,
		ConsumeTimeout: 100 * time.Millisecond,
		ProducerCount: 2,
		WorkerCount:   2,
		Repo:          repo,
		Sender:        sender,
	}

	fmt.Printf("Test started\n")
	ret := NewRetranslator(cfg)
	ret.Start()
	time.Sleep(400 * time.Millisecond)
	ret.Close()
}
