package db

import (
	"context"
	"scootin/models"
)

// Repository represents storage operations
type Repository interface {
	// BookScooter assign the scooter for a user.
	BookScooter(ctx context.Context, ScooterID, userID string) error

	// ListAvailableScooter lists all available scooters
	ListAvailableScooter(ctx context.Context) ([]models.ScooterInfo, error)

	// CreateUser stores a new user
	CreateUser(ctx context.Context, user *models.User) error

	// CreateScooter creates a new scooter
	CreateScooter(ctx context.Context, scooterID string) error

	// UpdateScooterCoordinates update the scooter coordinates
	UpdateScooterCoordinates(ctx context.Context, scooterID string, coordinates int64) error

	// Close closes the database connection
	Close()
}

var repositoryImpl Repository

// BookScooter ...
func BookScooter(ctx context.Context, ScooterID, userID string) error {
	return repositoryImpl.BookScooter(ctx, ScooterID, userID)
}

// ListAvailableScooter ...
func ListAvailableScooter(ctx context.Context) ([]models.ScooterInfo, error) {
	return repositoryImpl.ListAvailableScooter(ctx)
}

// CreateScooter ...
func CreateScooter(ctx context.Context, scooterID string) error {
	return repositoryImpl.CreateScooter(ctx, scooterID)
}

// CreateUser ...
func CreateUser(ctx context.Context, user *models.User) error {
	return repositoryImpl.CreateUser(ctx, user)
}

// UpdateScooterCoordinates ...
func UpdateScooterCoordinates(ctx context.Context, scooterID string, coordinates int64) error {
	return repositoryImpl.UpdateScooterCoordinates(ctx, scooterID, coordinates)
}

// Close ...
func Close() {
	repositoryImpl.Close()
}

// SetRepository sets the repo implementation
func SetRepository(repository Repository) {
	repositoryImpl = repository
}
