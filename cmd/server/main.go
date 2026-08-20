package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/static"

	"mahrec/internal/db"
	"mahrec/internal/handlers"
	"mahrec/internal/middleware"
	"mahrec/internal/store"
)

func main() {
	dbPath := envOr("DB_PATH", "data/mahrec.db")
	addr := envOr("ADDR", ":3000")

	conn, err := db.Open(dbPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer conn.Close()

	clientStore := store.NewClientStore(conn)
	goalStore := store.NewGoalStore(conn)
	goalHandler := handlers.NewGoalHandler(goalStore)

	app := fiber.New()

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use("/static", static.New("./internal/web/static"))
	app.Use(middleware.ClientID(clientStore))

	app.Get("/healthz", handlers.Health)

	app.Get("/", goalHandler.Index)
	app.Post("/goals", goalHandler.Create)
	app.Post("/goals/:id/increment", goalHandler.Increment)
	app.Post("/goals/:id/reset", goalHandler.Reset)
	app.Delete("/goals/:id", goalHandler.Delete)

	log.Printf("listening on %s", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatal(err)
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
