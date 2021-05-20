package main

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"main/be/functions"
	m "main/be/models"
	"net/http"
	"sync"
)

func Check(c echo.Context) error {
	Result := new(m.Results)

	wg := new(sync.WaitGroup)
    mu:=new(sync.Mutex)
	request := new(m.Url)
	if err := c.Bind(&request); err != nil {
		return err
	}
	Result.MainURL = request.Uri
	res, err := http.Get(request.Uri)
	defer res.Body.Close()
	if err != nil {
		fmt.Println(err)
	}
	ch := make(chan m.Url, 100)
	go functions.Links(res, request.Level, wg, ch, Result)
	for c := range ch {
		go functions.Ping(wg, c, ch, mu, Result)
	}
	fmt.Println(Result.BrokenLinks, len(Result.BrokenLinks), "BrokenLinks")
	fmt.Println(Result.AllLinks, len(Result.AllLinks), "AllLinks")

	return c.JSON(http.StatusOK, Result.BrokenLinks)
}
func main() {
	e := echo.New()
	e.Use(middleware.CORS())
	e.POST("/post", Check)
	e.Static("/", "../fe/dist/")
	e.Logger.Fatal(e.Start(":8000"))

}
