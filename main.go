package main

import (
	"bufio"
	"io/ioutil"
	"log"
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
	wg     sync.WaitGroup
	contor int = 0
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

		// contor += len(ll)
		// fmt.Println(contor)
	}
}

func unique(input []string) []string {
	u := make([]string, 0, len(input))
	m := make(map[string]bool)

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}

	return u
}
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func main() {
	content, _ := ioutil.ReadFile(filename)
	lines := strings.Split(string(content), "\n")

	for {
		shuffle.Strings(lines)
		for _, i := range lines {
			request_group(format(i))

			func() {
				u, _ := readLines(results)
				ll := len(u)
				u = unique(u)
				file, err := os.OpenFile("test.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

				if err != nil {
					log.Fatalf("failed creating file: %s", err)
				}

				datawriter := bufio.NewWriter(file)

				for _, data := range u {
					_, _ = datawriter.WriteString(data + "\n")
				}

				datawriter.Flush()
				file.Close()
				println("had:", ll, " Have:", len(u), "scrapped:", ll-len(u))
			}()
		}
	}

}
