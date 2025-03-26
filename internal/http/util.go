package httpsrv

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var errInvalidJSON = errors.New("invalid JSON")

func encodeJSON(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) error {
	jsonData, err := encodeJSON(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(jsonData)
	return err
}

func bindJSON(r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return errInvalidJSON
	}
	validate := validator.New()

	// TODO: customize error message
	err := validate.Struct(v)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		return errors
	}
	return nil
}

func errorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusInternalServerError
	if len(status) > 0 {
		statusCode = status[0]
	}

	errMsg := err.Error()

	resp := struct {
		Error string `json:"error"`
	}{
		Error: errMsg,
	}

	return writeJSON(w, statusCode, resp)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
