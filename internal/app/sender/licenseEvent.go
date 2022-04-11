package sender

import (
	"github.com/ozonmp/lic-license-api/internal/model/license"
)

type LicenseEventSender interface {
	Send(license *license.LicenseEvent) error
}
