package routes

import (
	"net/http"
	"encoding/json"

	"github.com/snake-game-api/models"
	"github.com/snake-game-api/database"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
)

func respondwithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func sendResponse(w http.ResponseWriter, result record.Result) {
	switch result.StatusCode {
		case 200:
			respondwithJSON(w, http.StatusOK, result.Data)
		case 201:
			respondwithJSON(w, http.StatusCreated, result.Data)
		case 400:
			respondwithJSON(w, http.StatusBadRequest, map[string]string{"message": "something went wrong"})
		case 404:
			respondwithJSON(w, http.StatusNotFound, map[string]string{"message": "record not found"})
		case 500:
			respondwithJSON(w, http.StatusInternalServerError, map[string]string{"message": "something went wrong"})
	}
}

func GetAllRecords(w http.ResponseWriter, r *http.Request) {
	records := record.SelectAllRows(cockroachdb.DB)
	sendResponse(w, records)
}

func GetRecord(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	resultRecord := record.SelectOneRow(cockroachdb.DB, username)
	sendResponse(w, resultRecord)
}

func AddRecord(w http.ResponseWriter, r *http.Request) {
	var newRecord record.Record
	json.NewDecoder(r.Body).Decode(&newRecord)  
	result := record.InsertRow(cockroachdb.DB, newRecord)
	sendResponse(w, result)
}

func UpdateRecord(w http.ResponseWriter, r *http.Request) {
	var newRecord record.Record
    username := chi.URLParam(r, "username")
    json.NewDecoder(r.Body).Decode(&newRecord)
	result := record.UpdateRow(cockroachdb.DB, username, newRecord)
	sendResponse(w, result)
}

func DeleteRecord(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	result := record.DeleteRow(cockroachdb.DB, username)
	sendResponse(w, result)
}