package main
import (
	"fmt"
	"log"
	"net/http"

	"github.com/snake-game-api/database"
	"github.com/snake-game-api/controllers"

	"github.com/go-chi/chi"
	"github.com/rs/cors"
	_ "github.com/lib/pq"
)

func main() {
	defer func() {
        if re := recover(); re != nil {
			fmt.Println("Something went wrong")
        }
	}()
	defer cockroachdb.DB.Close()
	fmt.Println("Everything went correctly!")
	r := registerRoutes()
	log.Fatal(http.ListenAndServe(":3060", r))
}

func registerRoutes() http.Handler {
	r := chi.NewRouter()
	// Allow all origins with all standard methods with any header and credentials.
	r.Use(cors.AllowAll().Handler)
	r.Route("/records", func(r chi.Router) {
		r.Get("/", routes.GetAllRecords)                //GET 	 /records
		r.Get("/{username}", routes.GetRecord)       	//GET 	 /records/mich99
		r.Post("/", routes.AddRecord)                   //POST   /records
		r.Put("/{username}", routes.UpdateRecord)    	//PUT 	 /records/mich99
		r.Delete("/{username}", routes.DeleteRecord) 	//DELETE /records/mich99
	})
	return r
}