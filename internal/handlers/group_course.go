package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/TwiLightDM/diploma-gateway/internal/dto"
	"github.com/TwiLightDM/diploma-gateway/internal/grpc/user-service"
	"github.com/labstack/echo/v4"
)

type GroupMemberHandler struct {
	userClient *user_service.UserClient
}

func NewGroupMemberHandler(userClient *user_service.UserClient) *GroupMemberHandler {
	return &GroupMemberHandler{userClient: userClient}
}

func (h *GroupMemberHandler) CreateGroupMember(c echo.Context) error {
	var request dto.GroupMemberRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	response, err := h.userClient.CreateGroupMember(ctx, request.UserId, request.GroupId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.GroupMemberResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.GroupMemberResponse{
		Id:      response.GroupMember.Id,
		UserId:  response.GroupMember.UserId,
		GroupId: response.GroupMember.GroupId,
	})
}

func (h *GroupMemberHandler) ReadAllGroupMembersByUserId(c echo.Context) error {
	userId := c.Param("id")
	if userId == "" {
		return c.JSON(http.StatusBadRequest, dto.GroupMemberResponse{Error: "invalid request"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	response, err := h.userClient.ReadAllGroupMembersByUserId(ctx, userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.GroupMemberResponse{Error: err.Error()})
	}

	groupMembers := make([]dto.GroupMemberResponse, 0, len(response.GroupMembers))
	for _, groupMember := range response.GroupMembers {
		groupMembers = append(groupMembers, dto.GroupMemberResponse{
			Id:      groupMember.Id,
			UserId:  groupMember.UserId,
			GroupId: groupMember.GroupId,
		})
	}

	return c.JSON(http.StatusOK, dto.GroupMemberListResponse{
		GroupMembers: groupMembers,
	})
}

func (h *GroupMemberHandler) ReadAllGroupMembersByGroupId(c echo.Context) error {
	groupId := c.Param("id")
	if groupId == "" {
		return c.JSON(http.StatusBadRequest, dto.GroupMemberResponse{Error: "invalid request"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	response, err := h.userClient.ReadAllGroupMembersByGroupId(ctx, groupId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.GroupMemberResponse{Error: err.Error()})
	}

	groupMembers := make([]dto.GroupMemberResponse, 0, len(response.GroupMembers))
	for _, groupMember := range response.GroupMembers {
		groupMembers = append(groupMembers, dto.GroupMemberResponse{
			Id:      groupMember.Id,
			UserId:  groupMember.UserId,
			GroupId: groupMember.GroupId,
		})
	}

	return c.JSON(http.StatusOK, dto.GroupMemberListResponse{
		GroupMembers: groupMembers,
	})
}

func (h *GroupMemberHandler) DeleteGroupMember(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, dto.GroupMemberResponse{Error: "invalid request"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	_, err := h.userClient.DeleteGroupMember(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.GroupMemberResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusNoContent, nil)
}
