package main
import (
	"fmt"
	"log"
	"net/http"

	"github.com/snake-game-api/database"
	"github.com/snake-game-api/controllers"

	"github.com/go-chi/chi"
	"github.com/rs/cors"
	"github.com/go-chi/chi/middleware"
	_ "github.com/lib/pq"
)

func init() {
	defer func() {
        if re := recover(); re != nil {
			fmt.Println("An error occurred while trying to connect to the database.")
        }
	}()
	cockroachdb.ConnectToDB()
}

func main() {
	defer cockroachdb.DB.Close()
	fmt.Println("The server is running")
	r := registerRoutes()
	log.Fatal(http.ListenAndServe(":3060", r))
}

func registerRoutes() http.Handler {
	r := chi.NewRouter()

	c := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		Debug: false,
	})
	r.Use(c.Handler)
	r.Use(middleware.AllowContentType("application/json"))
	r.Route("/records", func(r chi.Router) {
		r.Get("/", routes.GetAllRecords)                //GET 	 /records
		r.Get("/{username}", routes.GetRecord)       	//GET 	 /records/mich99
		r.Post("/", routes.AddRecord)                   //POST   /records
		r.Put("/{username}", routes.UpdateRecord)    	//PUT 	 /records/mich99
		r.Delete("/{username}", routes.DeleteRecord) 	//DELETE /records/mich99
	})
	return r
}