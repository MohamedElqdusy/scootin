package client

import (
	"scootin/models"
	"testing"

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
	////////////////////  create the scooters  //////////////////////
	uid, err = c.CreateScooter()
	assert.NoError(t, err)
	assert.True(t, isValidUUID(uid.ID))
	sc1 := models.NewScooter(uid.ID)

	uid, err = c.CreateScooter()
	assert.NoError(t, err)
	assert.True(t, isValidUUID(uid.ID))
	sc2 := models.NewScooter(uid.ID)

	uid, err = c.CreateScooter()
	assert.NoError(t, err)
	assert.True(t, isValidUUID(uid.ID))
	sc3 := models.NewScooter(uid.ID)

	////////////////////  ListAvailableScooter  //////////////////////
	// map all the scooter info we have for easy "access" comparison
	// the key is the scooter id and the value its info
	actualScooterMap := make(map[string]models.ScooterInfo, 0)
	actualScooterMap[sc1.Info.ID] = sc1.Info
	actualScooterMap[sc2.Info.ID] = sc2.Info
	actualScooterMap[sc3.Info.ID] = sc3.Info

	// checks all available scooters
	scs, err := c.ListAvailableScooter()
	assert.NoError(t, err)

	// converts scooter to map for easy checking
	expectedScootersMap := sliceToMap(scs)
	// compare the created scooter against the available in the database
	for _, x := range expectedScootersMap {
		assert.EqualValues(t, actualScooterMap[x.ID], x)
	}

	////////////////////////////   BookScooter  ///////////////////////////
	// book the scooter sc1 by the user u1
	err = c.BookScooter(sc1.Info.ID, u1.ID)
	assert.NoError(t, err)

	// we should receive an error if we tried to book the same scooter by another user
	err = c.BookScooter(sc1.Info.ID, u2.ID)
	assert.Error(t, err)

	// book the scooter sc2 by the user u2
	err = c.BookScooter(sc2.Info.ID, u2.ID)
	assert.NoError(t, err)

	// checks all available scooters, we booked 2, so we have 1 left available
	scs, err = c.ListAvailableScooter()
	assert.NoError(t, err)
	assert.Len(t, scs, 1)
	notOccupied := ""
	assert.Equal(t, scs[0].UserID, notOccupied)

	////////////////////////////   ReleaseScooter  ///////////////////////////
	err = c.ReleaseScooter(u1.ID)
	assert.NoError(t, err)
	// checks all available scooters, we have 1 booked, so we have 2 left available
	scs, err = c.ListAvailableScooter()
	assert.NoError(t, err)
	assert.Len(t, scs, 2)
	assert.Equal(t, scs[0].UserID, notOccupied)
	assert.Equal(t, scs[1].UserID, notOccupied)

	err = c.ReleaseScooter(u2.ID)
	assert.NoError(t, err)
	// checks all available scooters, we have 0 booked, so we have 3 left available
	scs, err = c.ListAvailableScooter()
	assert.NoError(t, err)
	assert.Len(t, scs, 3)
	assert.Equal(t, scs[0].UserID, notOccupied)
	assert.Equal(t, scs[1].UserID, notOccupied)
	assert.Equal(t, scs[2].UserID, notOccupied)
}

// isValidUUID validates uuid id
func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

// sliceToMap convert scooters' info from slice to a map
// the scooter's id as a key and the scooter's info as value
func sliceToMap(ps []models.ScooterInfo) map[string]models.ScooterInfo {
	m := make(map[string]models.ScooterInfo)
	for _, ps := range ps {
		m[ps.ID] = ps
	}
	return m
}
