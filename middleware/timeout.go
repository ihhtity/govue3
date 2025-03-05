package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestTimeout 是一个中间件函数，用于设置请求的超时时间。
// 它接收一个 time.Duration 类型的参数 timeout，表示超时时间。
func RequestTimeout(timeout time.Duration) gin.HandlerFunc {
	// 返回一个 gin.HandlerFunc 类型的函数，这个函数是中间件的具体实现。
	return func(c *gin.Context) {
		// 使用 context.WithTimeout 创建一个带有超时时间的上下文。
		// c.Request.Context() 获取当前请求的上下文。
		// timeout 是超时时间。
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		// 在函数结束时调用 cancel 函数，释放资源。
		defer cancel()
		// 将带有超时时间的上下文设置到当前请求中。
		c.Request = c.Request.WithContext(ctx)
		// 调用 c.Next() 继续处理后续的中间件或处理器。
		c.Next()
	}
}
