package utils

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

// InitConfig 初始化配置
func InitConfig() {
	//初始化配置文件
	viper.AddConfigPath("E:\\大二上\\数据库\\src\\Project\\config")
	viper.SetConfigName("app")
	viper.SetConfigType("yml")

	err := viper.ReadInConfig()
	if err != nil {
		// 打印更详细的错误信息
		fmt.Printf("读取配置文件时出错：%v\n", err)
		fmt.Printf("配置文件路径：%s\n", viper.ConfigFileUsed())
	} else {
		fmt.Println("配置文件已初始化")
	}
}

// InitMysql 初始化 MySQL 数据库连接
func InitMysql() error {
	// 自定义日志模板 打印 SQL 语句
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // 级别
			Colorful:      true,        // 彩色
		},
	)

	// 初始化数据库
	db, err := gorm.Open(mysql.Open(viper.GetString("mysql.dsn")), &gorm.Config{Logger: newLogger})
	if err != nil {
		return err
	}

	DB = db
	fmt.Println("MySQL 数据库已初始化")
	return nil
}
