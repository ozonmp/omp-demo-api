package repo

import (
	"github.com/ozonmp/lic-license-api/internal/model/license"
)

// TODO: Think about is it Event?
type LicenseEventRepo interface {
	Lock(n uint64) ([]license.LicenseEvent, error)
	Unlock(eventIDs []uint64) error

	Add(event []license.LicenseEvent) error // TODO: should trigger Created License Event?
	Remove(eventIDs []uint64) error
}
