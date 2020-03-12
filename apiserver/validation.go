package apiserver

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// use a single instance of Validate, as it caches struct info
var validate *validator.Validate

// use builtin url validation instead
// func validateURL(fl validator.FieldLevel) bool {
// 	return isURL(fl.Field().String())
// }

func validateURI(fl validator.FieldLevel) bool {
	re := regexp.MustCompile("^[A-Za-z0-9]([A-Za-z0-9-]*[A-Za-z0-9])?$")
	return re.MatchString(fl.Field().String())
}

func newValidator() *validator.Validate {
	validate = validator.New()
	validate.RegisterValidation("is-uri-path", validateURI)
	//validate.RegisterValidation("is-url", validateURL)
	return validate
}

// validateModel attempts to validate a model
// Returns the upstream validation error if any, as well as a slice of formatted error messages
func (s *server) validateModel(m interface{}) ([]string, error) {

	err := s.validator.Struct(m)
	if err != nil {

		switch err.(type) {
		case *validator.InvalidValidationError:
			// err = err.(*validator.InvalidValidationError)
			errMsgs := make([]string, 1)
			errMsgs[0] = "Internal Error"
			// s.logger.Debug().Msg(err.Error())
			s.logger.Error().Msg("validateModel: Internal validation failure")
			return errMsgs, err

		case validator.ValidationErrors:
			validationErrors := err.(validator.ValidationErrors)
			errMsgs := make([]string, len(validationErrors))

			for i, s := range validationErrors {
				var validationFailureReason string
				switch s.Tag() {
				case "max":
					validationFailureReason = " '" + s.Value().(string) + "' is too large or long"
				case "min":
					validationFailureReason = " '" + s.Value().(string) + "' is too small or short"
				case "is-uri":
					validationFailureReason = " '" + s.Value().(string) + "' is not a valid URI"
				case "url":
					validationFailureReason = " '" + s.Value().(string) + "' is not a valid URL"
				case "required":
					validationFailureReason = " is a required parameter"
				default:
					validationFailureReason = " '" + s.Value().(string) + "' has an unspecified error"
				}
				errMsgs[i] = s.Field() + validationFailureReason
			}
			return errMsgs, err
		}
	}
	return nil, nil
}
