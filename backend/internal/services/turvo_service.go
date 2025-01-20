package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"freight-broker/backend/internal/dto/tms"
)

const (
    sandboxAuthURL = "https://my-sandbox-publicapi.turvo.com/v1/oauth/token"
    prodAuthURL    = "https://publicapi.turvo.com/v1/oauth/token"
    baseShipmentsURL = "/shipments"
    baseCustomersURL = "/customers/list"

)

type TMSServiceConfig struct {
    APIKey       string
    ClientID     string
    ClientSecret string
    IsSandbox    bool
    TurvoUsername string
    TurvoPassword string
}

type TurvoService struct {
    config       TMSServiceConfig
    client       *http.Client
    authToken    string
    tokenExpiry  time.Time
    mu           sync.RWMutex
}

func NewTurvoService(config TMSServiceConfig) *TurvoService {
    return &TurvoService{
        config: config,
        client: &http.Client{
            Timeout: time.Second * 30,
        },
    }
}

func (s *TurvoService) Authenticate(ctx context.Context) error {

    authURL := prodAuthURL
    if s.config.IsSandbox {
        authURL = sandboxAuthURL
    }


    authReq := dto.TurvoAuthRequest{
        GrantType:    "password",
        ClientID:     s.config.ClientID,
        ClientSecret: s.config.ClientSecret,
        Username: s.config.TurvoUsername,
        Password: s.config.TurvoPassword,
        Scope:        "read+trust+write",
        Type:         "business",
    }


    jsonBody, err := json.Marshal(authReq)
    if err != nil {
        return fmt.Errorf("failed to marshal auth request: %w", err)
    }

    req, err := http.NewRequestWithContext(ctx, "POST", authURL, bytes.NewBuffer(jsonBody))
    if err != nil {
        return fmt.Errorf("failed to create request: %w", err)
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("x-api-key", s.config.APIKey)

    resp, err := s.client.Do(req)
    if err != nil {
        return fmt.Errorf("failed to make request: %w", err)
    }
    defer resp.Body.Close()

    bodyBytes, err := io.ReadAll(resp.Body)
    if err != nil {
        return fmt.Errorf("failed to read response body: %w", err)
    }

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("authentication failed with status: %d, body: %s", resp.StatusCode, string(bodyBytes))
    }

    var authResp dto.TurvoAuthResponse
    if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&authResp); err != nil {
        return fmt.Errorf("failed to decode response: %w", err)
    }

    s.mu.Lock()
    s.authToken = authResp.AccessToken
    s.tokenExpiry = time.Now().Add(time.Second * time.Duration(authResp.ExpiresIn))
    s.mu.Unlock()


    return nil
}

func (s *TurvoService) IsTokenValid() bool {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    if s.authToken == "" {
        return false
    }
    
    return time.Now().Add(5 * time.Minute).Before(s.tokenExpiry)
}

func (s *TurvoService) RefreshToken(ctx context.Context) error {
    return s.Authenticate(ctx)
}

func (s *TurvoService) GetAuthToken() string {
    s.mu.RLock()
    defer s.mu.RUnlock()
    return s.authToken
}

func (s *TurvoService) CreateShipment(ctx context.Context, req dto.CreateShipmentRequest) (*dto.ShipmentResponse, error) {
    url := s.getBaseURL() + baseShipmentsURL
    
    jsonData, err := json.Marshal(req)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request: %w", err)
    }
    
    log.Printf("Creating shipment with payload: %s", string(jsonData))

    httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }

    s.setAuthHeaders(httpReq)

    resp, err := s.client.Do(httpReq)
    if err != nil {
        return nil, fmt.Errorf("failed to make request: %w", err)
    }
    defer resp.Body.Close()

    bodyBytes, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response body: %w", err)
    }
    
    log.Printf("Response status: %d", resp.StatusCode)
    log.Printf("Response body: %s", string(bodyBytes))

    if resp.StatusCode != http.StatusOK {
        bodyBytes, _ := io.ReadAll(resp.Body)
        log.Printf("Failed request payload: %s", string(jsonData))
        log.Printf("Error response: %s", string(bodyBytes))
        return nil, fmt.Errorf("API returned status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
    }

    var shipmentResp dto.ShipmentResponse
    if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&shipmentResp); err != nil {
        return nil, fmt.Errorf("failed to decode response: %w", err)
    }

    return &shipmentResp, nil
}

func (s *TurvoService) ListShipments(ctx context.Context, page, pageSize int) (*dto.ListShipmentsResponse, error) {
    url := fmt.Sprintf("%s%s/list?start=%d&pageSize=%d", 
        s.getBaseURL(), baseShipmentsURL, page, pageSize)
    
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }

    s.setAuthHeaders(req)

    resp, err := s.client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to make request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
    }

    var listResp dto.ListShipmentsResponse
    if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
        return nil, fmt.Errorf("failed to decode response: %w", err)
    }

    return &listResp, nil
}

func (s *TurvoService) GetShipment(ctx context.Context, id string) (*dto.ShipmentResponse, error) {
    url := fmt.Sprintf("%s%s/%s", s.getBaseURL(), baseShipmentsURL, id)
    
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }

    s.setAuthHeaders(req)

    resp, err := s.client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to make request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
    }

    var shipment dto.ShipmentResponse
    if err := json.NewDecoder(resp.Body).Decode(&shipment); err != nil {
        return nil, fmt.Errorf("failed to decode response: %w", err)
    }

    return &shipment, nil
}
func (s *TurvoService) UpdateShipment(ctx context.Context, id string, req dto.CreateShipmentRequest) (*dto.ShipmentResponse, error) {
    url := fmt.Sprintf("%s%s/%s", s.getBaseURL(), baseShipmentsURL, id)
    
    jsonData, err := json.Marshal(req)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request: %w", err)
    }

    httpReq, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }

    s.setAuthHeaders(httpReq)

    resp, err := s.client.Do(httpReq)
    if err != nil {
        return nil, fmt.Errorf("failed to make request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        var errResp struct {
            Message string `json:"message"`
            Details string `json:"details,omitempty"`
        }
        if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
            return nil, fmt.Errorf("API returned status code: %d", resp.StatusCode)
        }
        return nil, fmt.Errorf("API error: %s - %s", errResp.Message, errResp.Details)
    }

    var shipmentResp dto.ShipmentResponse
    if err := json.NewDecoder(resp.Body).Decode(&shipmentResp); err != nil {
        return nil, fmt.Errorf("failed to decode response: %w", err)
    }

    return &shipmentResp, nil
}

func (s *TurvoService) DeleteShipment(ctx context.Context, id string) error {
    url := fmt.Sprintf("%s%s/%s", s.getBaseURL(), baseShipmentsURL, id)
    
    req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
    if err != nil {
        return fmt.Errorf("failed to create request: %w", err)
    }

    s.setAuthHeaders(req)

    resp, err := s.client.Do(req)
    if err != nil {
        return fmt.Errorf("failed to make request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
        var errResp struct {
            Message string `json:"message"`
            Details string `json:"details,omitempty"`
        }
        if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
            return fmt.Errorf("API returned status code: %d", resp.StatusCode)
        }
        return fmt.Errorf("API error: %s - %s", errResp.Message, errResp.Details)
    }

    return nil
}

// Helper methods
func (s *TurvoService) setAuthHeaders(req *http.Request) {
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("x-api-key", s.config.APIKey)
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.GetAuthToken()))
}

func (s *TurvoService) getBaseURL() string {
    if s.config.IsSandbox {
        return "https://my-sandbox-publicapi.turvo.com/v1"
    }
    return "https://publicapi.turvo.com/v1"
}
