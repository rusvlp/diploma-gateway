package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/TwiLightDM/diploma-gateway/internal/dto"
	"github.com/TwiLightDM/diploma-gateway/internal/grpc/user-service"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userClient *user_service.UserClient
}

func NewUserHandler(userClient *user_service.UserClient) *UserHandler {
	return &UserHandler{userClient: userClient}
}

func (h *UserHandler) Login(c echo.Context) error {
	var request dto.UserRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	response, err := h.userClient.Login(ctx, request.Email, request.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.LoginResponse{
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
	})
}

func (h *UserHandler) SignUp(c echo.Context) error {
	var request dto.UserRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	response, err := h.userClient.SignUp(ctx, request.FullName, "student", request.Email, request.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SignUpResponse{
		User: dto.UserResponse{
			Id:       response.User.Id,
			FullName: response.User.FullName,
			Role:     response.User.Role,
			Email:    response.User.Email,
		},
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
	})
}

func (h *UserHandler) CreateTeacher(c echo.Context) error {
	var request dto.UserRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	response, err := h.userClient.SignUp(ctx, request.FullName, "teacher", request.Email, request.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.SignUpResponse{
		User: dto.UserResponse{
			Id:       response.User.Id,
			FullName: response.User.FullName,
			Role:     response.User.Role,
			Email:    response.User.Email,
		},
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
	})
}

func (h *UserHandler) Refresh(c echo.Context) error {
	id := c.Get("user_id").(string)
	role := c.Get("role").(string)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	response, err := h.userClient.Refresh(ctx, id, role)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.LoginResponse{
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
	})
}

func (h *UserHandler) ReadUser(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	response, err := h.userClient.ReadUser(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.UserResponse{
		Id:       response.User.Id,
		FullName: response.User.FullName,
		Role:     response.User.Role,
		Email:    response.User.Email,
	})
}

func (h *UserHandler) ReadAllUser(c echo.Context) error {
	role := c.Get("role").(string)
	if role != "admin" {
		return c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "you need to be an admin"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	response, err := h.userClient.ReadUsers(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	users := make([]dto.UserResponse, 0, len(response.Users))
	for _, user := range response.Users {
		users = append(users, dto.UserResponse{
			Id:       user.Id,
			FullName: user.FullName,
			Role:     user.Role,
			Email:    user.Email,
		})
	}

	return c.JSON(http.StatusOK, dto.UserListResponse{
		Users: users,
	})
}

func (h *UserHandler) ReadSelf(c echo.Context) error {
	id := c.Get("user_id").(string)
	if id == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	response, err := h.userClient.ReadUser(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.UserResponse{
		Id:       response.User.Id,
		FullName: response.User.FullName,
		Role:     response.User.Role,
		Email:    response.User.Email,
	})
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	var request dto.UserRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	id := c.Get("user_id").(string)
	if id == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	response, err := h.userClient.UpdateUser(ctx, id, request.FullName, request.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.UserResponse{
		Id:       response.User.Id,
		FullName: response.User.FullName,
		Role:     response.User.Role,
		Email:    response.User.Email,
	})
}

func (h *UserHandler) ChangePassword(c echo.Context) error {
	var request dto.UserRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	id := c.Get("user_id").(string)
	if id == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	_, err := h.userClient.ChangePassword(ctx, id, request.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.NoContent(http.StatusOK)
}

func (h *UserHandler) UpdateUserRole(c echo.Context) error {
	var request dto.UserRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	role := c.Get("role").(string)
	if role != "admin" {
		return c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "you need to be an admin"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	response, err := h.userClient.UpdateUserRole(ctx, request.UserId, request.Role)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.UserResponse{
		Id:       response.User.Id,
		FullName: response.User.FullName,
		Role:     response.User.Role,
		Email:    response.User.Email,
	})
}
