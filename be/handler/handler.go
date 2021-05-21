package handler
import (
	"github.com/labstack/echo"
	f "main/be/functions"
	m "main/be/models"
	"net/http"
	"sync"
	"fmt"
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
	go f.Links(res, request.Level, wg, ch, Result)
	for c := range ch {
		go f.Ping(wg, c, ch, mu, Result)
	}
	fmt.Println(Result.BrokenLinks, len(Result.BrokenLinks), "BrokenLinks")
	fmt.Println(Result.AllLinks, len(Result.AllLinks), "AllLinks")

	return c.JSON(http.StatusOK, Result.BrokenLinks)
}