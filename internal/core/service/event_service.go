// internal/core/service/event_service.go
package service

import (
	"errors"
	"event-registration/internal/core/domain"
	"time"
)

// EventService handles the business logic for event registrations
type EventService struct {
	repo  domain.EventRepository
	cache domain.EventCache
	queue domain.EventQueue
}

// NewEventService creates a new instance of EventService
func NewEventService(repo domain.EventRepository, cache domain.EventCache, queue domain.EventQueue) *EventService {
	return &EventService{
		repo:  repo,
		cache: cache,
		queue: queue,
	}
}

// RegisterEvent registers a user for an event
func (s *EventService) RegisterEvent(eventID, userID string) error {
	// Check if the event exists in the cache
	cachedEvent, err := s.cache.Get(eventID)
	if err == nil && cachedEvent != nil {
		// Use the cached event
		return s.processRegistration(cachedEvent, userID)
	}

	// If not in cache, fetch from the repository
	event, err := s.repo.FindByID(eventID)
	if err != nil {
		return errors.New("event not found")
	}

	// Cache the event for future use
	s.cache.Set(eventID, event, time.Minute)

	// Process the registration
	return s.processRegistration(event, userID)
}

// processRegistration handles the actual registration logic
// TODO: FIX THIS
func (s *EventService) processRegistration(event *domain.Event, userID string) error {
	// Check if there are available slots
	if event.BookedSlots >= event.TotalSlots {
		return errors.New("no available slots")
	}

	// Update the booked slots
	event.BookedSlots++

	// Save the updated event to the repository
	err := s.repo.Update(event)
	if err != nil {
		return errors.New("failed to update event")
	}

	// Publish the event to the queue for further processing
	err = s.queue.Publish(event)
	if err != nil {
		return errors.New("failed to publish event")
	}

	return nil
}
