package configs

import (
	"context"
	"fmt"
	"go-jwt/models"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB        *gorm.DB
	Redis     *redis.Client
	RedisCtx  = context.Background()
)

func ConnectDB() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	// MySQL config
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	db.AutoMigrate(&models.User{}, &models.Product{}, &models.Category{}, &models.Transaction{}, &models.TransactionDetail{})
	DB = db
	log.Println("Connected to MySQL database")

	// Redis config
	redisHost := os.Getenv("REDIS_HOST")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	Redis = redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: redisPassword,
		DB:       0,
	})

	// Test Redis connection
	_, err = Redis.Ping(RedisCtx).Result()
	if err != nil {
		panic("failed to connect Redis: " + err.Error())
	}
	log.Println("Connected to Redis")
}
