package routes

import (
	"fmt"
	"net/http"
	"encoding/json"
	"database/sql"

	"github.com/snake-game-api/models"
	"github.com/snake-game-api/database"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
)

type ErrResponse struct {
	Err				error `json:"-"` // low-level runtime error
	HTTPStatusCode	int   `json:"-"` // http response status code
	Message			string `json:"message"`          // user-level error message
	Detail			string `json:"detail"`          // user-level error description
}

func respondwithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		// TODO: What status code is used when there is an encoding error?
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Error while parsing response body",
			"detail": err.Error(),
		})
	}
}
// returns true if an error occurred while decoding the request body for POST or PUT
// methods. This means that the ResponseWriter was used. Returns false if no error occurred.
func decodeRequestBody(w http.ResponseWriter, r *http.Request, record *record.Record) bool {
	if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
		errResponse := ErrResponse {
			Err: err,
			HTTPStatusCode: http.StatusBadRequest,
			Message: "Invalid json request body",
			Detail: err.Error(),
		}
		respondwithJSON(w, errResponse.HTTPStatusCode, errResponse)
		return false
	}
	return true
}

func serverErrResponse(w http.ResponseWriter, err error) {
	errResponse := ErrResponse {
		Err: err,
		HTTPStatusCode: http.StatusInternalServerError,
		Message: "Database Server error",
		Detail: err.Error(),
	}
	respondwithJSON(w, errResponse.HTTPStatusCode, errResponse)
}

func notFoundErrResponse(w http.ResponseWriter, err error) {
	errResponse := ErrResponse {
		Err: err,
		HTTPStatusCode: http.StatusNotFound,
		Message: "Record not found",
		Detail: fmt.Sprintf("%s, %s", "This username doesn't have a record", err.Error()),
	}
	respondwithJSON(w, errResponse.HTTPStatusCode, errResponse)
}

func GetAllRecords(w http.ResponseWriter, r *http.Request) {
	records, err := record.SelectAllRows(cockroachdb.DB)
	if err != nil {
		serverErrResponse(w, err)
	} else {
		respondwithJSON(w, http.StatusOK, records)
	}
}

func GetRecord(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	switch record, err := record.SelectOneRow(cockroachdb.DB, username); err {
		case sql.ErrNoRows:
			notFoundErrResponse(w, err)
		case nil:
			respondwithJSON(w, http.StatusOK, record)
		default:
			serverErrResponse(w, err)
	}
}

func AddRecord(w http.ResponseWriter, r *http.Request) {
	var newRecord record.Record
	if !decodeRequestBody(w, r, &newRecord) {
		return
	}
	switch _, err := record.SelectOneRow(cockroachdb.DB, newRecord.Username); err {
		// username must be unique
		case sql.ErrNoRows:
			if newRecord.Username == "" || newRecord.BestScore <= 0 {
				errResponse := ErrResponse {
					Err: err,
					HTTPStatusCode: http.StatusBadRequest,
					Message: "Error in request body data",
					Detail: "username and bestScore fields can't be empty. bestScore must be greater than 0.",
				}
				respondwithJSON(w, errResponse.HTTPStatusCode, errResponse)
				return
			}
		case nil:
			errResponse := ErrResponse {
				Err: err,
				HTTPStatusCode: http.StatusBadRequest,
				Message: "Error in request body data",
				Detail: "username already exists. username must be unique.",
			}
			respondwithJSON(w, errResponse.HTTPStatusCode, errResponse)
			return
		default:
			serverErrResponse(w, err)
			return
	}

	if err := record.InsertRow(cockroachdb.DB, newRecord); err != nil {
		serverErrResponse(w, err)
	} else {
		respondwithJSON(w, http.StatusCreated, map[string]string{"message": "record created"})
	}
}

func UpdateRecord(w http.ResponseWriter, r *http.Request) {
	var updatedRecord record.Record
    username := chi.URLParam(r, "username")
	if !decodeRequestBody(w, r, &updatedRecord) {
		return
	}
	switch oldRecord, err := record.SelectOneRow(cockroachdb.DB, username); err {
		case sql.ErrNoRows:
			// username must exists and can't be empty
			notFoundErrResponse(w, err)
			return
		case nil:
			// The new bestScore must be greater than the old bestScore
			if oldRecord.BestScore >= updatedRecord.BestScore ||
			updatedRecord.BestScore <= 0 {
				errResponse := ErrResponse {
					Err: err,
					HTTPStatusCode: http.StatusBadRequest,
					Message: "Error in request body data",
					Detail: "The updated best score must be greater than the current best score",
				}
				respondwithJSON(w, errResponse.HTTPStatusCode, errResponse)
				return
			}
		default:
			serverErrResponse(w, err)
			return
	}
	// record is valid, execute the update
	if err := record.UpdateRow(cockroachdb.DB, username, updatedRecord); err != nil {
		serverErrResponse(w, err)
	} else {
		respondwithJSON(w, http.StatusOK, map[string]string{"message": "record updated"})
	}
}

func DeleteRecord(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	// Check if a record with this username exists
	if _, err := record.SelectOneRow(cockroachdb.DB, username); err == sql.ErrNoRows {
		notFoundErrResponse(w, err)
		return
	} else if err != nil {
		serverErrResponse(w, err)
		return
	}
	// This code can be in the previous missed else to also return the deleted record
	if err := record.DeleteRow(cockroachdb.DB, username); err != nil {
		serverErrResponse(w, err)
	} else {
		respondwithJSON(w, http.StatusOK, map[string]string{"message": "record deleted"})
	}
}