package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func writePdf(c *fiber.Ctx, data io.Reader) error {
	ctx:=c.UserContext()
	
	if data == nil {
		log.Ctx(ctx).Info().Msg("nothing to writeout: pdf data empty")
		return c.SendStatus(http.StatusNoContent)
	}

	c.Set("Content-type", "application/pdf")
	c.Set("Content-disposition", "attachment; filename=\"document.pdf\"")

	return c.SendStream(data)
}

func writeJson(ctx context.Context, w http.ResponseWriter, data any) error {
	if data == nil {
		log.Ctx(ctx).Info().Msg("nothing to writeout: json data empty")
		w.WriteHeader(http.StatusNoContent)
		return nil
	}

	err := json.NewEncoder(w).Encode(data)

	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("cant write json to body")
		return err
	}

	w.Header().Set("Content-Type", "application/json")

	return nil
}

func getValueFromForm(formValues map[string][]string, key string) (string, bool) {
	vals, ok := formValues[key]
	if !ok || len(vals) != 1 {
		return "", false
	}

	return vals[0], true
}
