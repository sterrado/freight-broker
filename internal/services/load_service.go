package services

import (
    "context"
    "freight-broker/internal/interfaces"
    "freight-broker/internal/models"
    "freight-broker/internal/dto"
)

type LoadService struct {
    tmsProvider interfaces.TMSProvider
}

func NewLoadService(tmsProvider interfaces.TMSProvider) *LoadService {
    return &LoadService{
        tmsProvider: tmsProvider,
    }
}

func (s *LoadService) CreateLoad(ctx context.Context, createDTO *dto.CreateLoadDTO) (*dto.LoadResponseDTO, error) {
    // Implementation
}