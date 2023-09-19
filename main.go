package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/siddhant-deshmukh/google-form-clone-gin-postgres/form"
	"github.com/siddhant-deshmukh/google-form-clone-gin-postgres/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Unable to get Postgresql data source name (DSN)")
	}
	dsn := os.Getenv("PG_DATA_SOURCE_NAME")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Unable to connect to database")
	}

	user.SetUserTable(db)
	form.SetFormTable(db)

	router := gin.Default()

	userAuthRoutesGroup := router.Group("/")
	user.RegisterUserAuthRoutes(userAuthRoutesGroup)

	router.Run("localhost:8080")
}
