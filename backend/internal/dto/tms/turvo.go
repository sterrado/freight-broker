package dto

import "time"

type TurvoAuthRequest struct {
    GrantType    string `json:"grant_type"`
    ClientID     string `json:"client_id"`
    ClientSecret string `json:"client_secret"`
    Username     string `json:"username"`
    Password     string `json:"password"`
    Scope        string `json:"scope"`
    Type         string `json:"type"`
}

type TurvoAuthResponse struct {
    AccessToken  string `json:"access_token"`
    TokenType    string `json:"token_type"`
    ExpiresIn    int    `json:"expires_in"`
    Scope        string `json:"scope"`
    RefreshToken string `json:"refresh_token"`
    TenantRef    string `json:"tenant_ref"`
}

type DateInfo struct {
    Date     time.Time `json:"date"`
    TimeZone string    `json:"timeZone"`
}

type StatusCode struct {
    Key   string `json:"key"`
    Value string `json:"value"`
}

type Status struct {
    Code StatusCode `json:"code"`
}

type Lane struct {
    Start string `json:"start"`
    End   string `json:"end"`
}

type CustomerInfo struct {
    Name string `json:"name"`
}

type CustomerOrder struct {
    CustomerOrderSourceId string      `json:"customerOrderSourceId"`
    Customer             CustomerInfo `json:"customer"`
}

type CreateShipmentRequest struct {
    LTLShipment   bool            `json:"ltlShipment"`
    StartDate     DateInfo        `json:"startDate"`
    EndDate       DateInfo        `json:"endDate"`
    Status        Status          `json:"status"`
    Lane          Lane           `json:"lane"`
    CustomerOrder []CustomerOrder `json:"customerOrder"`
}

type ShipmentResponse struct {
    ID           int       `json:"id"`
    CustomID     string    `json:"customId"`
    Status       Status    `json:"status"`
    CustomerOrder []struct {
        ID       int `json:"id"`
        Customer struct {
            ID   int    `json:"id"`
            Name string `json:"name"`
            ParentAccount struct {
                Name string `json:"name"`
                Type string `json:"type"`
                ID   int    `json:"id"`
            } `json:"parentAccount"`
        } `json:"customer"`
    } `json:"customerOrder"`
    CarrierOrder []struct {
        ID      int `json:"id"`
        Carrier struct {
            ID   int    `json:"id"`
            Name string `json:"name"`
            ParentAccount struct {
                Name string `json:"name"`
                Type string `json:"type"`
                ID   int    `json:"id"`
            } `json:"parentAccount"`
        } `json:"carrier"`
    } `json:"carrierOrder"`
    Created      time.Time `json:"created"`
    Updated      time.Time `json:"updated"`
    LastUpdatedOn time.Time `json:"lastUpdatedOn"`
    CreatedDate  time.Time `json:"createdDate"`
}

type ListShipmentsResponse struct {
    Status  string `json:"Status"`
    Details struct {
        Pagination struct {
            Start              int  `json:"start"`
            PageSize           int  `json:"pageSize"`
            TotalRecordsInPage int  `json:"totalRecordsInPage"`
            MoreAvailable      bool `json:"moreAvailable"`
        } `json:"pagination"`
        Shipments []ShipmentResponse `json:"shipments"`
    } `json:"details"`
}