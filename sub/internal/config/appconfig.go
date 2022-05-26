package config

import (
	"html/template"
	"net/http"

	"github.com/yalagtyarzh/L0/sub/pkg/logging"
)

// AppConfig holds main components of application
type AppConfig struct {
	InProduction  bool
	UseCache      bool
	TemplateCache map[string]*template.Template
	Logger        *logging.Logger
	Router        http.Handler
}
