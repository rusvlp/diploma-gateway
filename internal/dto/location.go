package dto

import "time"

// BrushTree — manually painted tree, stored in world-space coordinates
// (terrain is always centered at origin, so world coords are stable per job).
type BrushTree struct {
	WX float64 `json:"wx"`
	WY float64 `json:"wy"`
	WZ float64 `json:"wz"`
	S  float64 `json:"s"`
}

// RemovedZone — normalized-coordinate zone where auto-generated trees are hidden.
type RemovedZone struct {
	X float64 `json:"x"`
	Z float64 `json:"z"`
	R float64 `json:"r"`
}

type CreateLocationRequest struct {
	JobID             string        `json:"job_id"`
	Name              string        `json:"name"`
	ScaleZ            float64       `json:"scale_z"`
	WaterEnabled      bool          `json:"water_enabled"`
	WaterLevel        float64       `json:"water_level"`
	TreeCountPercent  int           `json:"tree_count_percent"`
	CrownHeightScale  float64       `json:"crown_height_scale"`
	CrownRadiusScale  float64       `json:"crown_radius_scale"`
	TrunkRadiusScale  float64       `json:"trunk_radius_scale"`
	ExtraTrees        []BrushTree   `json:"extra_trees"`
	RemovedPositions  []RemovedZone `json:"removed_positions"`
}

type UpdateLocationRequest struct {
	Name              string        `json:"name"`
	ScaleZ            float64       `json:"scale_z"`
	WaterEnabled      bool          `json:"water_enabled"`
	WaterLevel        float64       `json:"water_level"`
	TreeCountPercent  int           `json:"tree_count_percent"`
	CrownHeightScale  float64       `json:"crown_height_scale"`
	CrownRadiusScale  float64       `json:"crown_radius_scale"`
	TrunkRadiusScale  float64       `json:"trunk_radius_scale"`
	ExtraTrees        []BrushTree   `json:"extra_trees"`
	RemovedPositions  []RemovedZone `json:"removed_positions"`
}

type LocationResponse struct {
	ID                string        `json:"id"`
	UserID            string        `json:"user_id"`
	JobID             string        `json:"job_id"`
	Name              string        `json:"name"`
	ScaleZ            float64       `json:"scale_z"`
	WaterEnabled      bool          `json:"water_enabled"`
	WaterLevel        float64       `json:"water_level"`
	TreeCountPercent  int           `json:"tree_count_percent"`
	CrownHeightScale  float64       `json:"crown_height_scale"`
	CrownRadiusScale  float64       `json:"crown_radius_scale"`
	TrunkRadiusScale  float64       `json:"trunk_radius_scale"`
	ExtraTrees        []BrushTree   `json:"extra_trees"`
	RemovedPositions  []RemovedZone `json:"removed_positions"`
	ModelURL          string        `json:"model_url,omitempty"`
	TextureURL        string        `json:"texture_url,omitempty"`
	TreesURL          string        `json:"trees_url,omitempty"`
	CreatedAt         time.Time     `json:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at"`
}
