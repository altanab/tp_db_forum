package main

import (
	"context"
	"db_forum/models"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"os/signal"
	"time"

	"db_forum/server"
)

func main() {
	configString := "host=localhost user=postgres password=Qwerty123 dbname=forum sslmode=disable"
	dbpool, err := pgxpool.Connect(context.Background(), configString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	defer dbpool.Close()
	models.DBConn = dbpool

	e := echo.New()

	e.POST("/forum/create", server.CreateForum)
	e.GET("/forum/:slug/details", server.GetForumDetails)
	e.POST("/forum/:slug/create", server.CreateThread)
	e.GET("forum/:slug/users", server.GetForumUsers)
	e.GET("forum/:slug/threads", server.GetForumThreads)
	e.GET("post/:id/details", server.GetPostDetails)
	e.POST("post/:id/details", server.UpdatePost)
	e.POST("/service/clear", server.Clear)
	e.GET("/service/status", server.GetInfo)
	e.POST("/thread/:id/create", server.CreatePost)
	e.GET("/thread/:id/details", server.GetThreadDetails)
	e.POST("/thread/:id/details", server.UpdateThread)
	e.GET("/thread/:id/posts", server.GetThreadPosts)
	e.POST("/thread/:id/vote", server.VoteThread)
	e.POST("/user/:nickname/create", server.CreateUser)
	e.GET("/user/:nickname/profile", server.GetUserProfile)
	e.POST("/user/:nickname/profile", server.UpdateUser)


	// Start server
	go func() {
		if err := e.Start(":5000"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
