package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func Encode[T any](w http.ResponseWriter, status int, v T) error {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		return fmt.Errorf("could not encode to json: %w", err)
	}
	return nil
}

func EncodeError[T any](w http.ResponseWriter, status int, data T) error {
	var statusMessage string
	if 400 <= status && status <= 499 {
		statusMessage = "bad request received"
	} else if 500 <= status && status <= 599 {
		statusMessage = "server unable to process"
	} else {
		return fmt.Errorf("bad status code provided %d", status)
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)

	mssg := map[string]any{"message": statusMessage, "reason": data}
	err := json.NewEncoder(w).Encode(mssg)
	if err != nil {
		return fmt.Errorf("could not encode to json: %w", err)
	}
	return nil
}

func EncodeValidationError(w http.ResponseWriter, ves validator.ValidationErrors) {
	errorMssgs := make(map[string]map[string]any)
	for _, ve := range ves {
		errorMssgs[ve.Field()] = map[string]any{"value": ve.Value(), "error": ve.Error()}
	}
	EncodeError(w, http.StatusBadRequest, errorMssgs)
}
