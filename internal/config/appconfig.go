package config

import (
	"html/template"

	"github.com/yalagtyarzh/L0/pkg/logging"
)

// AppConfig holds main components of application
type AppConfig struct {
	InProduction  bool
	UseCache      bool
	TemplateCache map[string]*template.Template
	Logger        *logging.Logger
}
