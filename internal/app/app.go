package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TwiLightDM/diploma-gateway/internal/config"
	"github.com/TwiLightDM/diploma-gateway/internal/grpc/course-service"
	"github.com/TwiLightDM/diploma-gateway/internal/grpc/user-service"
	"github.com/TwiLightDM/diploma-gateway/internal/handlers"
	"github.com/TwiLightDM/diploma-gateway/internal/middlewares"
	"github.com/TwiLightDM/diploma-gateway/internal/services"
	"github.com/TwiLightDM/diploma-gateway/package/databases/postgres"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Run(cfg *config.Config) error {
	err := postgres.RunMigrations(cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.Username, cfg.Postgres.Password, cfg.Postgres.Database)
	if err != nil {
		return err
	}

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.RequestLogger())

	jwtService := services.NewJWTService(cfg.JWTSecret)

	authMiddleware := middlewares.AuthMiddleware(jwtService)

	userClient := user_service.NewUserClient(cfg.UserGRPCAddr)
	userHandler := handlers.NewUserHandler(userClient)
	groupHandler := handlers.NewGroupHandler(userClient)
	groupMemberHandler := handlers.NewGroupMemberHandler(userClient)

	courseClient := course_service.NewCourseClient(cfg.CourseGRPCAddr)
	courseHandler := handlers.NewCourseHandler(courseClient)
	moduleHandler := handlers.NewModuleHandler(courseClient)
	lessonHandler := handlers.NewLessonHandler(courseClient)
	lessonFileHandler := handlers.NewLessonFileHandler(courseClient)
	groupCourseHandler := handlers.NewGroupCourseHandler(courseClient)

	defer func() {
		log.Println("Closing gRPC connection to user service")
		_ = userClient.Close()
	}()

	registerRoutes(e, authMiddleware, userHandler, groupHandler, groupMemberHandler, courseHandler, moduleHandler, lessonHandler, lessonFileHandler, groupCourseHandler)

	server := &http.Server{
		Addr:    ":" + cfg.GatewayPort,
		Handler: e,
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	go func() {
		log.Printf("Gateway started on :%s", cfg.GatewayPort)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down gateway...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return err
	}

	log.Println("Gateway stopped gracefully")
	return nil
}

func registerRoutes(e *echo.Echo,
	authMiddleware echo.MiddlewareFunc,
	userHandler *handlers.UserHandler,
	groupHandler *handlers.GroupHandler,
	groupMemberHandler *handlers.GroupMemberHandler,
	courseHandler *handlers.CourseHandler,
	moduleHandler *handlers.ModuleHandler,
	lessonHandler *handlers.LessonHandler,
	lessonFileHandler *handlers.LessonFileHandler,
	groupCourseHandler *handlers.GroupCourseHandler,
) {
	public := e.Group("/auth")
	public.POST("/login", userHandler.Login)
	public.POST("/signup", userHandler.SignUp)
	public.POST("/refresh", userHandler.Refresh, authMiddleware)

	users := e.Group("/users", authMiddleware)
	users.GET("/me", userHandler.ReadSelf)
	users.GET("/all", userHandler.ReadAllUser)
	users.GET("", userHandler.ReadUser)
	users.PATCH("", userHandler.UpdateUser)
	users.PATCH("/role", userHandler.UpdateUserRole)
	users.PATCH("/password", userHandler.ChangePassword)
	users.GET("/courses", courseHandler.ReadAllCoursesByOwnerId)
	users.GET("/groups", groupHandler.ReadAllGroupsByOwnerId)
	users.GET("/:id/group-members", groupMemberHandler.ReadAllGroupMembersByUserId)

	groups := e.Group("/groups", authMiddleware)
	groups.POST("", groupHandler.CreateGroup)
	groups.GET("/:id", groupHandler.ReadGroup)
	groups.PATCH("/:id", groupHandler.UpdateGroup)
	groups.DELETE("/:id", groupHandler.DeleteGroup)
	groups.GET("/:id/group-members", groupMemberHandler.ReadAllGroupMembersByGroupId)
	groups.GET("/:id/group-courses", groupCourseHandler.ReadAllGroupCoursesByGroupId)

	groupMembers := e.Group("/group-members", authMiddleware)
	groupMembers.POST("", groupMemberHandler.CreateGroupMember)
	groupMembers.DELETE("/:id", groupMemberHandler.DeleteGroupMember)
	users.GET("/me", userHandler.ReadSelf)
	users.GET("/:id", userHandler.ReadUser)
	users.PATCH("", userHandler.UpdateUser)
	users.PATCH("/password", userHandler.ChangePassword)

	courses := e.Group("/courses", authMiddleware)
	courses.POST("", courseHandler.CreateCourse)
	courses.GET("", courseHandler.ReadAllCourses)
	courses.GET("/my", courseHandler.ReadAllCoursesByOwnerId)
	courses.GET("/:id", courseHandler.ReadCourse)
	courses.PATCH("/:id", courseHandler.UpdateCourse)
	courses.PATCH("/:id/publish", courseHandler.UpdatePublishedAt)
	courses.DELETE("/:id", courseHandler.DeleteCourse)
	courses.GET("/:id/group-courses", groupCourseHandler.ReadAllGroupCoursesByCourseId)

	modules := e.Group("/modules", authMiddleware)
	modules.POST("", moduleHandler.CreateModule)
	modules.GET("/courses/:course_id", moduleHandler.ReadAllModulesByCourseId)
	modules.GET("/:id", moduleHandler.ReadModule)
	modules.PATCH("/:id", moduleHandler.UpdateModule)
	modules.DELETE("/:id", moduleHandler.DeleteModule)

	lessons := e.Group("/lessons", authMiddleware)
	lessons.POST("", lessonHandler.CreateLesson)
	lessons.GET("/modules/:module_id", lessonHandler.ReadAllLessonsByCourseId)
	lessons.GET("/:id", lessonHandler.ReadLesson)
	lessons.PATCH("/:id", lessonHandler.UpdateLesson)
	lessons.DELETE("/:id", lessonHandler.DeleteLesson)

	lessonFiles := e.Group("/lessons/files", authMiddleware)
	lessonFiles.POST("", lessonFileHandler.UploadFile)
	lessonFiles.GET("/:lesson_id", lessonFileHandler.GetFiles)
	lessonFiles.DELETE("/:id", lessonFileHandler.DeleteFile)

	groupCourses := e.Group("/group-courses", authMiddleware)
	groupCourses.POST("", groupCourseHandler.CreateGroupCourse)
	groupCourses.DELETE("/:id", groupCourseHandler.DeleteGroupCourse)
}
