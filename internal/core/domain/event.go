package domain

import (
	"time"
)

type Event struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	TotalSlots  int       `json:"total_slots"`
	BookedSlots int       `json:"booked_slots"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type EventRepository interface {
	FindByID(id string) (*Event, error)
	Save(event *Event) error
	Update(event *Event) error
}

type EventCache interface {
	Get(key string) (*Event, error)
	Set(key string, event *Event, expiration time.Duration) error
}

type EventQueue interface {
	Publish(event *Event) error
	Consume() (<-chan *Event, error)
}
