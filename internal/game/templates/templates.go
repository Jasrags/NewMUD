package templates

import (
	"fmt"
	"path/filepath"
	"text/template"
)

// TemplateCache stores parsed templates for reuse
var TemplateCache = map[string]*template.Template{}

// LoadTemplate dynamically loads and parses a template file.
func LoadTemplate(templateName string) (*template.Template, error) {
	// Check if the template is already cached
	if tmpl, ok := TemplateCache[templateName]; ok {
		return tmpl, nil
	}

	// Build the template path
	templatePath := filepath.Join("templates", templateName)

	// Parse the template file
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load template: %w", err)
	}

	// Cache the template for future use
	TemplateCache[templateName] = tmpl
	return tmpl, nil
}

// import (
// 	"log/slog"
// 	"os"
// 	"text/template"

// 	"github.com/Masterminds/sprig/v3"
// 	"github.com/muesli/termenv"
// 	"github.com/spf13/viper"
// )

// var (
// TemplateMgr = NewTemplateRenderer()
// )

// func NewTemplateRenderer() *template.Template {
// 	output := termenv.NewOutput(os.Stdout)

// 	for _, t := range viper.GetStringSlice("template_files") {
// 		slog.Debug("Template file",
// 			slog.String("file", t),
// 		)
// 	}

// 	tmpl, err := template.New("").
// 		Funcs(funcMap).
// 		Funcs(sprig.FuncMap()).
// 		Funcs(output.TemplateFuncs()).
// 		ParseFiles(viper.GetStringSlice("template_files")...)
// 	if err != nil {
// 		slog.Error("Error parsing templates",
// 			slog.Any("error", err))
// 	}

// 	slog.Debug("Templates defined",
// 		slog.Any("templates", tmpl.DefinedTemplates()))

// 	return tmpl
// }
