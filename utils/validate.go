package utils

import "github.com/go-playground/validator/v10"

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
