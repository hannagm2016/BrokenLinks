package functions

import (
	"fmt"
	m "main/be/models"
	"net/http"
	"strings"
	"sync"
    "runtime"
)


func Ping(wg *sync.WaitGroup, url m.Url, c chan m.Url, mu *sync.Mutex, result *m.Results) {
	url.Level--
	mu.Lock()
	result.AllLinks = append(result.AllLinks, url.Uri)
	mu.Unlock()
	defer func() {
		if err := recover(); err != nil {
       	}
	}()
	resp, err := http.Get(url.Uri)
	if err != nil {
		wg.Done()
    	close(c)
		return
	}
	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		mu.Lock()
		result.BrokenLinks = append(result.BrokenLinks, url.Uri)
		resp.Body.Close()
    	mu.Unlock()
	} else if url.Level > 0 && strings.HasPrefix(url.Uri, result.MainURL) {
		go Links(resp, url.Level, wg, c, result)
	} else {
		resp.Body.Close()
	}
	fmt.Println("Ping", url, runtime.NumGoroutine())
	wg.Done()
}
