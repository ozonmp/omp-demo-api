package workerpool

import (
	"github.com/ozonmp/lic-license-api/internal/app/repo"
	"github.com/ozonmp/lic-license-api/internal/mocks"
	"github.com/ozonmp/lic-license-api/internal/model/license"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestWorkerPool(t *testing.T) {
	workersCount := 5
	var repo repo.LicenseEventRepo
	NewWorkerPool(workersCount, repo)
}

func TestCreateRetranslatorWorkerPool(t *testing.T) {
	workersCount := 5
	var repo repo.LicenseEventRepo

	wpI := NewWorkerPool(workersCount, repo)
	_, ok := wpI.(RetranslatorWorkerPool)
	assert.True(t, ok)
}

type testSubmitter struct {
	submitIsCalled bool
	task           func()
}

func (s *testSubmitter) Submit(task func()) {
	s.submitIsCalled = true
	s.task()
}

func (s *testSubmitter) StopWait() {}
func TestCleanCreatedType(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockLicenseEventRepo(ctrl)

	event := license.LicenseEvent{
		ID:     1,
		Type:   license.Created,
		Status: license.Processed,
		Entity: &license.License{
			ID: 1,
		},
	}
	repo.EXPECT().Remove([]uint64{event.ID}).Return(nil).Times(1)

	wp := RetranslatorWorkerPool{
		repo: repo,
	}
	task := func() {
		wp.clean(event)
	}
	s := &testSubmitter{
		task: task,
	}
	wp.submitter = s

	err := wp.Clean(event)
	assert.Nil(t, err)
	assert.True(t, s.submitIsCalled)
}

func TestCleanNotCreatedType(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockLicenseEventRepo(ctrl)

	checkCleanNotCreatedType(repo, license.LicenseEvent{
		ID:     1,
		Type:   license.Updated,
		Status: license.Processed,
		Entity: &license.License{
			ID: 1,
		},
	}, t)
	checkCleanNotCreatedType(repo, license.LicenseEvent{
		ID:     1,
		Type:   license.Removed,
		Status: license.Processed,
		Entity: &license.License{
			ID: 1,
		},
	}, t)
}

func checkCleanNotCreatedType(repo *mocks.MockLicenseEventRepo, event license.LicenseEvent, t *testing.T) {
	repo.EXPECT().Remove([]uint64{event.ID}).Return(nil).Times(0)

	wp := RetranslatorWorkerPool{
		submitter: nil,
		repo:      repo,
	}
	task := func() {
		wp.clean(event)
	}
	s := &testSubmitter{
		task: task,
	}
	wp.submitter = s

	err := wp.Clean(event)
	assert.Nil(t, err)
	assert.True(t, s.submitIsCalled)
}

func TestUpdateCreatedType(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockLicenseEventRepo(ctrl)

	event := license.LicenseEvent{
		ID:     1,
		Type:   license.Created,
		Status: license.Processed,
		Entity: &license.License{
			ID: 1,
		},
	}
	repo.EXPECT().Unlock([]uint64{event.ID}).Return(nil).Times(1)

	wp := RetranslatorWorkerPool{
		repo: repo,
	}
	task := func() {
		wp.update(event)
	}
	s := &testSubmitter{
		submitIsCalled: false,
		task:           task,
	}
	wp.submitter = s

	err := wp.Update(event)
	assert.Nil(t, err)
	assert.True(t, s.submitIsCalled)
}

func TestUpdateNotCreatedType(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mocks.NewMockLicenseEventRepo(ctrl)

	checkUpdateNotCreatedType(repo, license.LicenseEvent{
		ID:     1,
		Type:   license.Updated,
		Status: license.Processed,
		Entity: &license.License{
			ID: 1,
		},
	}, t)
	checkUpdateNotCreatedType(repo, license.LicenseEvent{
		ID:     1,
		Type:   license.Removed,
		Status: license.Processed,
		Entity: &license.License{
			ID: 1,
		},
	}, t)
}

func checkUpdateNotCreatedType(repo *mocks.MockLicenseEventRepo, event license.LicenseEvent, t *testing.T) {
	repo.EXPECT().Unlock([]uint64{event.ID}).Return(nil).Times(0)

	wp := RetranslatorWorkerPool{
		repo: repo,
	}
	task := func() {
		wp.update(event)
	}
	s := &testSubmitter{
		submitIsCalled: false,
		task:           task,
	}
	wp.submitter = s

	err := wp.Update(event)
	assert.Nil(t, err)
	assert.True(t, s.submitIsCalled)
}
