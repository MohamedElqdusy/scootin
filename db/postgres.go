package db

import (
	"context"
	"database/sql"
	"fmt"
	"scootin/models"

	_ "github.com/lib/pq"
)

type PostgreRepository struct {
	db *sql.DB
}

func NewPostgre(url string) (*PostgreRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(scooterTable); err != nil {
		return nil, fmt.Errorf("couldn't initate the Scooter table: %s", err)
	}
	if _, err := db.Exec(UserTable); err != nil {
		return nil, fmt.Errorf("couldn't initate the User table: %s", err)
	}
	return &PostgreRepository{
		db,
	}, nil
}

func (p *PostgreRepository) Close() {
	p.db.Close()
}

func (p *PostgreRepository) CreateScooter(ctx context.Context, scooterID string) error {
	_, err := p.db.Exec("INSERT INTO scooters(id) VALUES($1)", scooterID)
	return err
}

func (p *PostgreRepository) CreateUser(ctx context.Context, userID string) error {
	_, err := p.db.Exec("INSERT INTO uers(id) VALUES($1)", userID)
	return err
}

// BookScooter ...
func (p *PostgreRepository) BookScooter(ctx context.Context, ScooterID, userID string) error {
	_, err := p.db.Exec("UPDATE scooter SET user_id = $1 Where id = $2", userID, ScooterID)
	return err
}

// UpdateScooterCoordinates ...
func (p *PostgreRepository) UpdateScooterCoordinates(ctx context.Context, scooterID string, coordinates int64) error {
	_, err := p.db.Exec("UPDATE scooter SET coordinate = $2 Where id = $1", scooterID, coordinates)
	return err
}

// ListAvailableScooter ...
func (p *PostgreRepository) ListAvailableScooter(ctx context.Context) ([]models.ScooterInfo, error) {
	var (
		rows *sql.Rows
		err  error
	)
	notOccupied := ""
	if rows, err = p.db.Query("SELECT * FROM scooter Where user_id = $1", notOccupied); err != nil {
		return nil, err
	}
	defer rows.Close()
	return extractScooterInfo(rows)
}

func extractScooterInfo(rows *sql.Rows) ([]models.ScooterInfo, error) {
	infx := make([]models.ScooterInfo, 0)
	for rows.Next() {
		info := models.ScooterInfo{}
		if err := rows.Scan(&info.ID, &info.Coordination, &info.UserID); err != nil {
			return nil, err
		}
		infx = append(infx, info)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return infx, nil
}
