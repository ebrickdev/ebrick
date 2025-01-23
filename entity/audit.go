package entity

import (
	"time"

	"gorm.io/gorm"
)

// Auditing provides common auditing fields for entities.
type Auditing struct {
	// CreatedAt records when the entity was created.
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	// CreatedBy records who created the entity.
	CreatedBy string `gorm:"type:varchar(255)" json:"created_by,omitempty"`

	// UpdatedAt records when the entity was last updated.
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// UpdatedBy records who last updated the entity.
	UpdatedBy string `gorm:"type:varchar(255)" json:"updated_by,omitempty"`

	// DeletedAt records when the entity was soft deleted.
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
