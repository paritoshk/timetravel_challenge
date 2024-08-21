package entity

import "time"

// Record represents a versioned record in the system
type Record struct {
	ID        int               `json:"id"`
	Data      map[string]string `json:"data"`
	Version   int               `json:"version"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// Copy creates a deep copy of the Record
func (r *Record) Copy() Record {
	newData := make(map[string]string, len(r.Data))
	for key, value := range r.Data {
		newData[key] = value
	}

	return Record{
		ID:        r.ID,
		Data:      newData,
		Version:   r.Version,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}
