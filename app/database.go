package app

import (
	"akadia/internal/shared/helper"
	"akadia/seeder"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase() *gorm.DB {
	dsn := "postgres://postgres:12345678@localhost:5432/ujian2_rematch?sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("Database Not Connected " + err.Error())
	}

	if err := helper.CreateSchemas(db); err != nil {
		panic("Schema Failed Create" + err.Error())
	}

	if err := db.AutoMigrate(
	// &model.User{},
	// &model.Attendance{},
	// &model.Department{},
	// &model.Employee{},
	// &model.Leave{},
	// &model.Position{},
	// &model.Salary{},
	// &model.BlackListedToken{},
	); err != nil {
		panic("Auto Migrate Failed " + err.Error())
	}

	if err := seeder.Run(db); err != nil {
		panic("Seeder FAILED" + err.Error())
	}

	return db
}
