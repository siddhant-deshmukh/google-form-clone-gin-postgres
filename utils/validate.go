package utils

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var validate = validator.New()

func ValidateFieldWithStruct(s interface{}) (string, error) {

	if err := validate.Struct(s); err != nil {
		res_msg := "Invalid "
		if _, ok := err.(*validator.InvalidValidationError); ok {
			res_msg += "user data"
		} else {
			for _, err := range err.(validator.ValidationErrors) {
				res_msg += err.StructField() + " "
			}
		}

		return res_msg, err
	}

	return "", nil
}

func GetTokenKey() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Unable to get token Key")
	}

	token_key := os.Getenv("TOKEN_KEY")
	if token_key == "" {
		log.Fatal("Please add TOKEN_KEY in .env")
	}
	return token_key
}

func GetFieldFromUrl(c *gin.Context, field string) (uint, error) {
	id := c.Param(field)
	uID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad form id format",
			"err":     err,
		})
		return 0, err
	}
	return uint(uID), nil
}

type Result struct {
	AuthorID uint `json:"author_id"`
}

func GetFormAuthor(formId uint, db *gorm.DB) (uint, error) {
	var result Result
	response := db.Raw("SELECT author_id FROM forms WHERE id = ?", formId).Scan(&result)

	if response.Error != nil {
		return 0, response.Error
	}
	if response.RowsAffected == 0 {
		return 0, gorm.ErrRecordNotFound
	}
	return result.AuthorID, nil
}

func GetQuestionAuthor(quesId uint, db *gorm.DB) (uint, error) {
	var result Result
	response := db.Raw("SELECT author_id FROM questions WHERE id = ?", quesId).Scan(&result)

	if response.Error != nil {
		return 0, response.Error
	}
	if response.RowsAffected == 0 {
		return 0, gorm.ErrRecordNotFound
	}
	return result.AuthorID, nil
}

func ArrayToString(arr []string) string {
	if len(arr) <= 0 {
		return ""
	} else if len(arr) == 1 {
		return arr[0]
	} else if len(arr) == 2 {
		return arr[0] + ", " + arr[1]
	} else {

		str := arr[0] + ", "
		for i := 0; i < len(arr)-1; i++ {
			str += arr[i] + ", "
		}
		return str + arr[len(arr)-1]
	}
}
