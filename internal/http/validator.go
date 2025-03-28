package httpsrv

import "github.com/go-playground/validator/v10"

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	currency := fl.Field().String()
	switch currency {
	// TODO: consider moving to the config
	case "EUR", "USD", "CAD":
		return true
	}
	return false
}
