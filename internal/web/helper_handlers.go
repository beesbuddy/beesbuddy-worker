package web

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/beesbuddy/beesbuddy-worker/internal/app"
	"github.com/beesbuddy/beesbuddy-worker/internal/web/dto"
	"github.com/gofiber/fiber/v2"
	"github.com/petaki/support-go/forms"
)

func ServerHTTPError(appCtx *app.Ctx, w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())

	if !appCtx.Config.GetCfg().IsProd {
		http.Error(w, trace, http.StatusInternalServerError)
		return
	}

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func ServeFiberError(appCtx *app.Ctx, fiberCtx *fiber.Ctx, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())

	if !appCtx.Config.GetCfg().IsProd {
		fiber.NewError(http.StatusInternalServerError, trace)
		return
	}

	fiber.NewError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}

func FormNetError(appCtx *app.Ctx, w http.ResponseWriter, err error) {
	var fe *forms.Error

	if errors.As(err, &fe) {
		http.Error(w, fe.Msg, fe.Status)
	} else {
		ServerHTTPError(appCtx, w, err)
	}
}

func AuthError(f *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return f.Status(fiber.StatusBadRequest).
			JSON(&dto.ResponseHTTP{
				Success: false,
				Data:    nil,
				Message: "Missing or malformed token",
			})
	}

	return f.Status(fiber.StatusUnauthorized).
		JSON(&dto.ResponseHTTP{
			Success: false,
			Data:    nil,
			Message: "Invalid or expired token",
		})
}
