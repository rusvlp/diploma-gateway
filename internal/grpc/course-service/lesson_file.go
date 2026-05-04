package course_service

import (
	"context"

	"github.com/TwiLightDM/diploma-course-service/proto/lessonfileservicepb"
)

func (c *CourseClient) UploadFile(ctx context.Context, fileName, contentType, lessonId string, size int64, file []byte) (*lessonfileservicepb.UploadFileResponse, error) {
	return c.lessonFile.UploadFile(ctx, &lessonfileservicepb.UploadFileRequest{
		FileName:    fileName,
		ContentType: contentType,
		File:        file,
		Size:        size,
		LessonId:    lessonId,
	})
}

func (c *CourseClient) GetLessonFiles(ctx context.Context, lessonId string) (*lessonfileservicepb.GetLessonFilesResponse, error) {
	return c.lessonFile.GetLessonFiles(ctx, &lessonfileservicepb.GetLessonFilesRequest{
		LessonId: lessonId,
	})
}

func (c *CourseClient) DeleteFile(ctx context.Context, id, objectName string) (*lessonfileservicepb.DeleteFileResponse, error) {
	return c.lessonFile.DeleteFile(ctx, &lessonfileservicepb.DeleteFileRequest{
		Id:         id,
		ObjectName: objectName,
	})
}
