package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"

	"test-task-photo-booth/src/entities/dtos"
)

type PhotoUseCase interface {
	AddInQueue(photo *dtos.Photo) error
	Create(photo *dtos.Photo) error
	GetByID(id, quality string) (dtos.Photo, error)
	Delete(id string) error
}

type PhotoHandler struct {
	photoUseCase PhotoUseCase
	log          *zerolog.Logger
}

func NewPhotoHandler(photoUseCase PhotoUseCase, log *zerolog.Logger) PhotoHandler {
	return PhotoHandler{
		photoUseCase: photoUseCase,
		log:          log,
	}
}

type CreatePhotoRequest struct {
	Data string `json:"data" validate:"required"`
}

func (h PhotoHandler) Create(w http.ResponseWriter, r *http.Request) {
	requestData := new(CreatePhotoRequest)
	if err := DecodeBody(r.Body, requestData); err != nil {
		RespondErr(w, h.log, fmt.Errorf("DecodeBody() failed: %w", err), http.StatusBadRequest)

		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(requestData); err != nil {
		RespondErr(w, h.log, fmt.Errorf("validate.Struct() failed: %w", err), http.StatusBadRequest)

		return
	}

	photo := &dtos.Photo{
		Data: requestData.Data,
	}

	if err := h.photoUseCase.AddInQueue(photo); err != nil {
		RespondErr(w, h.log, fmt.Errorf("photoUseCase.AddInQueue(): %w", err), http.StatusInternalServerError)

		return
	}

	RespondStatusOk(w, h.log)
}

func (h PhotoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		RespondErr(w, h.log, fmt.Errorf("id is required"), http.StatusBadRequest)

		return
	}

	quality := r.URL.Query().Get("quality")
	if err := validateQuality(quality); err != nil {
		RespondErr(w, h.log, fmt.Errorf("validateQuality() failed: %w", err), http.StatusBadRequest)

		return
	}

	photo, err := h.photoUseCase.GetByID(id, quality)
	if err != nil {
		RespondErr(w, h.log, fmt.Errorf("photoUseCase.GetByID(): %w", err), http.StatusInternalServerError)

		return
	}

	Respond(w, h.log, photo)
}

func validateQuality(quality string) error {
	isValid := true

	switch quality {
	case "100", "75", "50", "25": // valida quality for photo
	default:
		isValid = false
	}

	if !isValid {
		return fmt.Errorf("invalid quality: %s", quality)
	}

	return nil
}

func (h PhotoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		RespondErr(w, h.log, fmt.Errorf("id is required"), http.StatusBadRequest)

		return
	}

	if err := h.photoUseCase.Delete(id); err != nil {
		RespondErr(w, h.log, fmt.Errorf("photoUseCase.Delete(): %w", err), http.StatusInternalServerError)

		return
	}

	RespondStatusOk(w, h.log)
}
