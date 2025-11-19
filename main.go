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

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	classRepo := repository.NewClassRepository(db)
	classService := service.NewClassService(classRepo)
	classHandler := handler.NewClassHandler(classService)

	krsRepo := repository.NewKRSRepository(db)
	krsService := service.NewKRSService(krsRepo)
	krsHandler := handler.NewKRSHandler(krsService)

	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	router.SetupRoutes(app, authHandler, classHandler, krsHandler, userHandler)

	if err := app.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
