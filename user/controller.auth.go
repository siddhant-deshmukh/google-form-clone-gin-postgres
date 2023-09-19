package user

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

var token_key = getTokenKey()

func getTokenKey() string {
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

func userLogin(c *gin.Context) {
	var userData User
	var user User

	err := c.BindJSON(&userData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "enter user data in correct form",
			"err":     err,
		})
		return
	}

	result := db.Find(&User{}).Where("email = ?", userData.Email).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Some error occured",
			"err":     result.Error,
		})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userData.Password))
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"message": "Enter correct credentials", "err": err, "p": userData.Password, "c": user.Password})
		return
	}

	err = saveTokenString(c, strconv.FormatUint(uint64(user.ID), 10))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Some error occured while creating token",
			"error":   err,
		})
	}

	c.JSON(http.StatusCreated, gin.H{
		"user": User{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	})
}

func registerUser(c *gin.Context) {
	var newUserData User

	err := c.BindJSON(&newUserData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "enter user data in correct form",
			"err":     err,
		})
		return
	}

	var bytes []byte
	bytes, err = bcrypt.GenerateFromPassword([]byte(newUserData.Password), 14)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "While hashing password"})
		return
	}
	newUserData.Password = string(bytes)

	result := db.Create(&newUserData)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Some error occured",
			"err":     result.Error,
		})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusConflict, gin.H{
			"message": "Email already exist",
			"err":     result.Error,
		})
		return
	}

	err = saveTokenString(c, strconv.FormatUint(uint64(newUserData.ID), 10))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Some error occured while creating token",
			"error":   err,
		})
	}

	c.JSON(http.StatusCreated, gin.H{
		"user": User{
			ID:    newUserData.ID,
			Name:  newUserData.Name,
			Email: newUserData.Email,
		},
	})
}

func saveTokenString(c *gin.Context, ID string) error {

	// fmt.Println(ID)

	signing_key := []byte(token_key)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"_id": ID,
	})

	tokenString, err := token.SignedString(signing_key)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"message": "While creating token",
			"err":     err,
		})
		return err
	}

	c.SetCookie("gf_clone_auth_token", tokenString, 364000, "/", "http://www.localhost:8080.com", false, true)
	return nil
}
