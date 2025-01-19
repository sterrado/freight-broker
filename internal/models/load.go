package models

import (
    "github.com/jinzhu/gorm"
    "github.com/google/uuid"
)

type Load struct {
    gorm.Model
    ID               uuid.UUID `gorm:"type:uuid;primary_key;"`
    ExternalTMSLoadID string    `gorm:"type:varchar(100)"`
    FreightLoadID    string    `gorm:"type:varchar(100)"`
    Status           string    `gorm:"type:varchar(50)"`
    Customer         JSON      `gorm:"type:jsonb"`
    BillTo          JSON      `gorm:"type:jsonb"`
    Pickup          JSON      `gorm:"type:jsonb"`
    Consignee       JSON      `gorm:"type:jsonb"`
    Carrier         JSON      `gorm:"type:jsonb"`
    RateData        JSON      `gorm:"type:jsonb"`
    Specifications  JSON      `gorm:"type:jsonb"`
    InPalletCount   int
    OutPalletCount  int
    NumCommodities  int
    TotalWeight     float64
    BillableWeight  float64
    PoNums          string    `gorm:"type:varchar(255)"`
    Operator        string    `gorm:"type:varchar(100)"`
    RouteMiles      float64
}

// JSON is a wrapper for handling JSON fields
type JSON map[string]interface{}