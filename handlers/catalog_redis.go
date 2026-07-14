package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"goredis/services"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
)

type catalogHandlerRedis struct {
	catalogSrv  services.CatalogService
	redisClient *redis.Client
}

func NewCatalogHandlerRedis(catalogSrv services.CatalogService, redisClient *redis.Client) CatalogHandler {
	return catalogHandlerRedis{catalogSrv, redisClient}
}

func (h catalogHandlerRedis) GetProducts(c fiber.Ctx) error {

	key := "handler::GetProducts"

	//Redis GET
	if responseJson, err := h.redisClient.Get(context.Background(), key).Result(); err == nil {
		fmt.Println("from redis")
		c.Set("Content-Type", "application/json")
		return c.SendString(responseJson)
	}

	// Service
	products, err := h.catalogSrv.GetProducts()
	if err != nil {
		return err
	}

	response := fiber.Map{
		"status":   "ok",
		"products": products,
	}

	//Redis SET

	data, err := json.Marshal(response)
	if err == nil {
		h.redisClient.Set(context.Background(), key, string(data), time.Second*10)
	}

	fmt.Println("from database")
	return c.JSON(response)
}
