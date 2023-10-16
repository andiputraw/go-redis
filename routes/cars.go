package routes

import (
	"andiputraw/belajar-redis/common"
	"andiputraw/belajar-redis/model"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var ctx = context.Background()

type CarRoutes struct {
	db *gorm.DB
	redis *redis.Client
}

func RegisterCarRoutes(r *gin.RouterGroup, db *gorm.DB, redis *redis.Client){
	Car := NewCarRoutes(db, redis)
	
	r.POST("/", Car.NewCar)
	r.GET("/", Car.GetCarFilter)
}

func NewCarRoutes(db *gorm.DB,redis *redis.Client) *CarRoutes{
	return &CarRoutes{
		db: db,
		redis: redis,
	}
}

func(r *CarRoutes) NewCar(c *gin.Context){
	car := new(model.Car)
	if err := c.BindJSON(car); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error" : err.Error()}) 
		return
	}

	if err := r.db.Model(&model.Car{}).Create(car).Error; err != nil {
		c.JSON(http.StatusInternalServerError, common.BadRequestResponse(err))
		return
	}

	c.JSON(http.StatusOK, common.SuccesResponse(common.SuccessMessage))
}

func(r *CarRoutes) GetCar(c *gin.Context) {
	var cars []*model.Car

	if err := r.db.Find(cars).Error; err != nil {
		c.JSON(http.StatusInternalServerError, common.BadRequestResponse(err))
		return	
	}
	c.JSON(http.StatusOK, common.SuccesWithData(common.SuccessMessage, cars))
	
}
type CarFilter struct {
	Color string  `form:"color"` 
	Price string  `form:"price"`
}

func(r *CarRoutes) GetCarFilter(c *gin.Context){
	var filter CarFilter
	_ = c.ShouldBindQuery(&filter)

	

	finding := new(model.Car)

	
	
	if(filter.Color != ""){
		finding.Color = filter.Color
	}
	
	if(filter.Price != ""){
		intPrice, err := strconv.Atoi(filter.Price)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.BadRequestResponse(err))
			return
		}
		finding.Price = intPrice
	}


	var cars []*model.Car

	if filter.Color == "" && filter.Price == "" {
		val, err  := r.redis.Get(ctx, "all_car_cache").Bytes()
		if err == redis.Nil {
			if err := r.db.Find(&cars).Error; err != nil {
				c.JSON(http.StatusBadRequest, common.BadRequestResponse(err))
				return		
			}

			jsonBytes, err := json.Marshal(cars)

			if err != nil {
				c.JSON(http.StatusBadRequest, common.BadRequestResponse(err))
				return
			}
			_ , err = r.redis.Set(ctx, "all_car_cache", jsonBytes, time.Minute * 10).Result()
			if err == redis.Nil {
				c.JSON(http.StatusInternalServerError, common.BadRequestResponse(err))
				return
			}
			
			c.JSON(http.StatusOK, common.SuccesWithData(common.SuccessMessage, cars))
			return
		}
		err = json.Unmarshal(val,&cars)
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.BadRequestResponse(err))
			return
		}
		c.JSON(http.StatusOK, common.SuccesWithData(common.SuccessMessage, cars))
		return
	}

	if err := r.db.Where(finding).Find(&cars).Error; err != nil {
		c.JSON(http.StatusBadRequest, common.BadRequestResponse(err))
		return
	}

	c.JSON(http.StatusOK, common.SuccesWithData(common.SuccessMessage, cars))

}

