package model

type Keyword struct {
    ID uint64
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

type KeywordEvent struct {
    ID     uint64
    Type   EventType
    Status EventStatus
    Entity *Keyword
}
