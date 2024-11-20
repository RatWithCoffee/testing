package lab6

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

const maxBodySize = 1000000

const (
	unsupportedMediaTypeMsg = "Content-Type header is not application/json"
	badJson                 = "request body contains badly-formed JSON"
	emptyBody               = "request body must not be empty"
	bodyTooLarge            = "request body must not be larger than 1MB"
	severalJsonObjects      = "request body must only contain a single JSON object"
)

func (utils *HttpHandlerUtils) ReadJsonBody(w http.ResponseWriter, r *http.Request, value interface{}) error {
	ct := r.Header.Get("Content-Type")
	mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
	if mediaType != "application/json" {
		utils.StatusUnsupportedMediaType(w, unsupportedMediaTypeMsg)
		return errors.New(unsupportedMediaTypeMsg)
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(value)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			utils.BadRequest(w, msg)

		case errors.Is(err, io.ErrUnexpectedEOF):
			utils.BadRequest(w, badJson)

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			utils.BadRequest(w, msg)

		case errors.Is(err, io.EOF):
			msg := emptyBody
			utils.BadRequest(w, msg)

		case err.Error() == "http: request body too large":
			msg := bodyTooLarge
			utils.PayloadTooLarge(w, msg)

		default:
			utils.Log.Error("request body unmarshal error", zap.Error(err))
			utils.InternalServerError(w, "")
		}
		return err
	}

	err = validate.Struct(value)
	if err != nil {
		utils.BadRequest(w, err.Error())
		return err
	}

	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		utils.BadRequest(w, severalJsonObjects)
		return errors.New(severalJsonObjects)
	}
	return nil
}
