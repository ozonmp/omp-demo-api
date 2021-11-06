package model

type EventType uint8

type EventStatus uint8

const (
	Created EventType = iota
	Updated
	Removed

	Deferred EventStatus = iota
	Processed
)

func EventTypeToStr(eventType EventType) string {
	switch eventType {
	case Created: return "Created"
	case Updated: return "Updated"
	case Removed: return "Removed"
	}
	return "Undef event type"
}

func EventStatusToStr(eventStatus EventStatus) string {
	switch eventStatus {
	case Deferred: return "Deferred"
	case Processed: return "Processed"
	}
	return "Undef event status"
}

type EnumCardType uint8

const (
	DEBIT  EnumCardType = iota
	CREDIT
	UNDEF
)

func FromStrToEnum(in string) EnumCardType {
	switch in {
	case "DEBIT": return DEBIT
	case "CREDIT": return CREDIT
	default:
		return UNDEF
	}
}

func FromEnumToStr(in EnumCardType) string {
	switch in {
	case DEBIT: return "DEBIT"
	case CREDIT: return "CREDIT"
	default:
		return  "UNDEF TYPE"
	}
}

type Card struct {
	OwnerId uint64
	Number         string
	Cvv            string
	ExpirationDate string
	CardType       EnumCardType
}

func (p *Card) String() string {
	return FromEnumToStr(p.CardType) + " Card " + p.Number + " expiring " + p.ExpirationDate
}

type CardEvent struct {
	ID     uint64
	Type   EventType
	Status EventStatus
	Entity *Card
}

func (p *CardEvent) String() string {
	return p.Entity.String() + " " + EventTypeToStr(p.Type) + " and " + EventStatusToStr(p.Status)
}