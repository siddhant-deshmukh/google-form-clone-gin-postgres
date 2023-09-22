package question

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/siddhant-deshmukh/google-form-clone-gin-postgres/utils"
)

func newQuestion(c *gin.Context) {
	var question Question

	if err := c.BindJSON(&question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "bad format of question",
		})
		return
	}
	user_id := c.MustGet("user_id").(uint)
	question.AuthorID = user_id

	fmt.Println(question)

	if err_msg, err := utils.ValidateFieldWithStruct(question); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err_msg,
		})
		return
	}

	author_id, err := utils.GetFormAuthor(question.FormID, db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Internal Server error (checking form author)",
		})
		return
	} else if author_id != user_id {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   err.Error(),
			"message": "Not allowed",
		})
		return
	}

	results := db.Create(&question)
	if results.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something goes wrong while creating question",
			"error":   results.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Successful",
		"question": question,
	})
}

func editQuestion(c *gin.Context) {

}

func deleteQuestion(c *gin.Context) {

}

func getQuestions(c *gin.Context) {

}

// func CheckFormIdAndAuthorIdMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {

// 		token_key := utils.GetTokenKey()

// 		auth_token, err := c.Cookie("gf_clone_auth_token")
// 		if err != nil {
// 			c.JSON(http.StatusForbidden, gin.H{
// 				"message": "Error in cookie",
// 				"err":     err,
// 			})
// 			c.Abort()
// 			return
// 		}

// 		var form_id uint
// 		if form_id, err = utils.GetFieldFromUrl(c, "id"); err == nil {
// 			c.Abort()
// 			return
// 		}

// 		token, err := jwt.Parse(auth_token, func(token *jwt.Token) (interface{}, error) {
// 			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 			}
// 			return []byte(token_key), nil
// 		})
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{
// 				"error":   err,
// 				"message": "Internal Server error (Parse)",
// 			})
// 			c.Abort()
// 			return
// 		}

// 		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

// 			id := claims["_id"].(string)
// 			uid, err := strconv.ParseUint(id, 10, 32)
// 			if err != nil {
// 				c.JSON(http.StatusBadRequest, gin.H{
// 					"error":   err,
// 					"message": "Invalid token format",
// 				})
// 				c.Abort()
// 				return
// 			}

// 			// user, err := getUserById(uint(uid))
// 			if err != nil {
// 				c.JSON(http.StatusInternalServerError, gin.H{
// 					"error":   err,
// 					"message": "Internal Server error (Claim)",
// 				})
// 				c.Abort()
// 				return
// 			}

// 			// c.Set("user", user)

// 			author_id, err := utils.GetFormAuthor(form_id, db)
// 			if err != nil {
// 				c.JSON(http.StatusInternalServerError, gin.H{
// 					"error":   err,
// 					"message": "Internal Server error (checking form author)",
// 				})
// 				c.Abort()
// 				return
// 			}

// 			if author_id != uint(uid) {
// 				c.JSON(http.StatusForbidden, gin.H{
// 					"error":   err,
// 					"message": "Not allowed",
// 				})
// 				c.Abort()
// 				return
// 			}
// 			c.Set("user_id", uint(uid))
// 			c.Next()

// 		} else {
// 			c.JSON(http.StatusNotAcceptable, gin.H{
// 				"msg": "Invalid token",
// 			})
// 			c.Abort()
// 		}

// 	}4
// }
