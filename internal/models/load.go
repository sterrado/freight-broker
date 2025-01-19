package models

import (
    "github.com/jinzhu/gorm"
    "github.com/google/uuid"
)

type Load struct {
    gorm.Model
    ID          uuid.UUID `gorm:"type:uuid;primary_key;"`
    BrokerID    uuid.UUID `gorm:"type:uuid;not null"`
    Origin      string    `gorm:"not null"`
    Destination string    `gorm:"not null"`
    Status      string    `gorm:"not null"`
    // Add more fields as needed
}