package handlers

import (
	"context"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/TwiLightDM/diploma-gateway/internal/dto"
	terrain_service "github.com/TwiLightDM/diploma-gateway/internal/grpc/terrain-service"
	"github.com/labstack/echo/v4"
)

type TerrainHandler struct {
	client *terrain_service.TerrainClient
}

func NewTerrainHandler(client *terrain_service.TerrainClient) *TerrainHandler {
	return &TerrainHandler{client: client}
}

func (h *TerrainHandler) SubmitJob(c echo.Context) error {
	userID := c.Get("user_id").(string)

	file, err := c.FormFile("image")
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "image file required"})
	}
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "cannot open image"})
	}
	defer src.Close()

	imageBytes, err := io.ReadAll(src)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "cannot read image"})
	}

	scaleZ := float32(1.0)
	if v := c.FormValue("scale_z"); v != "" {
		if f, err := strconv.ParseFloat(v, 32); err == nil {
			scaleZ = float32(f)
		}
	}
	yUp := c.FormValue("y_up") == "true"
	textureMode := c.FormValue("texture_mode")
	if textureMode == "" {
		textureMode = "photo"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	jobID, err := h.client.SubmitJob(ctx, userID, imageBytes, scaleZ, yUp, textureMode)
	if err != nil {
		return grpcErr(c, err)
	}

	return c.JSON(http.StatusAccepted, dto.SubmitJobResponse{JobID: jobID})
}

func (h *TerrainHandler) GetJobStatus(c echo.Context) error {
	jobID := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := h.client.GetJobStatus(ctx, jobID)
	if err != nil {
		return grpcErr(c, err)
	}

	return c.JSON(http.StatusOK, dto.JobStatusResponse{
		JobID:      resp.JobId,
		Status:     resp.Status,
		ModelURL:   resp.ModelUrl,
		TextureURL: resp.TextureUrl,
		TreesURL:   resp.TreesUrl,
		Error:      resp.Error,
		CreatedAt:  resp.CreatedAt.AsTime(),
	})
}

func (h *TerrainHandler) GetHistory(c echo.Context) error {
	userID := c.Get("user_id").(string)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := h.client.GetUserJobs(ctx, userID)
	if err != nil {
		return grpcErr(c, err)
	}

	jobs := make([]dto.JobStatusResponse, 0, len(resp.Jobs))
	for _, j := range resp.Jobs {
		jobs = append(jobs, dto.JobStatusResponse{
			JobID:      j.JobId,
			Status:     j.Status,
			ModelURL:   j.ModelUrl,
			TextureURL: j.TextureUrl,
			TreesURL:   j.TreesUrl,
			Error:      j.Error,
			CreatedAt:  j.CreatedAt.AsTime(),
		})
	}
	return c.JSON(http.StatusOK, jobs)
}

func (h *TerrainHandler) GetMyLimit(c echo.Context) error {
	userID := c.Get("user_id").(string)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := h.client.GetLimit(ctx, userID)
	if err != nil {
		return grpcErr(c, err)
	}

	return c.JSON(http.StatusOK, dto.LimitResponse{
		UserID:     resp.UserId,
		DailyLimit: resp.DailyLimit,
		UsedToday:  resp.UsedToday,
		ResetAt:    resp.ResetAt.AsTime(),
	})
}

func (h *TerrainHandler) SetLimit(c echo.Context) error {
	var req dto.SetLimitRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := h.client.SetLimit(ctx, req.UserID, req.DailyLimit)
	if err != nil {
		return grpcErr(c, err)
	}

	return c.JSON(http.StatusOK, dto.LimitResponse{
		UserID:     resp.UserId,
		DailyLimit: resp.DailyLimit,
	})
}
