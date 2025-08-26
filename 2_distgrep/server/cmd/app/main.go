package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func grep(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "port to listen on")
	flag.Parse()

	e := echo.New()
	e.POST("/grep", grep)

	if err := e.Start(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatal(err)
	}
}
