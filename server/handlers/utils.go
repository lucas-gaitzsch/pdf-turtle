package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
)

func writePdf(ctx context.Context, w http.ResponseWriter, data io.Reader) error {
	if data == nil {
		log.Ctx(ctx).Info().Msg("nothing to writeout: pdf data empty")
		w.WriteHeader(http.StatusNoContent)
		return nil
	}

	w.Header().Set("Content-type", "application/pdf")
	w.Header().Set("Content-disposition", "attachment; filename=\"document.pdf\"")

	_, err := io.Copy(w, data)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("cant copy data to output stream")
		return err
	}

	return nil
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
