package dto

type LessonFileRequest struct {
	ObjectName string `json:"object_name"`
	File       string `json:"file"`
	LessonId   string `json:"lesson_id"`
}

type LessonFileResponse struct {
	Id         string `json:"id"`
	ObjectName string `json:"object_name"`
	Url        string `json:"url"`
}

type LessonFileListResponse struct {
	LessonFiles []LessonFileResponse `json:"files"`
}
