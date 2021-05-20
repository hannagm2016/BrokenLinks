package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"sync"
	"regexp"
	//"log"
	"time"
	"strings"
	"github.com/labstack/echo"
    	"github.com/labstack/echo/middleware"
	"os"
	//"runtime"
)

type Url struct {
	Uri   string
	Level int
}

var BrokenLinks []string
var AllLinks []string
var i,n, l int
var MainURL string
var LINKS_FILE   string= "links.txt"

func ping(wg *sync.WaitGroup, url Url, c chan Url) {
	//url:=<-c
	AllLinks = append(AllLinks, url.Uri)
	url.Level--
defer func() {
  if err := recover(); err != nil {
      fmt.Println("panic occurred:", err)
  }
  }()
	resp, err := http.Get(url.Uri)
	if err != nil {
		//fmt.Println(err, url)
		return
	}
	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		BrokenLinks = append(BrokenLinks, url.Uri)
		 links_file, _ := os.OpenFile(LINKS_FILE, os.O_APPEND|os.O_CREATE, 0666)
         	defer links_file.Close()
        links_file.WriteString(url.Uri + "\n")
		resp.Body.Close()

	} else if url.Level > 0 && strings.HasPrefix(url.Uri, MainURL) {
		//fmt.Println("to links", url, wg, i, l, c, runtime.NumGoroutine())
        		time.Sleep(2*time.Second)
        			ch := make(chan Url, 100)
                    	go links(resp, url.Level, wg, ch)
                    	for c := range ch {
                    		go ping(wg, c, ch)
                    	}
		//	go links(resp, url.Level, wg, c)
			} else {
		resp.Body.Close()
	}
	wg.Done()
}

func links(res *http.Response, level int, wg *sync.WaitGroup, c chan Url) {
defer func() {
  if err := recover(); err != nil {
      fmt.Println("panic occurred:", err)
  }
  }()

	var a = new(Url)
	doc, _ := goquery.NewDocumentFromReader(res.Body)
	doc.Find("a").Each(func(_ int, s *goquery.Selection) {

		i++
		a.Uri, _ = s.Attr("href")
		if a.Uri == "/" {
			return
		}
		a.Level = level
		matched, _ := regexp.MatchString(`^https?://.`, a.Uri)
		if !matched {
			a.Uri = MainURL + a.Uri
		}
		wg.Add(1)
		c <- *a
	})
	res.Body.Close()
//	wg.Done()
	wg.Wait()
	close(c)
	return
}
func Check(c echo.Context) error {
	wg := new(sync.WaitGroup)
	BrokenLinks=nil
AllLinks=nil
i=0
 links_file, _ := os.OpenFile(LINKS_FILE, os.O_APPEND|os.O_CREATE, 0666)
         	defer links_file.Close()
        links_file.WriteString("______________\n")

	request := new(Url)
	if err := c.Bind(&request); err != nil {
		return err
	}
	n = request.Level
	MainURL = request.Uri
	res, err := http.Get(request.Uri)
            	defer res.Body.Close()
            	if err != nil {
            		fmt.Println(err)
            	}
       ch := make(chan Url, 100)
       	go links(res, request.Level, wg, ch)
       	for c := range ch {
       		go ping(wg, c, ch)
       	}
	fmt.Println(BrokenLinks, i, len(BrokenLinks))
	fmt.Println(AllLinks, len(AllLinks),"***********")

	return c.JSON(http.StatusOK, BrokenLinks)
}
func main() {
//links_file, _ := os.OpenFile(LINKS_FILE, os.O_APPEND|os.O_CREATE, 0666)
e := echo.New()
	e.Use(middleware.CORS())
	e.POST("/post", Check)
	e.Static("/", "./dist/")
	e.Logger.Fatal(e.Start(":8000"))


}
