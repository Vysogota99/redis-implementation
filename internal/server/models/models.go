package models

// SetHashRequest ...
type SetHashRequest struct {
	Key   interface{}            `json:"key" binding:"required"`
	Value map[string]interface{} `json:"value" binding:"required"`
	TTL   int                    `json:"ttl"`
}

// SetListRequest ...
type SetListRequest struct {
	Key   interface{}   `json:"key" binding:"required"`
	Value []interface{} `json:"value" binding:"required"`
	TTL   int           `json:"ttl"`
}

// SetStringRequest ...
type SetStringRequest struct {
	Key   interface{} `json:"key" binding:"required"`
	Value string      `json:"value" binding:"required"`
	TTL   int         `json:"ttl"`
}

// ListElement - элемент массива для идентификации типа данных
type ListElement struct {
	Dtype string
	Data  string
}
