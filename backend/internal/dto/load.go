package dto

type CreateLoadRequest struct {
    ExternalTMSLoadID string                 `json:"externalTMSLoadID"`
    FreightLoadID     string                 `json:"freightLoadID"`
    Status           StatusDTO              `json:"status"`
    Customer         map[string]interface{} `json:"customer"`
    BillTo          map[string]interface{} `json:"billTo"`
    Pickup          map[string]interface{} `json:"pickup"`
    Consignee       map[string]interface{} `json:"consignee"`
    Carrier         map[string]interface{} `json:"carrier"`
    RateData        map[string]interface{} `json:"rateData"`
    Specifications  map[string]interface{} `json:"specifications"`
    InPalletCount   int                   `json:"inPalletCount"`
    OutPalletCount  int                   `json:"outPalletCount"`
    NumCommodities  int                   `json:"numCommodities"`
    TotalWeight     float64               `json:"totalWeight"`
    BillableWeight  float64               `json:"billableWeight"`
    PoNums          string                `json:"poNums"`
    Operator        string                `json:"operator"`
    RouteMiles      float64               `json:"routeMiles"`
}

type StatusDTO struct {
    Code        StatusCodeDTO `json:"code"`
    Notes       string        `json:"notes"`
    Description string        `json:"description"`
}

type StatusCodeDTO struct {
    Key   string `json:"key"`
    Value string `json:"value"`
}

type LoadResponse struct {
    ID               string                 `json:"id"`
    ExternalTMSLoadID string                 `json:"externalTMSLoadID"`
    FreightLoadID     string                 `json:"freightLoadID"`
    Status           StatusDTO              `json:"status"`
    Customer         map[string]interface{} `json:"customer"`
    BillTo          map[string]interface{} `json:"billTo"`
    Pickup          map[string]interface{} `json:"pickup"`
    Consignee       map[string]interface{} `json:"consignee"`
    Carrier         map[string]interface{} `json:"carrier"`
    RateData        map[string]interface{} `json:"rateData"`
    Specifications  map[string]interface{} `json:"specifications"`
    InPalletCount   int                   `json:"inPalletCount"`
    OutPalletCount  int                   `json:"outPalletCount"`
    NumCommodities  int                   `json:"numCommodities"`
    TotalWeight     float64               `json:"totalWeight"`
    BillableWeight  float64               `json:"billableWeight"`
    PoNums          string                `json:"poNums"`
    Operator        string                `json:"operator"`
    RouteMiles      float64               `json:"routeMiles"`
    CreatedAt       string                `json:"createdAt"`
    UpdatedAt       string                `json:"updatedAt"`
}

type ListLoadsResponse struct {
    Loads []LoadResponse `json:"loads"`
    Total int64         `json:"total"`
    Page  int           `json:"page"`
    Size  int           `json:"size"`
}

