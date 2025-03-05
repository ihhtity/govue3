package middleware

import (
	"fmt"
	"govue3/controllers"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter 结构体定义了令牌桶限流器
type RateLimiter struct {
	limiter        *rate.Limiter
	userLimiterMap map[string]*userLimiter
	mu             sync.Mutex
}

// userLimiter 存储每个用户的限流信息
type userLimiter struct {
	lastBlockedTime time.Time
	limiter         *rate.Limiter
}

// NewRateLimiter 创建一个新的 RateLimiter 实例
func NewRateLimiter(rateLimit rate.Limit, burst int) *RateLimiter {
	return &RateLimiter{
		limiter:        rate.NewLimiter(rateLimit, burst),
		userLimiterMap: make(map[string]*userLimiter),
	}
}

// Middleware 返回一个 Gin 中间件函数，用于限流
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	base := &controllers.BaseController{} // 移动到这里
	return func(c *gin.Context) {
		// 获取用户的唯一标识，这里假设使用 IP 地址作为唯一标识
		userID := c.ClientIP()

		// 获取或创建用户的限流器
		rl.mu.Lock()
		ul, exists := rl.userLimiterMap[userID]
		if !exists {
			ul = &userLimiter{
				lastBlockedTime: time.Time{},
				limiter:         rate.NewLimiter(rl.limiter.Limit(), rl.limiter.Burst()),
			}
			rl.userLimiterMap[userID] = ul
		}
		rl.mu.Unlock()

		now := time.Now()
		if !ul.limiter.Allow() {
			if ul.lastBlockedTime.IsZero() || now.Sub(ul.lastBlockedTime) >= 3*time.Minute {
				ul.lastBlockedTime = now
				base.Error(c, "请求过于频繁，请稍后再试!")
			} else {
				remainingTime := 3*time.Minute - now.Sub(ul.lastBlockedTime)
				base.Error(c, fmt.Sprintf("请求过于频繁，请在 %v 后再试!", remainingTime))
			}
			c.Abort()
			return
		}
		c.Next()
	}
}
