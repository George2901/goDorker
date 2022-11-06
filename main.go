package main

import (
	"io/ioutil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/imroc/req/v3"
	"github.com/shogo82148/go-shuffle"
)

const (
	one      string = "https://www.bing.com/search?q="
	second   string = "&first="
	filename string = "texts/btc.txt"
	results  string = "texts/results.txt"
)

var (
	wg sync.WaitGroup
)

func format(q string) []string {
	ll := make([]string, 30)
	for i := 0; i <= 30; i++ {
		plm := one + url.QueryEscape(q) + second + strconv.Itoa(i*10+1)
		ll = append(ll, plm)
	}
	return ll
}

func write(text string) {
	f, err := os.OpenFile(results, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(text); err != nil {
		panic(err)
	}
}

func request_group(ll []string) {

	for _, i := range ll {
		wg.Add(1)
		go func(i string) {
			defer wg.Done()
			resp, err := req.Get(i)
			if err == nil {
				doc, err := goquery.NewDocumentFromReader(resp.Body)
				_ = err
				doc.Find("a").Each(func(i int, s *goquery.Selection) {
					band, _ := s.Attr("href")
					if strings.Contains(band, "http") {
						if !strings.Contains(band, "micro") {
							if strings.Contains(band, "?") {
								write(band + "\n")
							}
						}
					}
				})
			}
		}(i)
		wg.Wait()
	}
}

func main() {
	content, _ := ioutil.ReadFile(filename)
	lines := strings.Split(string(content), "\n")

	for {
		shuffle.Strings(lines)
		for _, i := range lines {
			request_group(format(i))

			// reader := bufio.NewReader(os.Stdin)
			// char, _, err := reader.ReadRune()
			// _ = char
			// _ = err
		}
	}

}
