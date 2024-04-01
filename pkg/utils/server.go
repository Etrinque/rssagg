package utils

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"rssagg/internal"
	"rssagg/internal/database"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var root = internal.Path
var path = (root + "/rssagg/.env")
var Router = chi.NewRouter()
var corsMux = corsMiddleware(Router)

type ApiConfig struct {
	DB *database.Queries
}

func InitServer() {
	err := godotenv.Load(path)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT env not set")
	}
	fmt.Println(port)

	dbUrl := os.Getenv("CONNSTR")
	if dbUrl == "" {
		log.Fatal("error getting DB .env variable")
	}
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		fmt.Println("cannot connect to DB")
	}

	dbQueries := database.New(db)

	var apiCfg = ApiConfig{
		dbQueries,
	}

	Router.Get("/v1/readiness", handleReadiness)
	Router.Get("/v1/err", handleErr)
	Router.Post("/v1/users", apiCfg.handleUsersCreate)
	Router.Get("/v1/users", apiCfg.authMiddleware(apiCfg.handleUsersGet))
	Router.Post("/v1/feeds", apiCfg.authMiddleware(apiCfg.handleFeedCreate))
	Router.Get("/v1/feeds", apiCfg.handleGetAllFeeds)
	Router.Get("/v1/posts", apiCfg.authMiddleware(apiCfg.handleGetPostByUser))
	// Router.Get("/v1/posts", apiCfg.HandleGetFeedsFromUrl)
	Router.Post("/v1/feedfollow", apiCfg.authMiddleware(apiCfg.handleCreateFollowFeed))
	Router.Get("/v1/feedfollow", apiCfg.authMiddleware(apiCfg.handleGetFollowFeeds))
	Router.Delete("/v1/feedfollow/{feed_id}", apiCfg.authMiddleware(apiCfg.handleDeleteFollowFeed))

	server := &http.Server{
		Addr:         port,
		Handler:      corsMux,
		ReadTimeout:  500 * time.Second,
		WriteTimeout: 500 * time.Second,
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)

	go func() {
		sig := <-sigs
		fmt.Printf(" |--> Terminated with signal: %v <--|", sig)
		done <- true
		os.Exit(0)
	}()
	
	const collectionThreads = 10
	const collectionInterval = time.Minute
	go startScraper(dbQueries, collectionThreads, collectionInterval)

	fmt.Printf("listening on %v \n", port)
	log.Fatal(server.ListenAndServe())
	<-done
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			if strings.HasPrefix(origin, "http://") || strings.HasPrefix(origin, "https://") {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "*")
			w.Header().Set("Access-Control-Expose-Headers", "Link")
			w.Header().Set("Access-Control-Max-Age", "300")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
