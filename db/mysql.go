package db

import (
	"fmt"
	"govue3/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var MySQL *gorm.DB

func InitMySQL() {
	var err error
	MySQL, err = gorm.Open(mysql.Open(config.Cfg.Database.MySQL), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("MySQL连接失败: %v", err))
	}
	fmt.Println("MySQL连接成功")
}
