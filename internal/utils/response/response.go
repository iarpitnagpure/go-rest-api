package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Using json we can rename the response properties
type Response struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

const (
	StausOk     = "OK"
	StatusError = "ERROR"
)

func ResponseHandler(w http.ResponseWriter, status int, data interface{}) error {
	// Set Header by using ResponseWriter
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// Return JSON Response using json.NewEncoder, require ResponseWriter and data to convert into JSON
	return json.NewEncoder(w).Encode(data)
}

func ResponseErrorHandler(err error) Response {
	// Handle error scenario like 400, 404
	return Response{
		Error:  err.Error(),
		Status: StatusError,
	}
}

func ResponseValidationHandler(errs validator.ValidationErrors) Response {
	var errMessage []string

	// Validation Logic to check which field is missing in API payload
	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMessage = append(errMessage, fmt.Sprintf("field %s is required field", err.Field()))
		default:
			errMessage = append(errMessage, fmt.Sprintf("field %s is invalid", err.Field()))
		}
	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(errMessage, ", "),
	}
}
