package router

import (
	"course-planner-api/internal/handler"
	"os"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	jwt "github.com/golang-jwt/jwt/v4"
)

func SetupRoutes(app *fiber.App, authHandler *handler.AuthHandler, classHandler *handler.ClassHandler, krsHandler *handler.KRSHandler, userHandler *handler.UserHandler) {
	api := app.Group("/api")

	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	protected := api.Group("/me")
	protected.Use(jwtMiddleware())
	protected.Get("/", func(c *fiber.Ctx) error {
		// JWT middleware sudah mem-parse token, di sini hanya contoh respon
		return c.JSON(fiber.Map{"message": "authenticated endpoint"})
	})

	admin := api.Group("/admin")
	admin.Use(jwtMiddleware(), adminOnlyMiddleware())

	user := api.Group("/users")
    user.Get("/students", userHandler.ListStudents)

	krs := api.Group("/krs")
// ROUTE untuk mahasiswa, ubah dikit yaa untuk mastiin mana mahasiswa mana yang dosennya, soalnya nabrak2
	krsMahasiswa := krs.Group("/mahasiswa")
	krsMahasiswa.Use(jwtMiddleware(), roleOnlyMiddleware("mahasiswa"))
	krsMahasiswa.Get("/", krsHandler.GetTakenClasses)
	krsMahasiswa.Get("/available-classes", krsHandler.ListAvailableClasses)
	krsItems := krsMahasiswa.Group("/items")
	krsItems.Post("/", krsHandler.TakeClass)
	krsItems.Delete("/:classId", krsHandler.DropClass)
	krsItems.Patch("/:classId/request-cancellation", krsHandler.RequestCancellation) // Lihat kelas mahasiswa tapi masih orang lain selain dosen pa masi bisa liat keknya ni

	dosenKRS := api.Group("/krs/dosen")
	dosenKRS.Use(jwtMiddleware(), roleOnlyMiddleware("dosen"))
	dosenKRS.Delete("/mahasiswa/:mahasiswaId/items/:classId", krsHandler.DropClassByDosen)
	dosenKRS.Get("/mahasiswa/:id/classes", krsHandler.GetClassesByStudent)
	dosenKRS.Patch("/mahasiswa/:mahasiswaId/items/:itemId/verify", krsHandler.VerifyKRSItemByDosen)



	classes := admin.Group("/classes")
	classes.Get("/", classHandler.ListClasses)
	classes.Post("/", classHandler.CreateClass)
	classes.Get("/:id", classHandler.GetClass)
	classes.Patch("/:id", classHandler.UpdateClass)
	classes.Delete("/:id", classHandler.DeleteClass)
}

func jwtMiddleware() fiber.Handler {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "secret"
	}

	return jwtware.New(jwtware.Config{
		SigningKey: []byte(secret),
		ContextKey: "user",
	})
}

func adminOnlyMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := c.Locals("user").(*jwt.Token)
		if !ok || token == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		role, _ := claims["role"].(string)
		if role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
		}

		return c.Next()
	}
}

func roleOnlyMiddleware(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := c.Locals("user").(*jwt.Token)
		if !ok || token == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		role, _ := claims["role"].(string)
		if role != requiredRole {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden - required role: " + requiredRole})
		}

		return c.Next()
	}
}
