package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/TwiLightDM/diploma-gateway/internal/dto"
	"github.com/TwiLightDM/diploma-gateway/internal/entities"
	terrain_service "github.com/TwiLightDM/diploma-gateway/internal/grpc/terrain-service"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type LocationHandler struct {
	db      *gorm.DB
	terrain *terrain_service.TerrainClient
}

func NewLocationHandler(db *gorm.DB, terrain *terrain_service.TerrainClient) *LocationHandler {
	return &LocationHandler{db: db, terrain: terrain}
}

func (h *LocationHandler) Create(c echo.Context) error {
	userID := c.Get("user_id").(string)

	var req dto.CreateLocationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}
	if req.Name == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "name is required"})
	}

	extraJSON, removedJSON := marshalBrushFields(req.ExtraTrees, req.RemovedPositions)

	loc := &entities.Location{
		UserID:           userID,
		JobID:            req.JobID,
		Name:             req.Name,
		ScaleZ:           req.ScaleZ,
		WaterEnabled:     req.WaterEnabled,
		WaterLevel:       req.WaterLevel,
		TreeCountPercent: req.TreeCountPercent,
		CrownHeightScale: floatOr(req.CrownHeightScale, 1),
		CrownRadiusScale: floatOr(req.CrownRadiusScale, 1),
		TrunkRadiusScale: floatOr(req.TrunkRadiusScale, 1),
		ExtraTrees:       extraJSON,
		RemovedPositions: removedJSON,
	}
	if err := h.db.Create(loc).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	resp := locationToDTO(loc)
	h.fillJobURLs(c.Request().Context(), &resp)
	return c.JSON(http.StatusCreated, resp)
}

func (h *LocationHandler) List(c echo.Context) error {
	userID := c.Get("user_id").(string)

	var locs []entities.Location
	if err := h.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&locs).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	result := make([]dto.LocationResponse, 0, len(locs))
	for _, loc := range locs {
		r := locationToDTO(&loc)
		h.fillJobURLs(c.Request().Context(), &r)
		result = append(result, r)
	}
	return c.JSON(http.StatusOK, result)
}

func (h *LocationHandler) Update(c echo.Context) error {
	userID := c.Get("user_id").(string)
	locID := c.Param("id")

	var req dto.UpdateLocationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	var loc entities.Location
	if err := h.db.Where("id = ? AND user_id = ?", locID, userID).First(&loc).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "location not found"})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	extraJSON, removedJSON := marshalBrushFields(req.ExtraTrees, req.RemovedPositions)

	loc.Name = req.Name
	loc.ScaleZ = req.ScaleZ
	loc.WaterEnabled = req.WaterEnabled
	loc.WaterLevel = req.WaterLevel
	loc.TreeCountPercent = req.TreeCountPercent
	loc.CrownHeightScale = floatOr(req.CrownHeightScale, 1)
	loc.CrownRadiusScale = floatOr(req.CrownRadiusScale, 1)
	loc.TrunkRadiusScale = floatOr(req.TrunkRadiusScale, 1)
	loc.ExtraTrees = extraJSON
	loc.RemovedPositions = removedJSON

	if err := h.db.Save(&loc).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	resp := locationToDTO(&loc)
	h.fillJobURLs(c.Request().Context(), &resp)
	return c.JSON(http.StatusOK, resp)
}

func (h *LocationHandler) Delete(c echo.Context) error {
	userID := c.Get("user_id").(string)
	locID := c.Param("id")

	if err := h.db.Where("id = ? AND user_id = ?", locID, userID).Delete(&entities.Location{}).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *LocationHandler) fillJobURLs(ctx context.Context, r *dto.LocationResponse) {
	tCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if job, err := h.terrain.GetJobStatus(tCtx, r.JobID); err == nil {
		r.ModelURL = job.ModelUrl
		r.TextureURL = job.TextureUrl
		r.TreesURL = job.TreesUrl
	}
}

// marshalBrushFields serialises brush slices to JSON strings for storage.
// Returns "[]" (not null) when the slice is empty so the column always has valid JSON.
func marshalBrushFields(extra []dto.BrushTree, removed []dto.RemovedZone) (string, string) {
	if extra == nil {
		extra = []dto.BrushTree{}
	}
	if removed == nil {
		removed = []dto.RemovedZone{}
	}
	eb, _ := json.Marshal(extra)
	rb, _ := json.Marshal(removed)
	return string(eb), string(rb)
}

func floatOr(v, def float64) float64 {
	if v == 0 {
		return def
	}
	return v
}

func locationToDTO(loc *entities.Location) dto.LocationResponse {
	var extra []dto.BrushTree
	var removed []dto.RemovedZone
	_ = json.Unmarshal([]byte(loc.ExtraTrees), &extra)
	_ = json.Unmarshal([]byte(loc.RemovedPositions), &removed)
	if extra == nil {
		extra = []dto.BrushTree{}
	}
	if removed == nil {
		removed = []dto.RemovedZone{}
	}

	return dto.LocationResponse{
		ID:               loc.ID,
		UserID:           loc.UserID,
		JobID:            loc.JobID,
		Name:             loc.Name,
		ScaleZ:           loc.ScaleZ,
		WaterEnabled:     loc.WaterEnabled,
		WaterLevel:       loc.WaterLevel,
		TreeCountPercent: loc.TreeCountPercent,
		CrownHeightScale: floatOr(loc.CrownHeightScale, 1),
		CrownRadiusScale: floatOr(loc.CrownRadiusScale, 1),
		TrunkRadiusScale: floatOr(loc.TrunkRadiusScale, 1),
		ExtraTrees:       extra,
		RemovedPositions: removed,
		CreatedAt:        loc.CreatedAt,
		UpdatedAt:        loc.UpdatedAt,
	}
}
