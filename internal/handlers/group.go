package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/TwiLightDM/diploma-gateway/internal/dto"
	"github.com/TwiLightDM/diploma-gateway/internal/grpc/user-service"
	"github.com/labstack/echo/v4"
)

type GroupHandler struct {
	userClient *user_service.UserClient
}

func NewGroupHandler(userClient *user_service.UserClient) *GroupHandler {
	return &GroupHandler{userClient: userClient}
}

func (h *GroupHandler) CreateGroup(c echo.Context) error {
	var request dto.GroupRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	response, err := h.userClient.CreateGroup(ctx, request.Title, request.Description, request.OwnerId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.GroupResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.GroupResponse{
		Id:          response.Group.Id,
		Title:       response.Group.Title,
		Description: response.Group.Description,
		OwnerId:     response.Group.OwnerId,
	})
}

func (h *GroupHandler) ReadGroup(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, dto.GroupResponse{Error: "invalid request"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	response, err := h.userClient.ReadGroup(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.GroupResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.GroupResponse{
		Id:          response.Group.Id,
		Title:       response.Group.Title,
		Description: response.Group.Description,
		OwnerId:     response.Group.OwnerId,
	})
}

func (h *GroupHandler) ReadAllGroupsByOwnerId(c echo.Context) error {
	ownerId := c.Param("id")
	if ownerId == "" {
		return c.JSON(http.StatusBadRequest, dto.GroupResponse{Error: "invalid request"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	response, err := h.userClient.ReadAllGroupsByOwnerId(ctx, ownerId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.GroupResponse{Error: err.Error()})
	}

	groups := make([]dto.GroupResponse, 0, len(response.Groups))
	for _, group := range response.Groups {
		groups = append(groups, dto.GroupResponse{
			Id:          group.Id,
			Title:       group.Title,
			Description: group.Description,
			OwnerId:     group.OwnerId,
		})
	}

	return c.JSON(http.StatusOK, dto.GroupListResponse{
		Groups: groups,
	})
}

func (h *GroupHandler) UpdateGroup(c echo.Context) error {
	var request dto.GroupRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.GroupResponse{Error: "invalid request"})
	}

	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, dto.GroupResponse{Error: "invalid request"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	response, err := h.userClient.UpdateGroup(ctx, id, request.Title, request.Description)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.GroupResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.GroupResponse{
		Id:          response.Group.Id,
		Title:       response.Group.Title,
		Description: response.Group.Description,
		OwnerId:     response.Group.OwnerId,
	})
}

func (h *GroupHandler) DeleteGroup(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, dto.GroupResponse{Error: "invalid request"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	_, err := h.userClient.DeleteGroup(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.GroupResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusNoContent, nil)
}
