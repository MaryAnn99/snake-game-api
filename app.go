package main
import (
	"fmt"
	"log"
	"database/sql"
	"net/http"
	"encoding/json"

	"github.com/snake-game-api/record"
	"github.com/snake-game-api/cockroachdb"

	"github.com/go-chi/chi"
	"github.com/rs/cors"
	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	defer db.Close()
	db = cockroachdb.GetDatabaseConnection()
	fmt.Println("Everything went correctly!")
	r := registerRoutes()
	log.Fatal(http.ListenAndServe(":3060", r))
}

func respondwithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func registerRoutes() http.Handler {
	r := chi.NewRouter()
	// Allow all origins with all standard methods with any header and credentials.
	r.Use(cors.AllowAll().Handler)
	r.Route("/records", func(r chi.Router) {
		r.Get("/", getAllRecords)               //GET 	 /records
		r.Get("/{username}", getRecord)       	//GET 	 /records/mich99
		r.Post("/", addRecord)                  //POST   /records
		r.Put("/{username}", updateRecord)    	//PUT 	 /records/mich99
		r.Delete("/{username}", deleteRecord) 	//DELETE /records/mich99
	})
	return r
}

func getAllRecords(w http.ResponseWriter, r *http.Request) {
	records := record.SelectAllRows(db)
	respondwithJSON(w, http.StatusOK, records)
	defer func() {
        if re := recover(); re != nil {
			respondwithJSON(w, http.StatusInternalServerError, map[string]string{"message": "Something went wrong"})
        }
	}()
}

func getRecord(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	resultRecord := record.SelectOneRow(db, username)
	var zeroedRecord record.Record
	if resultRecord == zeroedRecord {
		respondwithJSON(w, http.StatusNotFound, map[string]string{"message": "record not found"})
	} else {
		respondwithJSON(w, http.StatusOK, resultRecord)
	}
	defer func() {
        if re := recover(); re != nil {
			respondwithJSON(w, http.StatusInternalServerError, map[string]string{"message": "Something went wrong"})
        }
	}()
}

func addRecord(w http.ResponseWriter, r *http.Request) {
	var newRecord record.Record
	json.NewDecoder(r.Body).Decode(&newRecord)  
	result := record.InsertRow(db, newRecord)
	if result == true {
		respondwithJSON(w, http.StatusCreated, map[string]string{"message": "successfully created"})
	} else {
		respondwithJSON(w, http.StatusBadRequest, map[string]string{"message": "something went wrong"})
	}
	defer func() {
        if re := recover(); re != nil {
			respondwithJSON(w, http.StatusInternalServerError, map[string]string{"message": "Something went wrong"})
        }
    }()
}

func updateRecord(w http.ResponseWriter, r *http.Request) {
	var newRecord record.Record
    username := chi.URLParam(r, "username")
    json.NewDecoder(r.Body).Decode(&newRecord)
	result := record.UpdateRow(db, username, newRecord)
	if result == true {
		respondwithJSON(w, http.StatusOK, map[string]string{"message": "update successfully"})
	} else {
		respondwithJSON(w, http.StatusBadRequest, map[string]string{"message": "something went wrong"})
	}
	defer func() {
        if re := recover(); re != nil {
			respondwithJSON(w, http.StatusInternalServerError, map[string]string{"message": "Something went wrong"})
        }
    }()
}

func deleteRecord(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	result := record.DeleteRow(db, username)
	if result == true {
		respondwithJSON(w, http.StatusOK, map[string]string{"message": "successfully deleted"})
	} else {
		respondwithJSON(w, http.StatusBadRequest, map[string]string{"message": "something went wrong"})
	}
	defer func() {
        if re := recover(); re != nil {
			respondwithJSON(w, http.StatusInternalServerError, map[string]string{"message": "Something went wrong"})
        }
    }()
}












