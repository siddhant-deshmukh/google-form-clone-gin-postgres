package user

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/siddhant-deshmukh/google-form-clone-gin-postgres/form"
)

func getData(c *gin.Context) {
	user_id := c.MustGet("user_id").(uint)

	user, err := getUserById(user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "While getting user",
			"err":     err,
		})
		return
	}

	var forms []uint
	result := db.Model(form.Form{}).Select("id").Where("author_id = ?", user_id).Find(&forms)
	if result.Error != nil {
		fmt.Println(result.Error.Error())

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "While getting forms of user",
			"err":     result.Error,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":    user,
		"formIds": forms,
	})
}

func AuthUserMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth_token, err := c.Cookie("gf_clone_auth_token")
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Error in cookie",
				"err":     err,
			})
			c.Abort()
			return
		}

		token, err := jwt.Parse(auth_token, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(token_key), nil
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   err,
				"message": "Internal Server error authentication (Parse)",
			})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

			id := claims["_id"].(string)
			uid, err := strconv.ParseUint(id, 10, 32)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   err,
					"message": "Invalid token format",
				})
				c.Abort()
				return
			}

			// user, err := getUserById(uint(uid))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   err,
					"message": "Internal Server error (Claim)",
				})
				c.Abort()
				return
			}

			// c.Set("user", user)
			c.Set("user_id", uint(uid))
			c.Next()

		} else {
			c.JSON(http.StatusNotAcceptable, gin.H{
				"msg": "Invalid token",
			})
			c.Abort()
		}

	}
}

func getUserById(id uint) (User, error) {
	var user User
	if result := db.First(&user, id); result.Error != nil {
		return user, result.Error
	}
	return User{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}, nil
}
