package entities

import "time"

type Location struct {
	ID                string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID            string    `gorm:"not null;type:uuid;index"`
	JobID             string    `gorm:"not null;type:uuid"`
	Name              string    `gorm:"not null"`
	ScaleZ            float64   `gorm:"not null;default:0.3"`
	WaterEnabled      bool      `gorm:"not null;default:false"`
	WaterLevel        float64   `gorm:"not null;default:0.1"`
	TreeCountPercent  int       `gorm:"not null;default:100"`
	CrownHeightScale  float64   `gorm:"not null;default:1"`
	CrownRadiusScale  float64   `gorm:"not null;default:1"`
	TrunkRadiusScale  float64   `gorm:"not null;default:1"`
	ExtraTrees        string    `gorm:"type:text;not null;default:'[]'"`
	RemovedPositions  string    `gorm:"type:text;not null;default:'[]'"`
	CreatedAt         time.Time `gorm:"autoCreateTime"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime"`
}
