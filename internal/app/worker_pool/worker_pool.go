package workerpool

import (
	"github.com/gammazero/workerpool"
	"github.com/ozonmp/lic-license-api/internal/app/repo"
	"github.com/ozonmp/lic-license-api/internal/model/license"
)

type WorkerPool interface {
	Clean(event license.LicenseEvent) error
	Update(event license.LicenseEvent) error
	Stop()
}

type submitter interface {
	Submit(task func())
	StopWait()
}

func NewWorkerPool(workerCount int, repo repo.LicenseEventRepo) WorkerPool {
	return RetranslatorWorkerPool{
		submitter: workerpool.New(workerCount),
		repo:      repo,
	}
}

type RetranslatorWorkerPool struct {
	submitter submitter
	repo      repo.LicenseEventRepo
}

func (wp RetranslatorWorkerPool) Clean(event license.LicenseEvent) error {
	wp.submitter.Submit(func() {
		wp.clean(event)
	})
	return nil
}

func (wp RetranslatorWorkerPool) clean(event license.LicenseEvent) {
	if event.Type == license.Created {
		wp.repo.Remove([]uint64{event.ID})
	}
}

func (wp RetranslatorWorkerPool) Update(event license.LicenseEvent) error {
	wp.submitter.Submit(func() {
		wp.update(event)
	})
	return nil
}

func (wp RetranslatorWorkerPool) update(event license.LicenseEvent) {
	if event.Type == license.Created {
		wp.repo.Unlock([]uint64{event.ID})
	}
}

func (wp RetranslatorWorkerPool) Stop() {
	wp.submitter.StopWait()
}
