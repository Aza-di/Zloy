package handlers

import (
	"context"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"zl0y.team/billing/internal/models"
)

type CreateReportRequest struct {
	ClientGeneratedID string `json:"client_generated_id" binding:"required"`
}

func MockCreateReport(c *gin.Context) {
	var req CreateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	mongoClient := c.MustGet("mongo").(*mongo.Client)
	reports := mongoClient.Database("billing").Collection("reports")
	report := models.Report{
		ReportID:          randomReportID(),
		ClientGeneratedID: req.ClientGeneratedID,
		IsPurchased:       false,
		CreatedAt:         time.Now(),
	}
	_, err := reports.InsertOne(context.TODO(), report)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "mongo error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "created", "report_id": report.ReportID})
}

func randomReportID() string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, 12)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
