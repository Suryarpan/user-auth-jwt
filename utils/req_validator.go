package utils

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
)

const (
	notDecode    string = "could not decode data"
	badValidator string = "could not process request due to server error"
)

func ValidateReq[T any](w http.ResponseWriter, r *http.Request, vdt *validator.Validate, val *T) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(val)
	if err != nil {
		slog.Error("could not decode create user data", "error", err)
		EncodeError(w, http.StatusBadRequest, notDecode)
		return err
	}

	err = vdt.Struct(val)
	if err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			slog.Error("bad validator definition", "error", err)
			EncodeError(w, http.StatusInternalServerError, badValidator)
		} else {
			slog.Error("bad data provided", "error", err)
			EncodeValidationError(w, validationErrors)
		}
		return err
	}
	return nil
}
