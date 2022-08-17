// Package trip is for all the runtime instances and their functionality.
package client

import (
	"context"
	"math/rand"
	"scootin/db"
	"scootin/logger"
	"scootin/models"
	"sync"
	"time"
)

// Scooter embodies the scooter functionality i.e, scooter runtime instance.
type Scooter struct {
	Info models.ScooterInfo
	Done chan bool
	mu   *sync.Mutex
}

// LocationUpdate contains the time, and geographical coordinates.
type LocationUpdate struct {
	ScooterID   string
	Time        time.Time
	Coordinates int64 // represents the scooter location update
}

// NewScooter returns a new scooter runtime instance
func NewScooter(ID string) *Scooter {
	return &Scooter{Info: models.ScooterInfo{
		ID:           ID,
		UserID:       models.NotOccupied,
		Coordination: 1,
	},
		Done: make(chan bool),
		mu:   &sync.Mutex{}}

}

// Start reports an event when a trip begins
func (s *Scooter) Start(ctx context.Context, userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Info.UserID = userID
	if err := db.BookScooter(ctx, s.Info.ID, s.Info.UserID); err != nil {
		return err
	}
	go s.Updates(ctx) // periodic updates
	return nil
}

// End reports an event when a trip ends
func (s *Scooter) End(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := db.ReleaseScooter(ctx, s.Info.UserID); err != nil {
		return err
	}
	s.Done <- true // signal end of the trip to the update() routine
	return nil
}

// Updates moves and reports  the scooter periodic updates
func (s *Scooter) Updates(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-s.Done:
			return
		case <-ticker.C:
			s.mu.Lock()
			s.Info.Coordination += randomDistance()
			// redirect the update report to the log
			logger.Infof("%+v", LocationUpdate{ScooterID: s.Info.ID, Coordinates: s.Info.Coordination, Time: time.Now()})

			// Persist the scooter coordination in the database
			if err := db.UpdateScooterCoordinates(ctx, s.Info.ID, s.Info.Coordination); err != nil {
				logger.Errorf("couldn't persist the scooter %s updates: %s", s.Info.ID, err)
			}
			s.mu.Unlock()
		}
	}
}

// randomDistance returns random distance
func randomDistance() int64 {
	max := 20 // max scooter speed
	min := 7  // min scooter speed

	// choose random speed in-between
	return int64(rand.Intn(max-min+1) + min)
}
