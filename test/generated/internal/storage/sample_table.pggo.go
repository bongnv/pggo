package model

import (
	"github.com/google/uuid"
)

// SampleTable represents sample_table table.
type SampleTable struct {
	ID   uuid.UUID
	Name string
}
