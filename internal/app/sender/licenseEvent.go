package sender

import (
	"github.com/ozonmp/omp-demo-api/internal/model/license"
)

type LicenseEventSender interface {
	Send(license *license.LicenseEvent) error
}
