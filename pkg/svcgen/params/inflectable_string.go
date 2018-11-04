package params

import (
	"strings"

	"github.com/jinzhu/inflection"
	"github.com/serenize/snaker"
)

type inflectableString struct {
	pluralCamel        string
	pluralCamelLower   string
	pluralSnake        string
	singularCamel      string
	singularCamelLower string
	singularSnake      string
}

func inflect(name string) inflectableString {
	infl := inflectableString{
		pluralCamel:   inflection.Plural(snaker.SnakeToCamel(name)),
		singularCamel: inflection.Singular(snaker.SnakeToCamel(name)),
	}
	infl.pluralCamelLower = strings.ToLower(string(infl.pluralCamel[0])) + infl.pluralCamel[1:]
	infl.pluralSnake = snaker.CamelToSnake(infl.pluralCamel)
	infl.singularCamelLower = strings.ToLower(string(infl.singularCamel[0])) + infl.singularCamel[1:]
	infl.singularSnake = snaker.CamelToSnake(infl.singularCamel)
	return infl
}
