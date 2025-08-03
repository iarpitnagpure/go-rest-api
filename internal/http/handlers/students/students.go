package students

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

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

		// Send new student id as API response
		response.ResponseHandler(w, http.StatusCreated, map[string]int64{"id": lastId})
	}
}

func GetStudentById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get API params from request by using PathValue method
		id := r.PathValue("id")
		fmt.Println("id", id)

		// Need to convert string id into int64 using strconv.ParseInt
		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.ResponseHandler(w, http.StatusBadRequest, response.ResponseErrorHandler(err))
			return
		}

		// Get student from database
		student, err := storage.GetStudentById(intId)
		if err != nil {
			response.ResponseHandler(w, http.StatusInternalServerError, response.ResponseErrorHandler(err))
			return
		}

		// Send student as API response
		response.ResponseHandler(w, http.StatusOK, student)
	}
}
