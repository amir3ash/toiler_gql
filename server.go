package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"toiler-graphql/auth"
	"toiler-graphql/database"
	"toiler-graphql/graph"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"
)

const defaultPort = "8080"

type Adapter func(http.Handler) http.Handler

type middlewareStack struct {
	middleWares []Adapter
}

func (s *middlewareStack) Then(h http.Handler) http.Handler {
	for _, adapter := range s.middleWares {
		h = adapter(h)
	}
	return h
}

// Setting up middlewares
func (s *middlewareStack) Adapt(adapters ...Adapter) {
	s.middleWares = adapters
}

func Notify(logger *log.Logger) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Println("before")
			defer logger.Println("after")
			h.ServeHTTP(w, r)
		})
	}
}

func UserIdEndPoint() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userId, err := auth.GetUserId(ctx)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return 
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "{\"user_id\": %d}", userId)
	})
	
}

func main() {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	dbUrl := fmt.Sprintf("app_user:%s@/toiler_db?parseTime=true", os.Getenv("MYSQL_PASSWORD"))
	db, err := database.Open(dbUrl)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	srv := handler.NewDefaultServer(
		graph.NewExecutableSchema(
			graph.Config{
				Resolvers: &graph.Resolver{
					Repository: *database.NewRepository(db),
				},
			},
		),
	)

	playgroundHandler := playground.Handler("GraphQL playground", "/gql/query")

	m := middlewareStack{}
	m.Adapt(auth.AuthMiddleware())

	mux := http.NewServeMux()
	mux.Handle("/gql/", m.Then(playgroundHandler))
	mux.Handle("/gql/query", m.Then(srv))
	mux.Handle("/gql/user_id", m.Then(UserIdEndPoint()))

	log.Printf("connect to http://localhost:%s/gql/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
