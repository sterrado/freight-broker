package models

import (
    "github.com/google/uuid"
    "time"
    "database/sql/driver"
    "encoding/json"
    "fmt"
)

type Load struct {
    ID               uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    CreatedAt        time.Time
    UpdatedAt        time.Time
    ExternalTMSLoadID string        `gorm:"type:varchar(100)"`
    FreightLoadID    string         `gorm:"type:varchar(100)"`
    Status           JSON           `gorm:"type:jsonb"`
    Customer         JSON           `gorm:"type:jsonb"`
    BillTo          JSON           `gorm:"type:jsonb"`
    Pickup          JSON           `gorm:"type:jsonb"`
    Consignee       JSON           `gorm:"type:jsonb"`
    Carrier         JSON           `gorm:"type:jsonb"`
    RateData        JSON           `gorm:"type:jsonb"`
    Specifications  JSON           `gorm:"type:jsonb"`
    InPalletCount   int
    OutPalletCount  int
    NumCommodities  int
    TotalWeight     float64
    BillableWeight  float64
    PoNums          string         `gorm:"type:varchar(255)"`
    Operator        string         `gorm:"type:varchar(100)"`
    RouteMiles      float64
}

// JSON is a wrapper for handling JSON fields
type JSON map[string]interface{}

func (j JSON) Value() (driver.Value, error) {
    if j == nil {
        return nil, nil
    }
    return json.Marshal(j)
}

func (j *JSON) Scan(value interface{}) error {
    if value == nil {
        *j = nil
        return nil
    }

    bytes, ok := value.([]byte)
    if !ok {
        return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
    }

    var result map[string]interface{}
    err := json.Unmarshal(bytes, &result)
    if err != nil {
        return err
    }

    *j = JSON(result)
    return nil
}