package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/digitalhouse-content/go-fundamentals-web-users/internal/user"
	"github.com/digitalhouse-content/go-fundamentals-web-users/pkg/bootstrap"
	"github.com/digitalhouse-content/go-fundamentals-web-users/pkg/handler"
)

func main() {
	server := http.NewServeMux()

	db := bootstrap.NewDB()

	logger := bootstrap.NewLogger()
	repo := user.NewRepo(db, logger)
	service := user.NewService(logger, repo)

	ctx := context.Background()

	handler.NewUserHTTPServer(ctx, server, user.MakeEndpoints(ctx, service))

	fmt.Println("Server started at port 8080")
	log.Fatal(http.ListenAndServe(":8080", server))
}
