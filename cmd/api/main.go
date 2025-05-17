package main

import (
	"fmt"
	"go-template/internal/common"
	internalhttp "go-template/internal/http"
	"go-template/internal/models"
	"go-template/internal/services"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	pg, err := services.ConnectPostgres()

	if err != nil {
		logger.Error(err.Error())
	}

	minio, err := services.ConnectMinio()
	if err != nil {
		logger.Error(err.Error())
	}

	common.App = &common.Application{
		Logger:     logger,
		Repository: &models.RepositoryModel{Pool: pg},
		Books:      &models.BookModel{Pool: pg, S3: minio, BucketName: common.GetEnv("BUCKET_NAME")},
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", common.GetEnv("PORT", "4000")),
		Handler:      internalhttp.Routes(),
		IdleTimeout:  time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	err = server.ListenAndServe()
	if err != nil {
		logger.Error(err.Error())
	}

	os.Exit(1)
}
