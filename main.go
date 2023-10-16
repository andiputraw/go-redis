package main

import (
	"andiputraw/belajar-redis/model"
	"andiputraw/belajar-redis/routes"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {	
	r := gin.Default()

	dsn := "host=localhost user=postgres password=root dbname=car-redis port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	db.AutoMigrate(model.Car{})
	if err != nil {
		log.Fatalf("ERROR: Cannot connect to database %s", err)
	}

	opts := redis.Options{
		Addr: "127.0.0.1:6379",
		Password: "",
		DB: 0,
		Protocol: 3,
	}

	redis := redis.NewClient(&opts)
	
	car := r.Group("/car")

	routes.RegisterCarRoutes(car, db, redis)
	
	r.Run(":8080")
}