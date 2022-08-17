package db

import (
	"context"
	"database/sql"
	"errors"
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
	if _, err := db.Exec(userTable); err != nil {
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
	_, err := p.db.Exec("INSERT INTO scooters(id,coordinate,user_id) VALUES($1,$2,$3)", scooterID, 1, models.NotOccupied)
	return err
}

func (p *PostgreRepository) CreateUser(ctx context.Context, user *models.User) error {
	_, err := p.db.Exec("INSERT INTO users(id,name,email) VALUES($1,$2,$3)", user.ID, user.Name, user.Email)
	return err
}

// BookScooter ...
func (p *PostgreRepository) BookScooter(ctx context.Context, ScooterID, userID string) error {
	var (
		txn *sql.Tx
		err error
		//row        *sql.Row
	)
	// start the transaction
	if txn, err = p.db.Begin(); err != nil {
		return err
	}

	// checks whether the scooter is already occupied by a user
	res, err := txn.Exec("Select user_id FROM scooters WHERE id = $1 ", ScooterID)
	/*res.
	if err := row.Scan(&occupation); err != nil {
		return err
	}
	if occupation != models.NotOccupied {
		return errors.New(fmt.Sprintf("we can't book the scooter %s for user $s as it's already occupied by user %s", ScooterID, userID, occupation))
	}
	*/
	// do the booking only if the scooter is not booked by another user
	if res, err = txn.Exec("UPDATE scooters SET user_id = $1 Where id = $2 AND user_id = $3", userID, ScooterID, models.NotOccupied); err != nil {
		return err
	}
	if rowsCountAffected, err := res.RowsAffected(); err != nil {
		return err
	} else if rowsCountAffected == 0 {
		return errors.New(fmt.Sprintf("we can't book the scooter %s for user %s as it's already occupied", ScooterID, userID))
	}
	return txn.Commit()
}

// ReleaseScooter ...
func (p *PostgreRepository) ReleaseScooter(ctx context.Context, userID string) error {
	_, err := p.db.Exec("UPDATE scooters SET user_id = $1 Where user_id = $2", models.NotOccupied, userID)
	return err
}

// UpdateScooterCoordinates ...
func (p *PostgreRepository) UpdateScooterCoordinates(ctx context.Context, scooterID string, coordinates int64) error {
	_, err := p.db.Exec("UPDATE scooters SET coordinate = $2 Where id = $1", scooterID, coordinates)
	return err
}

// ListAvailableScooter ...
func (p *PostgreRepository) ListAvailableScooter(ctx context.Context) ([]models.ScooterInfo, error) {
	var (
		rows *sql.Rows
		err  error
	)
	if rows, err = p.db.Query("SELECT * FROM scooters Where user_id = $1", models.NotOccupied); err != nil {
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
