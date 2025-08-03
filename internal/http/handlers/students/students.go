package students

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/iarpitnagpure/go-rest-api/internal/storage"
	"github.com/iarpitnagpure/go-rest-api/internal/types"
	"github.com/iarpitnagpure/go-rest-api/internal/utils/response"
)

func NewStudent(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var student types.Student

		// Add request body into student variable using json.NewDecoder, Accepts request body and struct pointer
		err := json.NewDecoder(r.Body).Decode(&student)
		// Check is there EOF error first, Since EOF is the error returned by Read when no more input is available.
		if errors.Is(err, io.EOF) {
			response.ResponseHandler(w, http.StatusBadRequest, response.ResponseErrorHandler(err))
			return
		}

		// Check any other error
		if err != nil {
			response.ResponseHandler(w, http.StatusBadRequest, response.ResponseErrorHandler(err))
			return
		}

		// request validation
		if err := validator.New().Struct(student); err != nil {
			// Type cast variables
			validateError := err.(validator.ValidationErrors)
			response.ResponseHandler(w, http.StatusBadRequest, response.ResponseValidationHandler(validateError))
			return
		}

		// Adding student entry in Database
		lastId, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)
		fmt.Println("lastId", lastId)
		if err != nil {
			response.ResponseHandler(w, http.StatusInternalServerError, err)
			return
		}

		response.ResponseHandler(w, http.StatusCreated, map[string]int64{"id": lastId})
	}
}
