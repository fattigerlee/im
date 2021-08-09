package db

import (
	"database/sql"
	"fmt"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"im/config"
	"im/pkg/logger"
)

var (
	DB       *gorm.DB
	RedisCli *redis.Client
)

// InitMysql 初始化MySQL
func InitMysql(cfg *config.Mysql) {
	logger.Info("init mysql")

	// 所有表默认以im_开头
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "im_" + defaultTableName
	}

	// 自动创建数据库
	autoCreateMysql(cfg)

	var err error
	if DB, err = gorm.Open("mysql", cfg.Dsn); err != nil {
		panic(err)
	}

	if cfg.Debug {
		DB = DB.Debug()
	}
	logger.Info("init mysql ok")
}

// InitRedis 初始化Redis
func InitRedis(cfg *config.Redis) {
	logger.Info("init redis")
	RedisCli = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		DB:       0,
		Password: cfg.Password,
	})

	if _, err := RedisCli.Ping().Result(); err != nil {
		panic(err)
	}
	logger.Info("init redis ok")
}

// InitByTest 初始化数据库配置，仅用在单元测试
func InitByTest() {
	logger.Info("init db")

	config.Init("config.yaml")
	logger.Init("im/db_test.log", logger.Console, "debug")

	InitMysql(config.GetMysql())
	InitRedis(config.GetRedis())
}

func autoCreateMysql(cfg *config.Mysql) {
	var (
		db  *sql.DB
		err error
	)

	if !cfg.AutoCreateDB {
		return
	}

	// 连接数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/mysql", cfg.Username, cfg.Password, cfg.Host, cfg.Port)
	if db, err = sql.Open("mysql", dsn); err != nil {
		logger.Error("open mysql failed, dsn: %s", dsn)
		return
	}
	defer func() {
		_ = db.Close()
	}()

	// 执行创建数据库命令
	query := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;", cfg.DBName)
	_, err = db.Exec(query)
	if err != nil {
		logger.Error("create database exec failed, err: %v", err)
		return
	}
	logger.Info("create database success")
}
