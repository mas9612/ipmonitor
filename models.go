package ipmonitor

import (
	"time"
)

// Host represents host table in DB
type Host struct {
	ID          uint       `gorm:"primary_key"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"-"`
	DeletedAt   *time.Time `json:"-" sql:"index"`
	Address     string     `json:"address"`
	Hostname    string     `json:"hostname"`
	Description string     `json:"description"`
}
