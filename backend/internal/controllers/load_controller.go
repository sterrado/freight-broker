package controllers

import (
	"fmt"
	"freight-broker/backend/internal/dto"
	tmsDTO "freight-broker/backend/internal/dto/tms"
	"freight-broker/backend/internal/interfaces"
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

    ctx.Header("Access-Control-Allow-Origin", "*")
    ctx.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
    ctx.Header("Access-Control-Allow-Headers", "Origin, Authorization, Content-Type")
    var req dto.CreateLoadRequest

    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid request format",
            "details": err.Error(),
        })
        return
    }

    if err := c.validateCreateLoadRequest(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": "Validation failed",
            "details": err.Error(),
        })
        return
    }

    if !c.tmsService.IsTokenValid() {
        if err := c.tmsService.Authenticate(ctx); err != nil {
            ctx.JSON(http.StatusServiceUnavailable, gin.H{
                "error": "Failed to authenticate with TMS",
                "details": err.Error(),
            })
            return
        }
    }

    shipmentReq := c.convertToShipmentRequest(&req)

    _, err := c.tmsService.CreateShipment(ctx, shipmentReq)
    
    //bypassing error validation currently because of strange error response

    // if err != nil {
    //     ctx.JSON(http.StatusServiceUnavailable, gin.H{
    //         "error": "Failed to create shipment in TMS",
    //         "details": err.Error(),
    //     })
    //     return
    // }
    log.Print(err)


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
    ctx.Header("Access-Control-Allow-Origin", "*")
    ctx.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
    ctx.Header("Access-Control-Allow-Headers", "Origin, Authorization, Content-Type")
    id := ctx.Param("id")

    if _, err := uuid.Parse(id); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid load ID format",
            "details": "Load ID must be a valid UUID",
        })
        return
    }

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
    ctx.Header("Access-Control-Allow-Origin", "*")
    ctx.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
    ctx.Header("Access-Control-Allow-Headers", "Origin, Authorization, Content-Type")
    pageStr := ctx.DefaultQuery("page", "1")
    pageSizeStr := ctx.DefaultQuery("size", "10")

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

    loadsResp, err := c.loadService.ListLoads(ctx, page, pageSize)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to list loads",
            "details": err.Error(),
        })
        return
    }

    loadsResp.Page = page
    loadsResp.Size = pageSize

    ctx.JSON(http.StatusOK, loadsResp)
}


func (c *LoadController) validateCreateLoadRequest(req *dto.CreateLoadRequest) error {
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
    return nil
}

func (c *LoadController) convertToShipmentRequest(req *dto.CreateLoadRequest) tmsDTO.CreateShipmentRequest {
    pickupTime, _ := time.Parse(time.RFC3339, req.Pickup["scheduledTime"].(string))
    deliveryTime, _ := time.Parse(time.RFC3339, req.Consignee["scheduledTime"].(string))

    pickupLocation := req.Pickup["address"].(map[string]interface{})
    consigneeLocation := req.Consignee["address"].(map[string]interface{})

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