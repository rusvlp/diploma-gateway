package handlers

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/TwiLightDM/diploma-gateway/internal/dto"
	"github.com/TwiLightDM/diploma-gateway/internal/grpc/course-service"
	"github.com/labstack/echo/v4"
)

type LessonFileHandler struct {
	courseClient *course_service.CourseClient
}

func NewLessonFileHandler(courseClient *course_service.CourseClient) *LessonFileHandler {
	return &LessonFileHandler{courseClient: courseClient}
}

func (h *LessonFileHandler) UploadFile(c echo.Context) error {
	file, header, err := c.Request().FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "file is required"})
	}
	defer file.Close()

	lessonId := c.FormValue("lesson_id")
	if lessonId == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "lesson_id is required"})
	}

	contentType := header.Header.Get("Content-Type")
	fileName := header.Filename

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 60*time.Second)
	defer cancel()

	response, err := h.courseClient.UploadFile(ctx, fileName, contentType, lessonId, header.Size, fileBytes)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, dto.LessonFileResponse{
		Id:         response.File.Id,
		ObjectName: response.File.ObjectName,
		Url:        response.File.Url,
	})
}

func (h *LessonFileHandler) GetFiles(c echo.Context) error {
	lessonId := c.Param("lesson_id")
	if lessonId == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	response, err := h.courseClient.GetLessonFiles(ctx, lessonId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	lessonFiles := make([]dto.LessonFileResponse, 0, len(response.Files))
	for _, lessonFile := range response.Files {
		lessonFiles = append(lessonFiles, dto.LessonFileResponse{
			Id:         lessonFile.Id,
			ObjectName: lessonFile.ObjectName,
			Url:        lessonFile.Url,
		})
	}

	return c.JSON(http.StatusOK, dto.LessonFileListResponse{
		LessonFiles: lessonFiles,
	})
}

func (h *LessonFileHandler) DeleteFile(c echo.Context) error {
	var request dto.LessonFileRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	_, err := h.courseClient.DeleteFile(ctx, id, request.ObjectName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusNoContent, nil)
}
