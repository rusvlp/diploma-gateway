package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/TwiLightDM/diploma-gateway/internal/dto"
	"github.com/TwiLightDM/diploma-gateway/internal/grpc/course-service"
	"github.com/labstack/echo/v4"
)

type GroupCourseHandler struct {
	courseClient *course_service.CourseClient
}

func NewGroupCourseHandler(courseClient *course_service.CourseClient) *GroupCourseHandler {
	return &GroupCourseHandler{courseClient: courseClient}
}

func (h *GroupCourseHandler) CreateGroupCourse(c echo.Context) error {
	var request dto.GroupCourseRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	response, err := h.courseClient.CreateGroupCourse(ctx, request.CourseId, request.GroupId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.GroupCourseResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.GroupCourseResponse{
		Id:       response.GroupCourse.Id,
		CourseId: response.GroupCourse.CourseId,
		GroupId:  response.GroupCourse.GroupId,
	})
}

func (h *GroupCourseHandler) ReadAllGroupCoursesByCourseId(c echo.Context) error {
	courseId := c.Param("id")
	if courseId == "" {
		return c.JSON(http.StatusBadRequest, dto.GroupCourseResponse{Error: "invalid request"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	response, err := h.courseClient.ReadAllGroupCoursesByCourseId(ctx, courseId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.GroupCourseResponse{Error: err.Error()})
	}

	groupCourses := make([]dto.GroupCourseResponse, 0, len(response.GroupCourses))
	for _, groupCourse := range response.GroupCourses {
		groupCourses = append(groupCourses, dto.GroupCourseResponse{
			Id:       groupCourse.Id,
			CourseId: groupCourse.CourseId,
			GroupId:  groupCourse.GroupId,
		})
	}

	return c.JSON(http.StatusOK, dto.GroupCourseListResponse{
		GroupCourses: groupCourses,
	})
}

func (h *GroupCourseHandler) ReadAllGroupCoursesByGroupId(c echo.Context) error {
	groupId := c.Param("id")
	if groupId == "" {
		return c.JSON(http.StatusBadRequest, dto.GroupCourseResponse{Error: "invalid request"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	response, err := h.courseClient.ReadAllGroupCoursesByGroupId(ctx, groupId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.GroupCourseResponse{Error: err.Error()})
	}

	groupCourses := make([]dto.GroupCourseResponse, 0, len(response.GroupCourses))
	for _, groupCourse := range response.GroupCourses {
		groupCourses = append(groupCourses, dto.GroupCourseResponse{
			Id:       groupCourse.Id,
			CourseId: groupCourse.CourseId,
			GroupId:  groupCourse.GroupId,
		})
	}

	return c.JSON(http.StatusOK, dto.GroupCourseListResponse{
		GroupCourses: groupCourses,
	})
}

func (h *GroupCourseHandler) DeleteGroupCourse(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, dto.GroupCourseResponse{Error: "invalid request"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	_, err := h.courseClient.DeleteGroupCourse(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.GroupCourseResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusNoContent, nil)
}
