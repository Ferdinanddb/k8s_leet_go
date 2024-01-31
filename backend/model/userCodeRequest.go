package model

import (
	// "gorm.io/gorm"
	"database/sql"
	"time"

	"k8s_leet_code/database"
)


type UserCodeRequest struct {
    UserID  uint `gorm:"type:bigint;not null"`
	InstanciationTS time.Time `gorm:"type:timestamp with time zone;not null"`
	RequestUUID string `gorm:"type:text;not null"`
	CodeContent string `gorm:"type:text;not null"`
	WorkerStatus sql.NullString `gorm:"type:text"`
	OutputResult sql.NullString `gorm:"type:text"`
}

func (ucr *UserCodeRequest) Save() (*UserCodeRequest, error) {
	err := database.Database.Create(&ucr).Error
	if err != nil {
		return &UserCodeRequest{}, err
	}
	return ucr, nil
}
