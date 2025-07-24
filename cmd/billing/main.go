package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"zl0y.team/billing/internal/db"
	"zl0y.team/billing/internal/handlers"
	"zl0y.team/billing/internal/middleware"
)

func main() {
	log.Println("Starting Billing service...")

	pgDsn := os.Getenv("POSTGRES_DSN")
	mongoUri := os.Getenv("MONGO_URI")

	pg, err := db.NewPostgres(pgDsn)
	if err != nil {
		log.Fatalf("Postgres connection error: %v", err)
	}
	if err := db.MigratePostgres(pg); err != nil {
		log.Fatalf("Postgres migration error: %v", err)
	}

	mongoClient, err := db.NewMongo(mongoUri)
	if err != nil {
		log.Fatalf("MongoDB connection error: %v", err)
	}
	defer func() {
		if err := mongoClient.Disconnect(context.TODO()); err != nil {
			log.Printf("MongoDB disconnect error: %v", err)
		}
	}()

	r := gin.Default()

	// Передаём подключения к БД в middleware Gin
	r.Use(func(c *gin.Context) {
		c.Set("pg", pg)
		c.Set("mongo", mongoClient)
		c.Next()
	})

	r.POST("/api/auth/register", handlers.Register)
	r.POST("/api/auth/login", handlers.Login)

	api := r.Group("/api/user", middleware.JWTAuth())
	api.POST("/link-anonymous", handlers.LinkAnonymous)
	api.GET("/reports", handlers.UserReports)

	r.POST("/api/reports/:report_id/purchase", middleware.JWTAuth(), handlers.PurchaseReport)
	r.POST("/api/mock/create-report", handlers.MockCreateReport)
	// TODO: добавить маршруты

	log.Fatal(r.Run(":8080"))
}
