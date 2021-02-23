package main

import (
	"log"
	"time"

	"github.com/bokimilinkovic/upp/geolocation"
	"github.com/bokimilinkovic/upp/handler"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"gopkg.in/tylerb/graceful.v1"
)

func main() {
	esPort := "http://localhost:9200"
	// Pdf reader
	geoLoc := &geolocation.Geo{Apikey: "537f18f94d1faf57db6fa565e43517cd"}
	elasticHandler := handler.NewElastic(esPort, geoLoc)
	uploader := handler.Uploader{Dir: "./books/"}

	e := echo.New()
	e.Use(echomiddleware.Logger(), echomiddleware.CORSWithConfig(
		echomiddleware.CORSConfig{
			AllowOrigins: []string{"*", "http://localhost:3000/"},
		},
	))

	e.Static("/static", "books")
	e.GET("/status", elasticHandler.CheckConnection)
	e.GET("/index", elasticHandler.CreateIndex)
	e.GET("/search", elasticHandler.Search)
	e.GET("/books", elasticHandler.GetAllBooks)
	e.POST("/book", elasticHandler.AddBook)
	e.POST("/upload", uploader.Upload)
	e.GET("/highlight", elasticHandler.HighlightSearch)
	e.GET("/users/index", elasticHandler.CreateUserIndex)
	e.GET("/users", elasticHandler.GetAllUsers)
	e.POST("/users", elasticHandler.AddUser)
	e.POST("/distance", elasticHandler.SearchDistance)

	e.Server.Addr = ":8080"
	graceful.DefaultLogger().Println("Application has successfully started at port: ", 8080)
	if err := graceful.ListenAndServe(e.Server, 5*time.Second); err != nil {
		log.Fatalf("Can not run server %s", err)
	}
}
