package cdmodel

import (
	"html/template"
)

type Page struct {
	Title			string
	Link			string
	Content			template.HTML
}
