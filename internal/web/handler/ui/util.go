package ui

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/beesbuddy/beesbuddy-worker/internal/app"
	"github.com/petaki/support-go/forms"
)

func serverError(ctx *app.Ctx, w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())

	if !ctx.Config.GetCfg().IsProd {
		http.Error(w, trace, http.StatusInternalServerError)
		return
	}

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func formError(ctx *app.Ctx, w http.ResponseWriter, err error) {
	var fe *forms.Error

	if errors.As(err, &fe) {
		http.Error(w, fe.Msg, fe.Status)
	} else {
		serverError(ctx, w, err)
	}
}
