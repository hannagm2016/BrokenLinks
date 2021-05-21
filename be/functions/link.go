package functions

import (
	"github.com/PuerkitoBio/goquery"
	m "main/be/models"
	"net/http"
	"regexp"
	"sync"
)

func  Links(res *http.Response, level int, wg *sync.WaitGroup, c chan m.Url, result *m.Results) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()

	var a = new(m.Url)
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
