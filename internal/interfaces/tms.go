package interfaces

import (
    "context"
    "freight-broker/internal/dto/tms"
)

type TMSService interface {
    // Auth methods
    Authenticate(ctx context.Context) error
    IsTokenValid() bool
    RefreshToken(ctx context.Context) error
    
    // Shipment methods
    CreateShipment(ctx context.Context, req dto.CreateShipmentRequest) (*dto.ShipmentResponse, error)
    GetShipment(ctx context.Context, id string) (*dto.ShipmentResponse, error)
    ListShipments(ctx context.Context, page, pageSize int) (*dto.ListShipmentsResponse, error)
    UpdateShipment(ctx context.Context, id string, req dto.CreateShipmentRequest) (*dto.ShipmentResponse, error)
    DeleteShipment(ctx context.Context, id string) error
}