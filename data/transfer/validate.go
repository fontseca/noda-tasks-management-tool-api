package transfer

import (
	"fmt"
	"noda/failure"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

func validate(s any) error {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if strings.Compare(name, "-") == 0 {
			return ""
		}
		return name
	})
	if err := validate.Struct(s); err != nil {
		errs := new(failure.AggregateDetails)
		for _, e := range err.(validator.ValidationErrors) {
			errs.Append(fmt.Sprintf("Validation for %q failed on: %s.",
				e.Field(), e.Tag()))
		}
		return errs
	}
	return nil
}
