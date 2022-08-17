// Package models contains all the data representation for the service layer "REST-API".
package models

// NotOccupied is a constant for non-users to indicate the scooter is available
var NotOccupied = "NOT_OCCUPIED"

// ScooterInfo has the scooter details
type ScooterInfo struct {
	ID           string
	Coordination int64 // represents the scooter location.
	UserID       string
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
