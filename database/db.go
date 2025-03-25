package database

import (
	"carrental/models"
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	username := viper.GetString("car.username")
	password := viper.GetString("car.password")
	dbHost := viper.GetString("car.db_host")
	dbPort := viper.GetInt("car.db_port")
	dbName := viper.GetString("car.db_name")

	dsn := username + ":" + password + "@tcp(" + dbHost + ":" + fmt.Sprintf("%d", dbPort) + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	db.AutoMigrate(&models.Customer{}, &models.Admin{}, &models.Vehicle{}, &models.Booking{}, &models.Payment{}, &models.Review{})

	DB = db
}

