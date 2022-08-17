package db

var (
	scooterTable = `CREATE TABLE IF NOT EXISTS scooters
(
    id           TEXT          NOT NULL PRIMARY KEY,
	coordinate   INT,
    user_id       TEXT
);`

	userTable = `CREATE TABLE IF NOT EXISTS users
(
    id           TEXT   NOT NULL PRIMARY KEY,
	Name         TEXT   NOT NULL,
    email        TEXT   NOT NULL
);`
)
