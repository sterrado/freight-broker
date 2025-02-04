package services

import (
	"context"
	"fmt"
	"freight-broker/backend/internal/dto"
	"freight-broker/backend/internal/interfaces"
	"freight-broker/backend/internal/models"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type LoadService struct {
    db         *gorm.DB
    tmsService interfaces.TMSService
}

func NewLoadService(db *gorm.DB, tmsService interfaces.TMSService) *LoadService {
    return &LoadService{
        db:         db,
        tmsService: tmsService,
    }
}

func (s *LoadService) CreateLoad(ctx context.Context, req *dto.CreateLoadRequest) (*dto.LoadResponse, error) {
    load := &models.Load{
        ID:               uuid.New(),
        ExternalTMSLoadID: req.ExternalTMSLoadID,
        FreightLoadID:     req.FreightLoadID,
        Status:           models.JSON(map[string]interface{}{
            "code": map[string]interface{}{
                "key":   req.Status.Code.Key,
                "value": req.Status.Code.Value,
            },
            "notes":       req.Status.Notes,
            "description": req.Status.Description,
        }),
        Customer:         models.JSON(req.Customer),
        BillTo:          models.JSON(req.BillTo),
        Pickup:          models.JSON(req.Pickup),
        Consignee:       models.JSON(req.Consignee),
        Carrier:         models.JSON(req.Carrier),
        RateData:        models.JSON(req.RateData),
        Specifications:  models.JSON(req.Specifications),
        InPalletCount:   req.InPalletCount,
        OutPalletCount:  req.OutPalletCount,
        NumCommodities:  req.NumCommodities,
        TotalWeight:     req.TotalWeight,
        BillableWeight:  req.BillableWeight,
        PoNums:          req.PoNums,
        Operator:        req.Operator,
        RouteMiles:      req.RouteMiles,
    }

    if err := s.db.Create(load).Error; err != nil {
        return nil, fmt.Errorf("failed to create load: %w", err)
    }

    return s.convertToLoadResponse(load)
}

func (s *LoadService) GetLoad(ctx context.Context, id string) (*dto.LoadResponse, error) {
    var load models.Load
    
    if err := s.db.Where("id = ?", id).First(&load).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, fmt.Errorf("load not found")
        }
        return nil, fmt.Errorf("failed to get load: %w", err)
    }

    return s.convertToLoadResponse(&load)
}

func (s *LoadService) ListLoads(ctx context.Context, page, pageSize int) (*dto.ListLoadsResponse, error) {
    var loads []models.Load
    var total int64

    offset := (page - 1) * pageSize

    if err := s.db.Model(&models.Load{}).Count(&total).Error; err != nil {
        return nil, fmt.Errorf("failed to count loads: %w", err)
    }

    if err := s.db.Offset(offset).Limit(pageSize).Find(&loads).Error; err != nil {
        return nil, fmt.Errorf("failed to list loads: %w", err)
    }

    loadResponses := make([]dto.LoadResponse, len(loads))
    for i, load := range loads {
        response, err := s.convertToLoadResponse(&load)
        if err != nil {
            return nil, err
        }
        loadResponses[i] = *response
    }

    return &dto.ListLoadsResponse{
        Loads: loadResponses,
        Total: total,
    }, nil
}

// Helper function to convert model to DTO
func (s *LoadService) convertToLoadResponse(load *models.Load) (*dto.LoadResponse, error) {
    statusCode := load.Status["code"].(map[string]interface{})
    
    return &dto.LoadResponse{
        ID:               load.ID.String(),
        ExternalTMSLoadID: load.ExternalTMSLoadID,
        FreightLoadID:     load.FreightLoadID,
        Status: dto.StatusDTO{
            Code: dto.StatusCodeDTO{
                Key:   statusCode["key"].(string), 
                Value: statusCode["value"].(string),
            },
            Notes:       load.Status["notes"].(string),
            Description: load.Status["description"].(string),
        },
        Customer:         load.Customer,
        BillTo:          load.BillTo,
        Pickup:          load.Pickup,
        Consignee:       load.Consignee,
        Carrier:         load.Carrier,
        RateData:        load.RateData,
        Specifications:  load.Specifications,
        InPalletCount:   load.InPalletCount,
        OutPalletCount:  load.OutPalletCount,
        NumCommodities:  load.NumCommodities,
        TotalWeight:     load.TotalWeight,
        BillableWeight:  load.BillableWeight,
        PoNums:          load.PoNums,
        Operator:        load.Operator,
        RouteMiles:      load.RouteMiles,
        CreatedAt:       load.CreatedAt.Format(time.RFC3339),
        UpdatedAt:       load.UpdatedAt.Format(time.RFC3339),
    }, nil
}