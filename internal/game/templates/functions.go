package templates

import (
	"text/template"

	"github.com/muesli/reflow/dedent"
	"github.com/muesli/reflow/padding"
	"github.com/muesli/reflow/wordwrap"
)

var (
	FuncMap = template.FuncMap{
		"wordWrap": func(max int, s string) string {
			return wordwrap.String(s, max)
		},
		"dedent": func(s string) string {
			return dedent.String(s)
		},
		"padding": func(width uint, s string) string {
			return padding.String(s, width)
		},
	}
)
