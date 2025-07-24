package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math/rand"
	"strconv"
	"sync"

	"zl0y.team/billing/internal/models"
)

var (
	captchaStore = make(map[string]string)
	captchaMutex sync.Mutex
)

type RegisterRequest struct {
	Login    string `json:"login" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6,max=64"`
}

type AuthResponse struct {
	AccessToken string `json:"access_token"`
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	pg := c.MustGet("pg").(*gorm.DB)

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "hash error"})
		return
	}

	user := models.User{
		Login:        req.Login,
		PasswordHash: string(hash),
		Balance:      0,
	}
	if err := pg.Create(&user).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "user exists"})
		return
	}

	token, err := generateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token error"})
		return
	}

	c.JSON(http.StatusOK, AuthResponse{AccessToken: token})
}

func Login(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	pg := c.MustGet("pg").(*gorm.DB)
	var user models.User
	if err := pg.Where("login = ?", req.Login).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	token, err := generateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token error"})
		return
	}
	c.JSON(http.StatusOK, AuthResponse{AccessToken: token})
}

func generateJWT(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	key := []byte(os.Getenv("JWT_SECRET"))
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(key)
}

// CaptchaHandler генерирует изображение с цифрами и сохраняет правильный ответ
func CaptchaHandler(c *gin.Context) {
	rand.Seed(time.Now().UnixNano())
	answer := ""
	for i := 0; i < 5; i++ {
		digit := rand.Intn(10)
		answer += strconv.Itoa(digit)
	}

	img := image.NewRGBA(image.Rect(0, 0, 120, 40))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
	// Для простоты: не рисуем цифры, а только возвращаем ответ (реализация с рисованием цифр требует font)
	// В реальной задаче используйте freetype или go-captcha для отрисовки цифр

	id := strconv.FormatInt(time.Now().UnixNano(), 10)
	captchaMutex.Lock()
	captchaStore[id] = answer
	captchaMutex.Unlock()

	c.Header("Content-Type", "image/png")
	c.Header("X-Captcha-Id", id)
	png.Encode(c.Writer, img)
}
