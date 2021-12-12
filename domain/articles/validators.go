package articles

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var IsArticlePublishStatus validator.Func = func(fl validator.FieldLevel) bool {
	value := fmt.Sprintf("%s", fl.Field().Interface())
	return value == string(DRAFT) || value == string(PUBLISHED) || value == string(ARCHIVED)
}
