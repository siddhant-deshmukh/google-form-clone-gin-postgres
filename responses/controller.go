package responses

import (
	"net/http"
	"os/user"

	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/siddhant-deshmukh/google-form-clone-gin-postgres/form"
	"github.com/siddhant-deshmukh/google-form-clone-gin-postgres/utils"
	"gorm.io/gorm"
)

func getResponses(c *gin.Context) {
	var responses []Response
	user_id := c.MustGet("user_id").(uint)

	form_id, err := utils.GetFieldFromUrl(c, "fId")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err":     err.Error(),
			"message": "Some error occured while getting author of form",
		})
		return
	}

	if author_id, err := utils.GetFormAuthor(form_id, db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err":     err.Error(),
			"message": "Some error occured while getting author of form",
		})
		return
	} else if author_id != user_id {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Not allowed",
		})
		return
	}

	result := db.Model(Response{}).Where("form_id = ?", form_id).Find(&responses)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err":     result.Error.Error(),
			"message": "Some error occured while getting responses",
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"responses": responses,
	})
}

func submitResponse(c *gin.Context) {
	var response Response

	user_id := c.MustGet("user_id").(uint)

	form_id, err := utils.GetFieldFromUrl(c, "fId")
	if err != nil {
		return
	}

	var form form.Form
	if result := db.First(&form, form_id); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err":     result.Error.Error(),
			"message": "Something went wrong",
		})
		return
	} else if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Not found",
		})
		return
	}

	var email string
	if result := db.Model(user.User{}).Select("email").Where("id = ?", user_id).Find(&email); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err":     result.Error.Error(),
			"message": "Something went wrong while getting email",
		})
		return
	} else if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Not found user",
		})
		return
	}

	var doesResponseExist bool = true
	if result := db.Model(Response{}).Where("form_id = ? AND user_email = ?", form_id, email).Find(&response); result.Error != nil {
		if err == gorm.ErrRecordNotFound {
			doesResponseExist = false
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"err":     err.Error(),
				"message": "Some error occured while checking if response exist",
			})
			doesResponseExist = false
			return
		}
	} else if result.RowsAffected == 0 {
		doesResponseExist = false
	}

	var answersRes NewResInput
	if err := c.BindJSON(&answersRes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Not correct format of input",
			"err":     err.Error(),
		})
		return
	}

	answersRes.UserEmail = email
	answersRes.FormID = form_id
	if err_msg, err := utils.ValidateFieldWithStruct(answersRes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err_msg,
			"err":     err.Error(),
		})
		return
	}

	var result *gorm.DB
	mapstructure.Decode(answersRes, &response)
	if doesResponseExist {
		result = db.Model(&response).Where("form_id = ? AND user_email = ?", form_id, email).Updates(response)
	} else {
		result = db.Create(&response)
	}
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err":     result.Error.Error(),
			"message": "Something went wrong final op",
			"does":    doesResponseExist,
		})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{
		"response": response,
	})
}
