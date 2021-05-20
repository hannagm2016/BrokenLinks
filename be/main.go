package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

type Url struct {
	Uri   string
	Level int
}
type Results struct {
	BrokenLinks []string
	AllLinks    []string
	MainURL     string
}

func (result *Results) ping(wg *sync.WaitGroup, url Url, c chan Url) {
	result.AllLinks = append(result.AllLinks, url.Uri)
	url.Level--
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("panic occurred:", err)
		}
	}()
	resp, err := http.Get(url.Uri)
	if err != nil {
		return
	}
	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		result.BrokenLinks = append(result.BrokenLinks, url.Uri)
		resp.Body.Close()

	} else if url.Level > 0 && strings.HasPrefix(url.Uri, result.MainURL) {
		go result.links(resp, url.Level, wg, c)
	} else {
		resp.Body.Close()
	}
	wg.Done()
}

func (result *Results) links(res *http.Response, level int, wg *sync.WaitGroup, c chan Url) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()

	var a = new(Url)
	doc, _ := goquery.NewDocumentFromReader(res.Body)
	doc.Find("a").Each(func(_ int, s *goquery.Selection) {
		a.Uri, _ = s.Attr("href")
		if a.Uri == "/" {
			return
		}
		a.Level = level
		matched, _ := regexp.MatchString(`^https?://.`, a.Uri)
		if !matched {
			a.Uri = result.MainURL + a.Uri
		}
		wg.Add(1)
		c <- *a
	})
	res.Body.Close()
	wg.Wait()
	close(c)
	return
}

func Check(c echo.Context) error {
	Result := new(Results)

	wg := new(sync.WaitGroup)

	request := new(Url)
	if err := c.Bind(&request); err != nil {
		return err
	}
	Result.MainURL = request.Uri
	res, err := http.Get(request.Uri)
	defer res.Body.Close()
	if err != nil {
		fmt.Println(err)
	}
	ch := make(chan Url, 100)
	go Result.links(res, request.Level, wg, ch)
	for c := range ch {
		go Result.ping(wg, c, ch)
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
