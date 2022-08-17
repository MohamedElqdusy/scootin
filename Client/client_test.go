package client

import (
	"context"
	"os"
	"scootin/db"
	"scootin/logger"
	"scootin/models"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	var (
		err error
		uid *models.UUIDResponse
	)

	// create the client
	baseUrl := "http://localhost:8080"
	c := NewClient(baseUrl)

	////////////////////  create Users  //////////////////////
	u1 := &models.User{Name: "David", Email: "david@scootin.com"}
	uid, err = c.CreateUser(u1)
	assert.NoError(t, err)
	u1.ID = uid.ID
	assert.True(t, isValidUUID(u1.ID))

	u2 := &models.User{Name: "Dan", Email: "dan@scootin.com"}
	uid, err = c.CreateUser(u2)
	assert.NoError(t, err)
	u2.ID = uid.ID
	assert.True(t, isValidUUID(u2.ID))

	u3 := &models.User{Name: "Sam", Email: "sam@scootin.com"}
	uid, err = c.CreateUser(u3)
	assert.NoError(t, err)
	u3.ID = uid.ID
	assert.True(t, isValidUUID(u3.ID))
	////////////////////  create the scooters  //////////////////////
	scooterIDs := createScootersIDs(t, c)

	////////////////////  ListAvailableScooter  //////////////////////
	// checks all available scooters
	scs, err := c.ListAvailableScooter()
	assert.NoError(t, err)

	// converts the available scooters to map for easy checking
	expectedScootersMap := scooterInfoSliceToMap(scs)
	assert.Equal(t, len(expectedScootersMap), len(scooterIDs))
	// compare the created scooter against the available in the database
	for _, scooterID := range scooterIDs {
		assert.True(t, expectedScootersMap[scooterID])
	}
	////////////////////////////   BookScooter  ///////////////////////////
	// book the scooter sc1 by the user u1
	err = c.BookScooter(scooterIDs[0], u1.ID)
	assert.NoError(t, err)

	// we should receive an error if we tried to book the same scooter by another user
	err = c.BookScooter(scooterIDs[0], u2.ID)
	assert.Error(t, err)

	// book the scooter sc2 by the user u2
	err = c.BookScooter(scooterIDs[1], u2.ID)
	assert.NoError(t, err)

	// checks all available scooters, we booked 2, so we have 1 left available
	scs, err = c.ListAvailableScooter()
	assert.NoError(t, err)
	assert.Len(t, scs, 1)
	assert.Equal(t, scs[0].UserID, models.NotOccupied)

	////////////////////////////   ReleaseScooter  ///////////////////////////
	err = c.ReleaseScooter(u1.ID)
	assert.NoError(t, err)
	// checks all available scooters, we have 1 booked, so we have 2 left available
	scs, err = c.ListAvailableScooter()
	assert.NoError(t, err)
	assert.Len(t, scs, 2)
	assert.Equal(t, scs[0].UserID, models.NotOccupied)
	assert.Equal(t, scs[1].UserID, models.NotOccupied)

	err = c.ReleaseScooter(u2.ID)
	assert.NoError(t, err)
	// checks all available scooters, we have 0 booked, so we have 3 left available
	scs, err = c.ListAvailableScooter()
	assert.NoError(t, err)
	assert.Len(t, scs, 3)
	assert.Equal(t, scs[0].UserID, models.NotOccupied)
	assert.Equal(t, scs[1].UserID, models.NotOccupied)
	assert.Equal(t, scs[2].UserID, models.NotOccupied)

	////////////////////////////   Scooters' Simulation  ///////////////////////////
	// initialize postgres connection and logging for the simulation
	setupEnv(t)
	err = db.InitiatePostgre()
	assert.NoError(t, err)

	log := logger.NewLogger()
	logger.InitLogger(log)
	defer logger.Sync()

	// create runtime scooters
	scooters := make([]*Scooter, 0)
	for _, x := range scs {
		scooters = append(scooters, NewScooter(x.ID))
	}

	ctx := context.Background()

	// book the first available scooter by the user u1
	err = scooters[0].Start(ctx, u1.ID)
	assert.NoError(t, err)

	// book the second available scooter by the user u2
	err = scooters[1].Start(ctx, u2.ID)
	assert.NoError(t, err)

	// book the third available scooter by the user u3
	err = scooters[2].Start(ctx, u3.ID)
	assert.NoError(t, err)

	// checks all available scooters, all of them are booked, so we have 0 left available
	scs, err = c.ListAvailableScooter()
	assert.NoError(t, err)
	assert.Len(t, scs, 0)

	// end the scooters' trips  after 10s
	ticker := time.NewTicker(10 * time.Second)
	endTrips(t, ctx, ticker, scooters)

	// checks all available scooters, none of them are booked, so we have 3 left available
	scs, err = c.ListAvailableScooter()
	assert.NoError(t, err)
	assert.Len(t, scs, 3)
}

func createScootersIDs(t *testing.T, c *Client) []string {
	l := make([]string, 0)
	uid, err := c.CreateScooter()
	assert.NoError(t, err)
	assert.True(t, isValidUUID(uid.ID))
	l = append(l, uid.ID)

	uid, err = c.CreateScooter()
	assert.NoError(t, err)
	assert.True(t, isValidUUID(uid.ID))
	l = append(l, uid.ID)

	uid, err = c.CreateScooter()
	assert.NoError(t, err)
	assert.True(t, isValidUUID(uid.ID))
	l = append(l, uid.ID)

	return l
}

// isValidUUID validates uuid id
func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

// scooterInfoSliceToMap convert scooters' info from slice to a map
// the scooter's id as a key and true as value
func scooterInfoSliceToMap(ps []models.ScooterInfo) map[string]bool {
	m := make(map[string]bool)
	for _, ps := range ps {
		m[ps.ID] = true
	}
	return m
}

func setupEnv(t *testing.T) {
	err := os.Setenv("POSTGRES_USER", "postgres")
	assert.NoError(t, err)
	err = os.Setenv("POSTGRES_PASSWORD", "12345")
	assert.NoError(t, err)
	err = os.Setenv("POSTGRES_DATABASE", "dev_db")
	assert.NoError(t, err)
	err = os.Setenv("POSTGRES_HOST", "localhost")
	assert.NoError(t, err)
	err = os.Setenv("POSTGRES_PORT", "5432")
	assert.NoError(t, err)
}

// endTrips ends all the scooter trip after the time specified for the simulation
func endTrips(t *testing.T, ctx context.Context, ticker *time.Ticker, scooters []*Scooter) {
	for {
		select {
		case <-ticker.C:
			for _, x := range scooters {
				err := x.End(ctx)
				assert.NoError(t, err)
			}
			return
		}
	}
}
