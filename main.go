package main

import (
	"course-planner-api/config"
	"course-planner-api/internal/handler"
	"course-planner-api/internal/models"
	"course-planner-api/internal/repository"
	"course-planner-api/internal/router"
	"course-planner-api/internal/service"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env jika ada (abaikan kalau file tidak ditemukan)
	_ = godotenv.Load()

	db := config.LoadDatabase()
	if err := db.AutoMigrate(
		&models.User{},
		&models.Course{},
		&models.Room{},
		&models.Class{},
		&models.KRS{},
		&models.KRSItem{},
	); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	app := fiber.New()

	// Middleware logger: log setiap request ke terminal (mirip morgan di Express)
	app.Use(logger.New())

	// Serve dokumentasi OpenAPI via Swagger UI pada /docs
	app.Get("/docs", func(c *fiber.Ctx) error {
		return c.SendFile("docs/index.html")
	})
	app.Static("/docs", "./docs")

	// Auth
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	// Class
	classRepo := repository.NewClassRepository(db)
	classService := service.NewClassService(classRepo)
	classHandler := handler.NewClassHandler(classService)

	// KRS
	krsRepo := repository.NewKRSRepository(db)
	krsService := service.NewKRSService(krsRepo)
	krsHandler := handler.NewKRSHandler(krsService)
	dosenService := service.NewDosenPAService(krsRepo)
	dosenHandler := handler.NewDosenHandler(dosenService)

	// Course (Admin)
	courseRepo := repository.NewCourseRepository(db)
	courseService := service.NewCourseService(courseRepo)
	courseHandler := handler.NewCourseHandler(courseService)

	// Room (Admin)
	roomRepo := repository.NewRoomRepository(db)
	roomService := service.NewRoomService(roomRepo)
	roomHandler := handler.NewRoomHandler(roomService)

	// Dosen Management (Admin)
	dosenRepo := repository.NewDosenRepository(db)
	dosenMgmtService := service.NewDosenManagementService(dosenRepo)
	dosenMgmtHandler := handler.NewDosenManagementHandler(dosenMgmtService)

	router.SetupRoutes(
		app,
		authHandler,
		classHandler,
		krsHandler,
		dosenHandler,
		courseHandler,
		roomHandler,
		dosenMgmtHandler,
	)

	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
