package handlers

import (
	"context"
	"net/http"

	"strconv"

	"zl0y.team/billing/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

type LinkAnonymousRequest struct {
	ClientGeneratedID string `json:"client_generated_id" binding:"required"`
}

func LinkAnonymous(c *gin.Context) {
	var req LinkAnonymousRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	userID := c.GetInt("user_id")
	mongoClient := c.MustGet("mongo").(*mongo.Client)
	reports := mongoClient.Database("billing").Collection("reports")
	filter := bson.M{"client_generated_id": req.ClientGeneratedID, "user_id": bson.M{"$exists": false}}
	update := bson.M{"$set": bson.M{"user_id": userID}}
	res, err := reports.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "mongo error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"updated": res.ModifiedCount})
}

func UserReports(c *gin.Context) {
	userID := c.GetInt("user_id")
	mongoClient := c.MustGet("mongo").(*mongo.Client)
	reports := mongoClient.Database("billing").Collection("reports")

	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)

	filter := bson.M{"user_id": userID}
	findOpts := options.Find().SetLimit(limit).SetSkip(offset)
	cursor, err := reports.Find(c, filter, findOpts)
	if err != nil {
		c.JSON(500, gin.H{"error": "mongo error"})
		return
	}
	defer cursor.Close(c)

	var result []models.Report
	for cursor.Next(c) {
		var r models.Report
		if err := cursor.Decode(&r); err == nil {
			result = append(result, r)
		}
	}
	c.JSON(200, gin.H{"reports": result})
}

func PurchaseReport(c *gin.Context) {
	userID := c.GetInt("user_id")
	reportID := c.Param("report_id")
	pg := c.MustGet("pg").(*gorm.DB)
	mongoClient := c.MustGet("mongo").(*mongo.Client)
	reports := mongoClient.Database("billing").Collection("reports")

	// 1. Проверяем баланс пользователя
	var user models.User
	if err := pg.First(&user, userID).Error; err != nil {
		c.JSON(404, gin.H{"error": "user not found"})
		return
	}
	const price = 100 // цена отчёта в центах (пример)
	if user.Balance < price {
		c.JSON(402, gin.H{"error": "insufficient balance"})
		return
	}

	// 2. Имитируем транзакцию: списываем баланс и обновляем отчёт
	tx := pg.Begin()
	if err := tx.Model(&user).Update("balance", gorm.Expr("balance - ?", price)).Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "balance update error"})
		return
	}
	filter := bson.M{"report_id": reportID, "user_id": userID}
	update := bson.M{"$set": bson.M{"is_purchased": true}}
	res, err := reports.UpdateOne(c, filter, update, options.Update())
	if err != nil || res.ModifiedCount == 0 {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "mongo update error"})
		return
	}
	tx.Commit()
	c.JSON(200, gin.H{"status": "purchased"})
}
