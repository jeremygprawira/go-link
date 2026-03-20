package models

const (
	TYPE_HEALTH = "health"
)

// HealthStatus represents the health check response body.
type (
	HealthResponse struct {
		Description  string                 `json:"description"`
		Dependencies []HealthDetailResponse `json:"dependencies"`
	}

	HealthDetailResponse struct {
		Type        string `json:"type"`
		Component   string `json:"component"`
		Status      string `json:"status"`
		Description string `json:"description,omitempty"`
	}
)

