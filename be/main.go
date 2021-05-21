package main

import (
	h "main/be/handler"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.CORS())
	e.POST("/post", h.Check)
	e.Static("/", "../fe/dist/")
	e.Logger.Fatal(e.Start(":8000"))

}
