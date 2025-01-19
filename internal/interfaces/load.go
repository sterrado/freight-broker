package interfaces

import (
    "context"
    "freight-broker/internal/dto"
)

type LoadService interface {
    CreateLoad(ctx context.Context, req *dto.CreateLoadRequest) (*dto.LoadResponse, error)
    GetLoad(ctx context.Context, id string) (*dto.LoadResponse, error)
    ListLoads(ctx context.Context, page, pageSize int) (*dto.ListLoadsResponse, error)
}