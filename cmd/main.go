package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/digitalhouse-content/go-fundamentals-web-users/internal/user"
	"github.com/digitalhouse-content/go-fundamentals-web-users/pkg/bootstrap"
	"github.com/digitalhouse-content/go-fundamentals-web-users/pkg/handler"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	db, err := bootstrap.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	logger := bootstrap.NewLogger()
	repo := user.NewRepo(db, logger)
	service := user.NewService(logger, repo)

	ctx := context.Background()

	h := handler.NewUserHTTPServer(user.MakeEndpoints(ctx, service))

	port := os.Getenv("PORT")
	fmt.Println("Server started at port ", port)

	address := fmt.Sprintf("127.0.0.1:%s", port)

	srv := &http.Server{
		Handler: accessControl(h),
		Addr: address,
	}
	log.Fatal(srv.ListenAndServe())
}


func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS, HEAD, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept,Authorization,Cache-Control,Content-Type,DNT,If-Modified-Since,Keep-Alive,Origin,User-Agent,X-Requested-With")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}