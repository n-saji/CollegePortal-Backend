package daos

import (
	"gorm.io/gorm"
)

type Daos struct {
	dbConn *gorm.DB
}

func New(conn *gorm.DB) *Daos {
	return &Daos{dbConn: conn}
}
