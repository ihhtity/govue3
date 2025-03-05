package db

import (
	"context"
	"time"
)

type HealthStatus struct {
	MySQL   bool `json:"mysql"`
	Redis   bool `json:"redis"`
	MongoDB bool `json:"mongo"`
}

func HealthCheck() HealthStatus {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	status := HealthStatus{
		MySQL:   MySQL != nil,
		Redis:   Redis != nil,
		MongoDB: MongoDB != nil,
	}

	if status.MySQL {
		if err := MySQL.WithContext(ctx).Exec("SELECT 1").Error; err != nil {
			status.MySQL = false
		}
	}

	if status.Redis {
		if _, err := Redis.Ping(ctx).Result(); err != nil {
			status.Redis = false
		}
	}

	if status.MongoDB {
		if err := MongoDB.Ping(ctx, nil); err != nil {
			status.MongoDB = false
		}
	}

	return status
}
