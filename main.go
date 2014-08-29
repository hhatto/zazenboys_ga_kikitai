package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	var doc *goquery.Document
	var err error
	quit := make(chan string)
	count := 0

	baseUrl := "http://mukaishutoku.com/download.html"
	if len(os.Args) == 2 {
		baseUrl = os.Args[1]
	}

	parsedUrl, err := url.Parse(baseUrl)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		m := []string{"|", "/", "-", "\\", "|"}
		for {
			for _, v := range m {
				fmt.Print(v)
				time.Sleep(120 * time.Millisecond)
				fmt.Print("\x1b[1D")
			}
		}
	}()

	if doc, err = goquery.NewDocument(baseUrl); err != nil {
		log.Fatal(err)
	}

	doc.Find("a").Each(func(_ int, s *goquery.Selection) {
		mp3Url, _ := s.Attr("href")
		var downloadUrl string
		if match, err := regexp.MatchString("http://.+.mp3", mp3Url); err == nil && match {
			downloadUrl = mp3Url
			//fmt.Println(mp3Url)
		} else if match, err := regexp.MatchString(".+.mp3", mp3Url); err == nil && match {
			downloadUrl = fmt.Sprintf("%s://%s/%s", parsedUrl.Scheme, parsedUrl.Host, mp3Url)
			//fmt.Println(downloadUrl)
		} else {
			return
		}
		count++

		go func(dlUrl string) {
			res, err := http.Get(dlUrl)
			if err != nil {
				log.Fatal(err)
			}
			defer res.Body.Close()
			contents, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.Fatal(err)
			}

			mp3ParsedUrl, err := url.Parse(dlUrl)
			if err != nil {
				log.Fatal(err)
			}
			_, target := path.Split(mp3ParsedUrl.Path)
			ioutil.WriteFile(target, contents, 0755)
			quit <- target
		}(downloadUrl)
	})

	var ret []string
	for count > 0 {
		ret = append(ret, <-quit)
		count--
	}

	fmt.Println("\x1b[1D[download file]")
	for _, name := range ret {
		fmt.Println(name)
	}
}
