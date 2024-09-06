package validators

import (
	"errors"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var validate *validator.Validate
var translator ut.Translator

// Initialise initialises the values of `validate` and `translator`.
// The internal variable naming is a bit confusing, but it's mostly
// copied from documentation
func Initialise() {

	localeTranslator := en.New()
	uni := ut.New(localeTranslator, localeTranslator)

	translator, _ = uni.GetTranslator("en")
	validate = validator.New()

	if err := en_translations.RegisterDefaultTranslations(validate, translator); err != nil {
		panic(err)
	}

	// Returns error messages with json tag names
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

}

// ValidateStruct validates a struct and returns a formatted string
func ValidateStruct(obj interface{}) error {
	err := validate.Struct(obj)
	if err != nil {
		return errors.New(errorsString(err.(validator.ValidationErrors)))
	}
	return nil
}

// errorsString returns a formatted string from the input validator.ValidationErrors
func errorsString(err validator.ValidationErrors) string {

	errsMap := err.Translate(translator)
	errsArray := []string{}
	for _, v := range errsMap {
		errsArray = append(errsArray, v)
	}
	return strings.Join(errsArray, "\n")
}
