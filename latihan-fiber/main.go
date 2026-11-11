package main

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"latihan-fiber/app/repository"
	"latihan-fiber/config"
	"latihan-fiber/database"
)

func main() {
	// 1. Konfigurasi
	config.LoadEnv()

	// 2. Koneksi basis data
	pool, err := database.NewPool(context.Background())
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer pool.Close()

	// 3. Perakitan: pool -> repository -> handler
	studentRepository := repository.NewStudentRepository(pool)
	studentHandler := NewStudentHandler(studentRepository)

	// 4. Aplikasi Fiber
	app := fiber.New(fiber.Config{
		AppName: "latihan-fiber",
	})
	app.Use(requestid.New())
	app.Use(logger.New())
	app.Use(cors.New())

	api := app.Group("/api/v1")

	// /health ikut memeriksa kondisi basis data.
	// Bila PostgreSQL mati, Ping mengembalikan error dan /health
	// membalas 503 Service Unavailable. Inilah titik pengujian
	// untuk tangkapan layar "layanan PostgreSQL dimatikan".
	api.Get("/health", func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.UserContext(), 2*time.Second)
		defer cancel()

		if err := pool.Ping(ctx); err != nil {
			return fail(c, fiber.StatusServiceUnavailable,
				"database tidak dapat dihubungi")
		}
		return ok(c, "server dan database berjalan", nil)
	})

	s := api.Group("/students", requireJSON)
	s.Get("/", studentHandler.List)
	s.Get("/:id", studentHandler.Get)
	s.Post("/", studentHandler.Create)
	s.Put("/:id", studentHandler.Replace)
	s.Patch("/:id", studentHandler.Patch)
	s.Delete("/:id", studentHandler.Delete)

	port := config.GetEnv("APP_PORT", "3000")
	log.Fatal(app.Listen(":" + port))
}