package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog"

	"test-task-photo-booth/src/entities"
)

func Respond(w http.ResponseWriter, log *zerolog.Logger, data any) {
	if data == nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var buf bytes.Buffer
	if err := EncodeBody(&buf, data); err != nil {
		err = fmt.Errorf("encoding to buffer failed: %w", err)
		RespondErr(w, log, err, http.StatusInternalServerError)

		return
	}

	if _, err := buf.WriteTo(w); err != nil {
		err = fmt.Errorf("writing response failed: %w", err)
		RespondErr(w, log, err, http.StatusInternalServerError)

		return
	}
}

func RespondStatusOk(w http.ResponseWriter, log *zerolog.Logger) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var buf bytes.Buffer
	if err := EncodeBody(&buf, map[string]string{"status": "ok"}); err != nil {
		err = fmt.Errorf("encoding to buffer failed: %w", err)
		RespondErr(w, log, err, http.StatusInternalServerError)

		return
	}

	if _, err := buf.WriteTo(w); err != nil {
		err = fmt.Errorf("writing response failed: %w", err)
		RespondErr(w, log, err, http.StatusInternalServerError)

		return
	}
}

func RespondNativeJSON(w http.ResponseWriter, log *zerolog.Logger, data []byte) {
	if data == nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(data)
	if err != nil {
		err = fmt.Errorf("writing response failed: %w", err)
		RespondErr(w, log, err, http.StatusInternalServerError)

		return
	}
}

func RespondErr(w http.ResponseWriter, log *zerolog.Logger, err error, statusCode int) {
	log.Error().Err(err).Msg("responding with error")

	errResponse := entities.Error{
		Code:    statusCode,
		Message: err.Error(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response, err := json.Marshal(errResponse)
	if err != nil {
		log.Error().Err(fmt.Errorf("json.Marshal(errResponse) failed: %w", err)).Msg("responding with error")

		return
	}

	_, err = w.Write(response)
	if err != nil {
		log.Error().Err(fmt.Errorf("writing response failed: %w", err)).Msg("responding with error")

		return
	}
}

func RespondNativeErr(w http.ResponseWriter, log *zerolog.Logger, err error, statusCode int) {
	log.Error().Err(err).Msg("responding with error")

	http.Error(w, err.Error(), statusCode)
}

// DecodeBody reads data from a body and converts it to any
func DecodeBody(body io.Reader, data any) error {
	if err := json.NewDecoder(body).Decode(data); err != nil {
		return fmt.Errorf("decoding body failed: %w", err)
	}

	return nil
}

// EncodeBody writes data to a writer after converting it to JSON.
func EncodeBody(w io.Writer, data any) error {
	if err := json.NewEncoder(w).Encode(data); err != nil {
		return fmt.Errorf("encoding body failed: %w", err)
	}

	return nil
}
