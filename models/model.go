package models

import (
	logger "scootin/loger"
	"sync"
	"time"
)

// Scooter embodies the scooter functionality
type Scooter struct {
	Info *ScooterInfo
	Done chan bool
	mu   *sync.Mutex
}

// ScooterInfo has the scooter details
type ScooterInfo struct {
	ID           string
	Coordination int64 // represents the scooter location.
	UserID       string
}

// LocationUpdate contains the time, and geographical coordinates.
type LocationUpdate struct {
	ScooterID   string
	Time        time.Time
	Coordinates int64 // represents the scooter location update
}

type User struct {
	ID string
}

// Event represents the scooter's trip event, End is false for Start()
// and true for End()
type Event struct {
	ScooterID string
	UserID    string
	End       bool
}

// IsOccupied returns true if the scooter is being used, false otherwise.
func (s *Scooter) IsOccupied() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Info.UserID != ""
}

// Start reports an event when a trip begins
func (s *Scooter) Start(userID string) Event {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Info.UserID = userID
	go s.Updates() // periodic updates
	return Event{ScooterID: s.Info.ID, UserID: userID, End: false}
}

// End reports an event when a trip ends
func (s *Scooter) End() Event {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Info.UserID = "" // Not Occupied, make it available
	s.Done <- true     // signal end of the trip to the update() routine
	return Event{ScooterID: s.Info.ID, UserID: s.Info.UserID, End: true}
}

// Updates report scooter periodic updates
func (s *Scooter) Updates() {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-s.Done:
			return
		case <-ticker.C:
			s.mu.Lock()
			// redirect the update report to the log
			logger.Infof("%+v", LocationUpdate{ScooterID: s.Info.ID, Coordinates: s.Info.Coordination, Time: time.Now()})
			s.mu.Unlock()
		}
	}
}
