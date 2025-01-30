package gorm

import (
	"event-registration/internal/core/domain"

	"gorm.io/gorm"
)

type EventRepo struct {
	db *gorm.DB
}

func NewEventRepo(db *gorm.DB) *EventRepo {
	return &EventRepo{db: db}
}

func (r *EventRepo) FindByID(id string) (*domain.Event, error) {
	var event domain.Event
	err := r.db.First(&event, "id = ?", id).Error
	return &event, err
}

func (r *EventRepo) Save(event *domain.Event) error {
	return r.db.Create(event).Error
}

func (r *EventRepo) Update(event *domain.Event) error {
	return r.db.Save(event).Error
}
