package db

import (
	"flag"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DBHost = flag.String("DBHost", "127.0.0.1", "数据库地址")
var DBUser = flag.String("DBUser", "postgres", "数据库账号名")
var DBPassword = flag.String("DBPassword", "postgres", "数据库密码")
var DBName = flag.String("DBName", "postgres", "要连接的库名")
var DBPort = flag.String("DBPort", "5432", "数据库端口")

func NewPG() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", *DBHost, *DBUser, *DBPassword, *DBName, *DBPort)
	pg, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败")
	}

	return pg
}
