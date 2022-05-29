package helpers

import (
	"fmt"
	"net/http"

	"github.com/yalagtyarzh/L0/internal/config"
)

var app *config.AppConfig

// ClientError writes error for client
func ClientError(w http.ResponseWriter, status int) {
	errstr := fmt.Sprintf("Client error with status of %d", status)
	app.Logger.Error(errstr)
	http.Error(w, http.StatusText(status), status)
}
