package main

import (
	"cmkids/adapter"
	basicApp "cmkids/application/basic"
	basicHandler "cmkids/interfaces/basic"
	"cmkids/interfaces/routing"
	"cmkids/repository/basic"
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
	"os"

	"github.com/rs/cors"
	"go.uber.org/zap"
)

func runServer(addr string) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	sugarLogger := logger.Sugar()

	err := godotenv.Load("db.env")
	if err != nil {
		sugarLogger.Fatal("Could not load db.env file", zap.String("error", err.Error()))
	}
	host := os.Getenv("DB_HOST")
	pass := os.Getenv("DB_PASS")

	conn, err := adapter.InitDB(host, pass, logger)
	if err != nil {
		sugarLogger.Fatal("Can not init db connection", zap.String("error", err.Error()))
	}

	basicRep := basic.NewRepository(conn)
	basicApp := basicApp.NewBasicApp(basicRep, logger)
	basicHandler := basicHandler.NewBasicHandler(basicApp, logger)

	os.Setenv("CSRF_ON", "false")
	os.Setenv("HTTPS_ON", "false")

	r := routing.CreateRouter(basicHandler, os.Getenv("CSRF_ON") == "true", os.Getenv("HTTPS_ON") == "true")

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
