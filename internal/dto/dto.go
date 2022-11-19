package dto

// Metric is full dto both for requests and responses
type Metric struct {
	ID    string   `json:"id"`              // metric name
	MType string   `json:"type"`            // metric type
	Delta *int64   `json:"delta,omitempty"` // counter value
	Value *float64 `json:"value,omitempty"` // gauge value
}

// GetMetricRequest to get metric by name (ID) and type
type GetMetricRequest struct {
	ID    string `json:"id"`   // metric name
	MType string `json:"type"` // metric type
}

// ErrorResponse is common dto for error responses
type ErrorResponse struct {
	Error string `json:"error,omitempty"` // error message
}
