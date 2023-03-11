package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"toiler-graphql/auth"
	"toiler-graphql/cache"
	"toiler-graphql/database"
	"toiler-graphql/dataloaders"
	"toiler-graphql/graph"
	"toiler-graphql/graph/model"

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

func genOutputToken() (done bool) {
	userIdFlag := flag.Int("user_id", 0, "generate jwt for user_id")
	flag.Parse()

	if userIdFlag != nil && *userIdFlag != 0 {
		token, err := auth.GenToken(int32(*userIdFlag))
		if err != nil {
			panic(err)
		}
		fmt.Println(token)
		return true
	}
	return
}

func main() {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	if genOutputToken(){
		return
	}

	model.SetAvatarsPrefixPath(os.Getenv("AVATARS_S3_URL"))

	dbUrl := fmt.Sprintf("app_user:%s@/toiler_db?parseTime=true", os.Getenv("MYSQL_PASSWORD"))
	db, err := database.Open(dbUrl)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	repo := database.NewRepository(db)
	dl := dataloaders.NewRetriever()
	lruCache, err := cache.NewLRU(1024)
	if err != nil {
		panic(err)
	}

	redisDB, err := cache.NewRedisDB(
		os.Getenv("PUBSUB_ADDR"),
		os.Getenv("PUBSUB_PASSWORD"),
		0,
		lruCache)
	if err != nil {
		panic(err)
	}

	redisDB.ConsumeEvents()

	srv := handler.NewDefaultServer(
		graph.NewExecutableSchema(
			graph.Config{
				Resolvers: &graph.Resolver{
					Repository:  repo,
					Dataloaders: dl,
					Cache:       lruCache,
				},
			},
		),
	)

	dlMiddleware := dataloaders.Middleware(repo)

	playgroundHandler := playground.Handler("GraphQL playground", "/gql/query")

	m := middlewareStack{}
	m.Adapt(auth.AuthMiddleware(), dlMiddleware)

	mux := http.NewServeMux()
	mux.Handle("/gql/", m.Then(playgroundHandler))
	mux.Handle("/gql/query", m.Then(srv))
	mux.Handle("/gql/user_id", m.Then(UserIdEndPoint()))

	log.Printf("connect to http://localhost:%s/gql/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
