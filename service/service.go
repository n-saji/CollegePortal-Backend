package service

import (
	"CollegeAdministration/daos"

	"gorm.io/gorm"
)

type Service struct {
	daos *daos.Daos
}

func New(dbConn *gorm.DB) *Service {
	return &Service{
		daos: daos.New(dbConn),
	}
}
