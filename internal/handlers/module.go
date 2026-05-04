package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/TwiLightDM/diploma-gateway/internal/dto"
	"github.com/TwiLightDM/diploma-gateway/internal/grpc/course-service"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ModuleHandler struct {
	courseClient *course_service.CourseClient
}

func NewModuleHandler(courseClient *course_service.CourseClient) *ModuleHandler {
	return &ModuleHandler{courseClient: courseClient}
}

func (h *ModuleHandler) CreateModule(c echo.Context) error {
	var request dto.ModuleRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	response, err := h.courseClient.CreateModule(ctx, request.Title, request.Description, request.CourseId)
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.AlreadyExists:
				return c.JSON(http.StatusConflict, dto.ErrorResponse{Error: "course already exists"})
			default:
				return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
			}
		}
	}

	return c.JSON(http.StatusOK, dto.ModuleResponse{
		Id:          response.Module.Id,
		Title:       response.Module.Title,
		Description: response.Module.Description,
		Position:    response.Module.Position,
		CourseId:    response.Module.CourseId,
	})
}

func (h *ModuleHandler) ReadModule(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	response, err := h.courseClient.ReadModule(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.ModuleResponse{
		Id:              response.Module.Id,
		Title:           response.Module.Title,
		Description:     response.Module.Description,
		Position:        response.Module.Position,
		CourseId:        response.Module.CourseId,
		AmountOfLessons: int(response.Module.AmountOfLessons),
	})
}

func (h *ModuleHandler) ReadAllModulesByCourseId(c echo.Context) error {
	courseId := c.Param("course_id")
	if courseId == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	response, err := h.courseClient.ReadAllModulesByCourseId(ctx, courseId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	modules := make([]dto.ModuleResponse, 0, len(response.Modules))
	for _, module := range response.Modules {
		modules = append(modules, dto.ModuleResponse{
			Id:              module.Id,
			Title:           module.Title,
			Description:     module.Description,
			Position:        module.Position,
			CourseId:        module.CourseId,
			AmountOfLessons: int(module.AmountOfLessons),
		})
	}

	return c.JSON(http.StatusOK, dto.ModuleListResponse{
		Modules: modules,
	})
}

func (h *ModuleHandler) UpdateModule(c echo.Context) error {
	var request dto.ModuleRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	response, err := h.courseClient.UpdateModule(ctx, id, request.Title, request.Description, request.Position)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.ModuleResponse{
		Id:          response.Module.Id,
		Title:       response.Module.Title,
		Description: response.Module.Description,
		Position:    response.Module.Position,
		CourseId:    response.Module.CourseId,
	})
}

func (h *ModuleHandler) DeleteModule(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	_, err := h.courseClient.DeleteModule(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusNoContent, nil)
}
