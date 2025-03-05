package services

import (
	"govue3/db"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type BaseService struct {
	MySQL   *gorm.DB
	Redis   *redis.Client
	MongoDB *mongo.Client
}

func NewBaseService() *BaseService {
	return &BaseService{
		MySQL:   db.MySQL,
		Redis:   db.Redis,
		MongoDB: db.MongoDB,
	}
}
