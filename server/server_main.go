package main

import (
	"cmkids/interfaces/basic"
	"cmkids/interfaces/routing"
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

	// dockerStatus := os.Getenv("CONTAINER_PREFIX")
	// if dockerStatus != "DOCKER" && dockerStatus != "LOCALHOST" {
	// 	sugarLogger.Fatalf("Wrong prefix: %s , should be DOCKER or LOCALHOST", dockerStatus)
	// }

	// repoUser := protoUser.NewUserClient(sessionUser)
	// repoAuth := protoAuth.NewAuthClient(sessionAuth)
	// repoPins := protoPins.NewPinsClient(sessionPins)
	// repoComments := protoComments.NewCommentsClient(sessionComments)
	// repoNotification := persistance.NewNotificationRepository(tarantoolConn)
	// repoChat := persistance.NewChatRepository(tarantoolConn)
	// cookieApp := application.NewCookieApp(repoAuth, 40, 10*time.Hour)
	// boardApp := application.NewBoardApp(repoPins)
	// s3App := application.NewS3App(sess, os.Getenv("BUCKET_NAME"))
	// userApp := application.NewUserApp(repoUser, boardApp)
	// authApp := application.NewAuthApp(repoAuth, userApp, cookieApp,
	// 	os.Getenv("VK_CLIENT_ID"), os.Getenv("VK_CLIENT_SECRET"))
	// pinApp := application.NewPinApp(repoPins, boardApp)
	// followApp := application.NewFollowApp(repoUser, pinApp)
	// commentApp := application.NewCommentApp(repoComments)
	// websocketApp := application.NewWebsocketApp(userApp)
	// notificationApp := application.NewNotificationApp(repoNotification, userApp, websocketApp)
	// chatApp := application.NewChatApp(repoChat, userApp, websocketApp)

	// boardInfo := board.NewBoardInfo(boardApp, logger)
	// authInfo := auth.NewAuthInfo(userApp, authApp, cookieApp, s3App, boardApp, websocketApp, logger)
	// profileInfo := profile.NewProfileInfo(userApp, authApp, cookieApp, followApp, s3App, notificationApp, logger)
	// followInfo := follow.NewFollowInfo(userApp, followApp, notificationApp, logger)
	// pinInfo := pin.NewPinInfo(pinApp, followApp, notificationApp, userApp, boardApp, s3App, logger,
	// 	pinEmailTemplate, os.Getenv("EMAIL_USERNAME"), os.Getenv("EMAIL_PASSWORD"))
	// commentsInfo := comment.NewCommentInfo(commentApp, pinApp, logger)
	// websocketInfo := websocket.NewWebsocketInfo(notificationApp, chatApp, websocketApp, os.Getenv("CSRF_ON") == "true", logger)
	// notificationInfo := notification.NewNotificationInfo(notificationApp, logger)
	// chatInfo := chat.NewChatnfo(chatApp, userApp, logger)

	basicHandler := basic.NewBasicHandler(logger)

	os.Setenv("CSRF_ON", "false")
	os.Setenv("HTTPS_ON", "true")

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
