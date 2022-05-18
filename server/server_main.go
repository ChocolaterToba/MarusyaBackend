package main

import (
	"cmkids/adapter"
	authapp "cmkids/application/auth"
	quizapp "cmkids/application/quiz"
	marusiaHandler "cmkids/interfaces/marusia"
	quizhandler "cmkids/interfaces/quiz"
	"cmkids/interfaces/routing"
	settings "cmkids/models/settings"
	quizrepo "cmkids/repository/quiz"
	userrepo "cmkids/repository/user"
	"fmt"
	"net/http"
	"os"

	"github.com/rs/cors"
	"go.uber.org/zap"
)

func runServer(addr string) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	sugarLogger := logger.Sugar()

	config, err := settings.NewConfig("settings/values_local.yaml")
	if err != nil {
		sugarLogger.Fatal("Can not load config", zap.String("error", err.Error()))
	}

	conn, err := adapter.InitDB(config.Secrets.DBHost, config.Secrets.DBPassword)
	if err != nil {
		sugarLogger.Fatal("Can not init db connection", zap.String("error", err.Error()))
	}

	userRepo := userrepo.NewUserRepo(conn)
	quizRepo := quizrepo.NewQuizRepo(conn)

	authApp := authapp.NewAuthApp(userRepo, quizRepo, config, logger)
	quizApp := quizapp.NewQuizApp(authApp, quizRepo, config, logger)

	marusiaHandler := marusiaHandler.NewMarusiaHandler(quizApp, logger)
	quizHandler := quizhandler.NewQuizHandler(quizApp, logger)

	os.Setenv("HTTPS_ON", "false")

	r := routing.CreateRouter(marusiaHandler, quizHandler)

	allowedOrigins := make([]string, 0)
	switch os.Getenv("HTTPS_ON") {
	case "true":
		allowedOrigins = append(allowedOrigins, "https://127.0.0.1:8081", "https://51.250.14.186")
	case "false":
		allowedOrigins = append(allowedOrigins, "http://127.0.0.1:8081", "http://51.250.14.186")
	default:
		sugarLogger.Fatal("HTTPS_ON variable is not set")
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
	})

	handler := c.Handler(r)
	fmt.Printf("Starting server at localhost%s\n", addr)

	switch os.Getenv("HTTPS_ON") {
	case "true":
		sugarLogger.Fatal(http.ListenAndServeTLS(addr, "cert.pem", "key.pem", handler))
	case "false":
		sugarLogger.Fatal(http.ListenAndServe(addr, handler))
	}
}

func main() {
	runServer(":8080")
}
