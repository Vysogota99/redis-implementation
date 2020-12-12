package models

// SetHashRequest ...
type SetHashRequest struct {
	Key   string                 `json:"key" binding:"required"`
	Value map[string]interface{} `json:"value" binding:"required"`
	TTL   int                    `json:"ttl"`
}

// SetListRequest ...
type SetListRequest struct {
	Key   string        `json:"key" binding:"required"`
	Value []interface{} `json:"value" binding:"required"`
	TTL   int           `json:"ttl"`
}

// SetStringRequest ...
type SetStringRequest struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
	TTL   int    `json:"ttl"`
}
