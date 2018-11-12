package util

import (
	"strings"

	"github.com/jinzhu/inflection"
	"github.com/serenize/snaker"
)

type String struct {
	Camel struct {
		Plural   string
		Singular string
	}
	CamelLower struct {
		Plural   string
		Singular string
	}
	Snake struct {
		Plural   string
		Singular string
	}
}

func Inflect(in string) (out String) {
	out.Camel.Singular = inflection.Singular(snaker.SnakeToCamel(in))
	out.Camel.Plural = inflection.Plural(out.Camel.Singular)
	out.CamelLower.Singular = strings.ToLower(string(out.Camel.Singular[0])) + out.Camel.Singular[1:]
	out.CamelLower.Plural = strings.ToLower(string(out.Camel.Plural[0])) + out.Camel.Plural[1:]
	out.Snake.Singular = snaker.CamelToSnake(out.Camel.Singular)
	out.Snake.Plural = snaker.CamelToSnake(out.Camel.Plural)
	return
}
