package license

type License struct {
	ID    uint64
	Title string
}

type EventType uint8

type EventStatus uint8

const (
	Created EventType = iota
	Updated
	Removed

	Deferred EventStatus = iota
	Processed
)

type LicenseEvent struct {
	ID     uint64
	Type   EventType
	Status EventStatus
	Entity *License
}
