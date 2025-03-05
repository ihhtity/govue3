package db

import (
	"context"
	"fmt"
	"govue3/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDB *mongo.Client

func InitMongoDB() {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(config.Cfg.Database.MongoDB))
	if err != nil {
		panic(fmt.Sprintf("MongoDB连接失败: %v", err))
	}

	if err := client.Ping(context.Background(), nil); err != nil {
		panic(fmt.Sprintf("MongoDB Ping失败: %v", err))
	}
	MongoDB = client
	fmt.Println("MongoDB连接成功")
}
