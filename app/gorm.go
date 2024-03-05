package app

import (
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

func NewDatabaseConnection(dsn string) *gorm.DB {

	/// SqlServer
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	/// Postgres
	//db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
	//	Logger: logger.Default.LogMode(logger.Info),
	//})
	if err != nil {
		log.Fatal(err)
	}

	connection, err := db.DB()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	connection.SetMaxIdleConns(10)
	connection.SetMaxOpenConns(100)
	connection.SetConnMaxLifetime(time.Second * time.Duration(300))
	return db
}
