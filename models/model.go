package models

import (
	"context"
	"scootin/db"
	"scootin/logger"
	"sync"
	"time"
)

// Scooter embodies the scooter functionality i.e, scooter runtime instance.
type Scooter struct {
	Info ScooterInfo
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

// User represents the user details
type User struct {
	ID    string
	Name  string
	Email string
}

// UUIDResponse used to return the uuid for the new user's and scooter's creation
type UUIDResponse struct {
	ID string
}

// Event represents the scooter's trip event, End is false for Start()
// and true for End()
type Event struct {
	ScooterID string
	UserID    string
	End       bool
}

// NewScooter returns a new scooter runtime instance
func NewScooter(ID string) *Scooter {
	return &Scooter{Info: ScooterInfo{
		ID: ID,
	},
		Done: make(chan bool),
		mu:   &sync.Mutex{}}

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

			// Persist the scooter coordination in the database
			if err := db.UpdateScooterCoordinates(context.Background(), s.Info.ID, s.Info.Coordination); err != nil {
				logger.Errorf("couldn't persist the scooter %s updates: %s", s.Info.ID, err)
			}
			s.mu.Unlock()
		}
	}
}
