/*
 * code in this file is to avoid hard-coding constant for 'currency' like "IDR, USD, and EUR".
 * The reason is to make easy if in the future this API want to support >100 types of currency,
 * and then there're also duplications of the currency because 'currency' can appear in many
 * different APIs.
 * custom validator is use to solve the problem.
 */

package api

import (
	"simple-bank-system/util"

	"github.com/go-playground/validator/v10"
)

// 'validator.Func' is a function that takes a 'validator.FieldLevel' interface as input and
// return "true" when validation succeeds. This is an interface that contains all informations
// and helper functions to validate a field.

/*var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(currency)
	}
	return false
}*/

func validCurrency(fl validator.FieldLevel) bool {
	currency := fl.Field().String()
	return util.IsSupportedCurrency(currency)
}
