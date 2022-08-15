package db

var (
	scooterTable = `CREATE TABLE IF NOT EXISTS scooter
(
    id           TEXT          NOT NULL PRIMARY KEY,
	coordinate   INT           NOT NULL,
    user_id       TEXT,
);`

	UserTable = `CREATE TABLE IF NOT EXISTS users
(
    id           TEXT   NOT NULL PRIMARY KEY,
	Name         TEXT   NOT NULL,
    email        TEXT   NOT NULL PRIMARY KEY,

);`
)
