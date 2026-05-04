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

type CourseHandler struct {
	courseClient *course_service.CourseClient
}

func NewCourseHandler(courseClient *course_service.CourseClient) *CourseHandler {
	return &CourseHandler{courseClient: courseClient}
}

func (h *CourseHandler) CreateCourse(c echo.Context) error {
	var request dto.CourseRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	ownerId := c.Get("user_id").(string)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	response, err := h.courseClient.CreateCourse(ctx, request.Title, request.Description, request.AccessType, ownerId)
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

	return c.JSON(http.StatusOK, dto.CourseResponse{
		Id:          response.Course.Id,
		Title:       response.Course.Title,
		Description: response.Course.Description,
		AccessType:  response.Course.AccessType,
		PublishedAt: response.Course.PublishedAt,
		OwnerId:     response.Course.OwnerId,
	})
}

func (h *CourseHandler) ReadCourse(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	response, err := h.courseClient.ReadCourse(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.CourseResponse{
		Id:              response.Course.Id,
		Title:           response.Course.Title,
		Description:     response.Course.Description,
		AccessType:      response.Course.AccessType,
		PublishedAt:     response.Course.PublishedAt,
		OwnerId:         response.Course.OwnerId,
		AmountOfModules: int(response.Course.AmountOfModules),
		AmountOfLessons: int(response.Course.AmountOfLessons),
	})
}

func (h *CourseHandler) ReadAllCourses(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	response, err := h.courseClient.ReadAllCourses(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	courses := make([]dto.CourseResponse, 0, len(response.Courses))
	for _, course := range response.Courses {
		courses = append(courses, dto.CourseResponse{
			Id:              course.Id,
			Title:           course.Title,
			Description:     course.Description,
			AccessType:      course.AccessType,
			PublishedAt:     course.PublishedAt,
			OwnerId:         course.OwnerId,
			AmountOfModules: int(course.AmountOfModules),
			AmountOfLessons: int(course.AmountOfLessons),
		})
	}

	return c.JSON(http.StatusOK, dto.CourseListResponse{
		Courses: courses,
	})
}

func (h *CourseHandler) ReadAllCoursesByOwnerId(c echo.Context) error {
	ownerId := c.Get("user_id").(string)
	if ownerId == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	response, err := h.courseClient.ReadAllCoursesByOwnerId(ctx, ownerId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	courses := make([]dto.CourseResponse, 0, len(response.Courses))
	for _, course := range response.Courses {
		courses = append(courses, dto.CourseResponse{
			Id:              course.Id,
			Title:           course.Title,
			Description:     course.Description,
			AccessType:      course.AccessType,
			PublishedAt:     course.PublishedAt,
			OwnerId:         course.OwnerId,
			AmountOfModules: int(course.AmountOfModules),
			AmountOfLessons: int(course.AmountOfLessons),
		})
	}

	return c.JSON(http.StatusOK, dto.CourseListResponse{
		Courses: courses,
	})
}

func (h *CourseHandler) UpdateCourse(c echo.Context) error {
	var request dto.CourseRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	response, err := h.courseClient.UpdateCourse(ctx, id, request.Title, request.Description, request.AccessType)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.CourseResponse{
		Id:          response.Course.Id,
		Title:       response.Course.Title,
		Description: response.Course.Description,
		AccessType:  response.Course.AccessType,
		PublishedAt: response.Course.PublishedAt,
		OwnerId:     response.Course.OwnerId,
	})
}

func (h *CourseHandler) UpdatePublishedAt(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	response, err := h.courseClient.UpdatePublishedAt(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.CourseResponse{
		Id:          response.Course.Id,
		Title:       response.Course.Title,
		Description: response.Course.Description,
		AccessType:  response.Course.AccessType,
		PublishedAt: response.Course.PublishedAt,
		OwnerId:     response.Course.OwnerId,
	})
}

func (h *CourseHandler) DeleteCourse(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	_, err := h.courseClient.DeleteCourse(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusNoContent, nil)
}
