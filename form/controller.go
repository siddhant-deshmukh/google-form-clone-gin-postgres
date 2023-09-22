package form

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/siddhant-deshmukh/google-form-clone-gin-postgres/utils"
)

func createForm(c *gin.Context) {
	var newForm Form
	user_id := c.MustGet("user_id").(uint)
	err := c.BindJSON(&newForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad form data format",
			"err":     err,
		})
		return
	}

	if res_msg, err := utils.ValidateFieldWithStruct(newForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": res_msg,
			"err":     err,
		})
		return
	}

	newForm.AuthorID = user_id

	result := db.Create(&newForm)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "While creating data",
			"err":     result.Error,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"form": newForm,
	})
}

func getFormById(c *gin.Context) {
	var isAuthor bool
	var form Form
	user_id := c.MustGet("user_id").(uint)

	formId, err := utils.GetFieldFromUrl(c, "id")
	if err == nil {
		return
	}

	result := db.First(&form, formId)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "While creating result",
			"err":     result.Error,
		})
		return
	}

	isAuthor = user_id == form.AuthorID

	if isAuthor {
		c.JSON(http.StatusCreated, gin.H{
			"form": form,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"form": gin.H{
				"ID":           form.ID,
				"Title":        form.Title,
				"Description":  form.Description,
				"Quiz_Setting": form.Quiz_Setting,
			},
		})
	}
}

func editForm(c *gin.Context) {
	var form EditForm
	user_id := c.MustGet("user_id").(uint)

	formId, err := utils.GetFieldFromUrl(c, "id")
	if err == nil {
		return
	}

	err = c.BindJSON(&form)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad form data format",
			"err":     err,
		})
		return
	}
	if res_msg, err := utils.ValidateFieldWithStruct(form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": res_msg,
			"err":     err,
		})
		return
	}

	jsonData, err := json.Marshal(form)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "While reading form",
			"err":     err,
		})
		return
	}

	var form_map map[string]interface{}
	err = json.Unmarshal(jsonData, &form_map)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "While reading form",
			"err":     err,
		})
		return
	}
	keys := []string{}
	for key, value := range form_map {
		if eq := reflect.DeepEqual(value, Quiz_Setting{}); !eq {
			keys = append(keys, key)
		}
	}
	fmt.Println(keys)
	fmt.Println()
	fmt.Println(form_map)

	result := db.Model(Form{}).Select(keys).Where("id = ? AND author_id = ?", formId, user_id).Updates(form_map)
	// result := db.Model(Form{}).Select("title").Where "send_res_copy": true("id = ? AND author_id = ?", formId, user_id).Updates(map[string]interface{}{"title": "Meow"})

	// result := db.First(&form, formId)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "While updating row",
			"err":     result.Error,
			"keys":    keys,
		})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "While updating row",
			"err":     result.Error,
		})
		return
	}
	// if user_id != form.AuthorID {
	// 	c.JSON(http.StatusForbidden, gin.H{
	// 		"message": "Permission Denied",
	// 	})
	// 	return
	// }
	c.JSON(http.StatusAccepted, gin.H{
		"message": "Updated",
		"form":    form,
	})
}

func deleteForm(c *gin.Context) {
	var form Form
	user_id := c.MustGet("user_id").(uint)

	formId, err := utils.GetFieldFromUrl(c, "id")
	if err == nil {
		return
	}

	result := db.First(&form, formId)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "While creating result",
			"err":     result.Error,
		})
		return
	}

	if user_id != form.AuthorID {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Permission Denied",
		})
		return
	}
}
