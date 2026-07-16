package main

import (
	"log"

	"goredis/config"
	"goredis/handlers"
	"goredis/repositories"
	"goredis/services"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	db, err := initDatabase(cfg)
	if err != nil {
		log.Fatalf("database: %v", err)
	}

	redisClient := initRedis(cfg)
	// _ = redisClient

	productRepo := repositories.NewProductRepositoryDB(db) //ไม่มี redis
	// productRepo := repositories.NewProductRepositoryRedis(db, redisClient)  		//มี redis

	productservice := services.NewCatalogService(productRepo) //ไม่มี redis
	// productservice := services.NewCatalogServiceRedis(productRepo, redisClient)	//มี redis

	// productHandler := handlers.NewCatalogHandler(productservice)					//ไม่มี redis
	productHandler := handlers.NewCatalogHandlerRedis(productservice, redisClient) //มี redis

	app := fiber.New()

	app.Get("/products", productHandler.GetProducts)

	app.Listen(":" + cfg.AppPort)
}

func initDatabase(cfg *config.Config) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
}

func initRedis(cfg *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})
}
