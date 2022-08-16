package db

import (
	"fmt"
	"scootin/config"
)

func InitiatePostgre() error {
	pc, err := config.IniatilizePostgreConfig()
	if err != nil {
		return err
	}
	return setUpPostgre(pc)
}

func setUpPostgre(pc *config.PostgreConfig) error {
	postgersAddress := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", pc.PostgresHost, pc.PostgresPort, pc.PostgresUser, pc.PostgresPassword, pc.PostgresDataBase)
	repository, err := NewPostgre(postgersAddress)
	if err != nil {
		return err
	}
	SetRepository(repository)
	return nil
}
