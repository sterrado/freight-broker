package controllers

import (
	"fmt"
	"freight-broker/internal/dto"
	tmsDTO "freight-broker/internal/dto/tms"
	"freight-broker/internal/interfaces"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LoadController struct {
    loadService interfaces.LoadService
    tmsService  interfaces.TMSService
}

func NewLoadController(loadService interfaces.LoadService, tmsService interfaces.TMSService) *LoadController {
    return &LoadController{
        loadService: loadService,
        tmsService:  tmsService,
    }
}

func (c *LoadController) CreateLoad(ctx *gin.Context) {
    var req dto.CreateLoadRequest

    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid request format",
            "details": err.Error(),
        })
        return
    }

    // Validate required nested structures
    if err := c.validateCreateLoadRequest(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": "Validation failed",
            "details": err.Error(),
        })
        return
    }

    // Ensure TMS is authenticated
    if !c.tmsService.IsTokenValid() {
        if err := c.tmsService.Authenticate(ctx); err != nil {
            ctx.JSON(http.StatusServiceUnavailable, gin.H{
                "error": "Failed to authenticate with TMS",
                "details": err.Error(),
            })
            return
        }
    }

    // Try to convert to shipment request
    shipmentReq := c.convertToShipmentRequest(&req)

    // Create shipment in TMS
    _, err := c.tmsService.CreateShipment(ctx, shipmentReq)
    // if err != nil {
    //     ctx.JSON(http.StatusServiceUnavailable, gin.H{
    //         "error": "Failed to create shipment in TMS",
    //         "details": err.Error(),
    //     })
    //     return
    // }
    log.Print(err)

    // Update the request with the TMS ID

    // Create load in local database
    loadResp, err := c.loadService.CreateLoad(ctx, &req)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to create load in local database",
            "details": err.Error(),
        })
        return
    }

    ctx.JSON(http.StatusCreated, loadResp)
}

func (c *LoadController) GetLoad(ctx *gin.Context) {
    // Get load ID from URL parameter
    id := ctx.Param("id")

    // Validate UUID format
    if _, err := uuid.Parse(id); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid load ID format",
            "details": "Load ID must be a valid UUID",
        })
        return
    }

    // Get load from service
    loadResp, err := c.loadService.GetLoad(ctx, id)
    if err != nil {
        if err.Error() == "load not found" {
            ctx.JSON(http.StatusNotFound, gin.H{
                "error": "Load not found",
                "details": err.Error(),
            })
            return
        }
        
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to get load",
            "details": err.Error(),
        })
        return
    }

    ctx.JSON(http.StatusOK, loadResp)
}

func (c *LoadController) ListLoads(ctx *gin.Context) {
    // Get pagination parameters from query string
    pageStr := ctx.DefaultQuery("page", "1")
    pageSizeStr := ctx.DefaultQuery("size", "10")

    // Convert and validate pagination parameters
    page, err := strconv.Atoi(pageStr)
    if err != nil || page < 1 {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid page parameter",
            "details": "Page must be a positive integer",
        })
        return
    }

    pageSize, err := strconv.Atoi(pageSizeStr)
    if err != nil || pageSize < 1 || pageSize > 100 {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid size parameter",
            "details": "Size must be a positive integer between 1 and 100",
        })
        return
    }

    // Get loads from service
    loadsResp, err := c.loadService.ListLoads(ctx, page, pageSize)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to list loads",
            "details": err.Error(),
        })
        return
    }

    // Add pagination info to response
    loadsResp.Page = page
    loadsResp.Size = pageSize

    ctx.JSON(http.StatusOK, loadsResp)
}


func (c *LoadController) validateCreateLoadRequest(req *dto.CreateLoadRequest) error {
    // Add validation logic based on your business rules
    if req.FreightLoadID == "" {
        return fmt.Errorf("freight load ID is required")
    }
    if req.Customer == nil {
        return fmt.Errorf("customer information is required")
    }
    if req.Pickup == nil {
        return fmt.Errorf("pickup information is required")
    }
    if req.Consignee == nil {
        return fmt.Errorf("consignee information is required")
    }
    // Add more validation as needed
    return nil
}

func (c *LoadController) convertToShipmentRequest(req *dto.CreateLoadRequest) tmsDTO.CreateShipmentRequest {
    // Parse scheduled times
    pickupTime, _ := time.Parse(time.RFC3339, req.Pickup["scheduledTime"].(string))
    deliveryTime, _ := time.Parse(time.RFC3339, req.Consignee["scheduledTime"].(string))

    // Get locations
    pickupLocation := req.Pickup["address"].(map[string]interface{})
    consigneeLocation := req.Consignee["address"].(map[string]interface{})

    // Convert status from load DTO to TMS DTO
    status := tmsDTO.Status{
        Code: tmsDTO.StatusCode{
            Key:   req.Status.Code.Key,
            Value: req.Status.Code.Value,
        },
    }

    return tmsDTO.CreateShipmentRequest{
        LTLShipment: true,
        StartDate: tmsDTO.DateInfo{
            Date:     pickupTime,
            TimeZone: "America/New_York",
        },
        EndDate: tmsDTO.DateInfo{
            Date:     deliveryTime,
            TimeZone: "America/New_York",
        },
        Status: status,
        Lane: tmsDTO.Lane{
            Start: fmt.Sprintf("%s, %s", 
                pickupLocation["city"].(string),
                pickupLocation["state"].(string)),
            End: fmt.Sprintf("%s, %s",
                consigneeLocation["city"].(string),
                consigneeLocation["state"].(string)),
        },
        CustomerOrder: []tmsDTO.CustomerOrder{{
            CustomerOrderSourceId: req.FreightLoadID,
            Customer: tmsDTO.CustomerInfo{
                Name: req.Customer["name"].(string),
            },
        }},
    }
}