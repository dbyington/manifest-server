package main

import (
    "net/http"
    "os"

    "github.com/go-chi/chi/middleware"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/cors"

    "github.com/dbyington/manifest-server/handlers"
)

const (
    clientDir = "./client/build"
	serverPortEnv = "PORT"
)

func main() {
    r := chi.NewRouter()
    serverPort := os.Getenv(serverPortEnv)
    addr := ":" + serverPort

    r.Use(middleware.Logger)

    // Basic CORS
    // for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
    corsHandler := cors.Handler(cors.Options{
        // AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
        AllowedOrigins:   []string{"*"},
        // AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
        AllowedMethods:   []string{"GET", "OPTIONS"},
        AllowedHeaders:   []string{"X-PINGOTHER", "Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
        ExposedHeaders:   []string{"Link"},
        AllowCredentials: false,
        MaxAge:           300, // Maximum value not ignored by any of major browsers
    })
    r.Use(corsHandler)



    h := handlers.NewPkgHandler()
    r.Route("/manifest", func(r chi.Router) {
        r.Use(Cors)
        r.Get("/", h.ServeHTTP)
    })

    
    fs := http.FileServer(http.Dir(clientDir))
    r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
        if _, err := os.Stat(clientDir + r.RequestURI); os.IsNotExist(err) {
            http.StripPrefix(r.RequestURI, fs).ServeHTTP(w, r)
        } else {
            fs.ServeHTTP(w, r)
        }
    })

    if err := http.ListenAndServe(addr, r); err != nil {
        panic(err)
    }

}

func Cors(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "*")
        w.Header().Set("Access-Control-Allow-Headers", "*")
        w.Header().Set("Access-Control-Allow-Credentials", "true")

        if r.Method == "OPTIONS" {
            return
        }
        next.ServeHTTP(w, r)
    })
}
