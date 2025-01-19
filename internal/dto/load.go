package dto

type CreateLoadDTO struct {
    Origin      string `json:"origin" binding:"required"`
    Destination string `json:"destination" binding:"required"`
    // Add more fields as needed
}

type LoadResponseDTO struct {
    ID          string `json:"id"`
    Origin      string `json:"origin"`
    Destination string `json:"destination"`
    Status      string `json:"status"`
    // Add more fields as needed
}