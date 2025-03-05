package utils

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// InitLogger 初始化日志记录器
func InitLogger() (*zap.Logger, error) {
	// 确保 logs 目录存在
	if err := os.MkdirAll("./logs", os.ModePerm); err != nil {
		return nil, fmt.Errorf("创建日志目录失败: %v", err)
	}

	// 自定义 lumberjack.Logger
	lumberjackLogger := &lumberjack.Logger{
		Filename:   getLogFileName(),
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
		LocalTime:  true,
	}

	// 配置日志写入文件
	fileEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	fileWriterSyncer := zapcore.AddSync(lumberjackLogger)

	// 配置日志写入控制台
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	consoleWriterSyncer := zapcore.AddSync(os.Stdout)

	// 配置日志级别
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, fileWriterSyncer, zap.InfoLevel),
		zapcore.NewCore(consoleEncoder, consoleWriterSyncer, zap.DebugLevel),
	)

	// 创建日志记录器
	logger := zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(logger)

	// 每天检查并更新日志文件名
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			lumberjackLogger.Filename = getLogFileName()
			lumberjackLogger.Rotate()
		}
	}()

	return logger, nil
}

// getLogFileName 返回当前日期的日志文件名
func getLogFileName() string {
	return fmt.Sprintf("./logs/%s.log", time.Now().Format("2006-01-02"))
}

// GracefulShutdown 优雅关闭服务
func GracefulShutdown(srv *http.Server) {
	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Info("正在关闭服务...")

	// 设置关闭超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Sugar().Fatalf("服务关闭失败: %v", err)
	}

	zap.L().Info("服务已关闭")
}
