package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/siddhant-deshmukh/google-form-clone-gin-postgres/form"
	"github.com/siddhant-deshmukh/google-form-clone-gin-postgres/question"
	"github.com/siddhant-deshmukh/google-form-clone-gin-postgres/responses"
	"github.com/siddhant-deshmukh/google-form-clone-gin-postgres/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Unable to get Postgresql data source name (DSN)")
	}
	dsn := os.Getenv("PG_DATA_SOURCE_NAME")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Unable to connect to database")
	}

	if _, ok := binding.Validator.Engine().(*validator.Validate); ok {
		fmt.Println("Okay")
		// v.RegisterValidation("bookabledate", bookableDate)
	}

	user.SetUserTable(db)
	question.SetQuestionTable(db)
	form.SetFormTable(db)
	responses.SetResponseTable(db)

	router := gin.Default()

	userAuthRoutesGroup := router.Group("/")
	user.RegisterUserAuthRoutes(userAuthRoutesGroup)

	userRoutesGroup := router.Group("/u")
	userRoutesGroup.Use(user.AuthUserMiddleWare())
	user.RegisterUserRoutes(userRoutesGroup)

	formRoutesGroup := router.Group("/f")
	formRoutesGroup.Use(user.AuthUserMiddleWare())
	form.RegisterFormRoutes(formRoutesGroup)

	questionRoutesGroup := router.Group("/q")
	questionRoutesGroup.Use(user.AuthUserMiddleWare())
	question.RegisterQuestionRoutes(questionRoutesGroup)

	resRoutesGroup := router.Group("/r")
	resRoutesGroup.Use(user.AuthUserMiddleWare())
	responses.RegisterResponseRoutes(resRoutesGroup)

	// srv := handler.NewDefaultServer(generated.NewExecutableSchema())

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))

	router.Run("localhost:8080")
}
